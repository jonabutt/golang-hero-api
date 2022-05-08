package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("temp todo create"))
}
// DELETE HERO -- /heros/{id} -- DELETE
func (h *superHeroHandler) delete(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("temp todo delete"))
}
// UPDATE HERO -- /heros -- PUT
func (h *superHeroHandler) update(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("temp todo update"))
}

func internalServerError(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error":"internal server error"}`))
}

func notFound(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"not found"}`))
}

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
	log.Fatal(http.ListenAndServe(":8081",mux))
}