package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	getAndDeleteHeroRe = regexp.MustCompile(`^\/heros\/(\d+)$`)
	createUpdateAndListHeroRe = regexp.MustCompile(`^\/heros[\/]*$`)
)


type superhero struct {
	ID string `json:"id"`
	Name string `json:"name"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Place string `json:"place"`
}

type datastore struct {
	m map[string]superhero
	*sync.RWMutex
}

type superHeroHandler struct {
	store *datastore
}

type credentials struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

type userClaims struct{
	Username string `json:"username"`
	jwt.StandardClaims
}

type authHandler struct {
	key []byte
}

func validateCredentials(loginDetails *credentials) bool {
	return loginDetails.Username == "admin" && loginDetails.Password == "secret"
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	switch{
	case r.Method == http.MethodPost:
		h.auth(w,r)
		return
	default:
		notFound(w,r)
		return
	}
}

// POST AUTH issue an JWT TOKEN 
func (h *authHandler) auth(w http.ResponseWriter, r *http.Request){
	// get json object
	decoder := json.NewDecoder(r.Body)
	var loginCred credentials
	er := decoder.Decode(&loginCred)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	// check username and password are valid
	if(!validateCredentials(&loginCred)){
		unAuthorized(w,r)
		return
	}
	// generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,userClaims{
		Username: loginCred.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0,0,1).Unix(),
			Issuer: "http://localhost:8081",
		},
	})
	tkn,err := token.SignedString(h.key)
	if(err != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	jwtJson, er := json.Marshal(struct{JWT string}{JWT: tkn})
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jwtJson)
}

func (h *superHeroHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type","application/json")
	switch{
	case r.Method == http.MethodGet && createUpdateAndListHeroRe.MatchString(r.URL.Path):
		h.list(w,r)
		return
	case r.Method == http.MethodGet && getAndDeleteHeroRe.MatchString(r.URL.Path):
		h.get(w,r)
		return
	case r.Method == http.MethodPost && createUpdateAndListHeroRe.MatchString(r.URL.Path):
		h.create(w,r)
		return
	case r.Method == http.MethodDelete && getAndDeleteHeroRe.MatchString(r.URL.Path):
		h.delete(w,r)
		return
	case r.Method == http.MethodPut && createUpdateAndListHeroRe.MatchString(r.URL.Path):
		h.update(w,r)
		return
	default:
		notFound(w,r)
		return
	}
}

// GET HEROS -- /heros -- GET
func (h *superHeroHandler) list(w http.ResponseWriter, r *http.Request){
	// convert map to slice
	h.store.RLock()
	heros := make([]superhero,0,len(h.store.m))
	for _,hero := range h.store.m{
		heros = append(heros, hero)
	}
	h.store.RUnlock()
	// convert slice to json and return it
	herosJson, er := json.Marshal(heros)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(herosJson)
}
// GET HERO -- /heros/{id} -- GET
func (h *superHeroHandler) get(w http.ResponseWriter, r *http.Request){
	// get id from request
	id := strings.TrimPrefix(r.URL.Path, "/heros/")
	// get hero object from map
	h.store.RLock()
	hero, ok := h.store.m[id];
	h.store.RUnlock()
	if !ok {
		notFound(w,r)
		return
	}
	// change hero object to json bytes
	heroJson, er := json.Marshal(hero)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(heroJson)
}
// ADD HERO -- /heros -- POST
func (h *superHeroHandler) create(w http.ResponseWriter, r *http.Request){
	// get hero json data from request to object
	decoder := json.NewDecoder(r.Body)
	var newHero superhero
	er := decoder.Decode(&newHero)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	// add hero to the map
	h.store.RWMutex.Lock()
	h.store.m[newHero.ID] = newHero
	h.store.RWMutex.Unlock()
	// return the hero object as json
	heroJson, er := json.Marshal(newHero)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(heroJson)
}
// DELETE HERO -- /heros/{id} -- DELETE
func (h *superHeroHandler) delete(w http.ResponseWriter, r *http.Request){
	// get id from request
	id := strings.TrimPrefix(r.URL.Path, "/heros/")
	// delete hero from map
	h.store.RWMutex.Lock()
	delete(h.store.m, id)
	h.store.RWMutex.Unlock()	
	w.WriteHeader(http.StatusNoContent)
}
// UPDATE HERO -- /heros -- PUT
func (h *superHeroHandler) update(w http.ResponseWriter, r *http.Request){
	// get hero json data from request to object
	decoder := json.NewDecoder(r.Body)
	var updatedHero superhero
	er := decoder.Decode(&updatedHero)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	// get hero from map
	h.store.RWMutex.RLock()
	_, ok  := h.store.m[updatedHero.ID]
	h.store.RWMutex.RUnlock()
	// if there is no hero return not found
	if(!ok){
		notFound(w,r)
		return
	}
	// update hero
	h.store.RWMutex.Lock()
	h.store.m[updatedHero.ID] = updatedHero
	h.store.RWMutex.Unlock()
	// return updated hero as json
	heroJson, er := json.Marshal(updatedHero)
	if(er != nil){
		// return server error
		internalServerError(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(heroJson)
}

func internalServerError(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error":"internal server error"}`))
}

func notFound(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"not found"}`))
}

func unAuthorized(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"unauthorized"}`))
}

var secretKey = []byte("")

func main(){
	mux := http.NewServeMux()
	superHeroH := &superHeroHandler{
		store : &datastore{
			m : map[string]superhero{
				"1":{ID: "1",Name: "SuperMan",FirstName: "Clark Joseph",LastName: "Kent",Place: "Smallville"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	mux.Handle("/heros", superHeroH)
	mux.Handle("/heros/", superHeroH)
	authH := &authHandler{
		key: secretKey,
	}
	mux.Handle("/auth/",authH)
	mux.Handle("/auth", authH)
	log.Fatal(http.ListenAndServe(":8081",mux))
}