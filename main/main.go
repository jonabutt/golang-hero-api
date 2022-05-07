package main

import (
	"log"
	"net/http"
)

// GET HEROS -- /heros -- GET
// GET HERO -- /heros/{id} -- GET
// ADD HERO -- /heros -- POST
// DELTE HERO -- /heros/{id} -- DELETE
// UPDATE HERO -- /heros -- PUT


func main(){
	mux := http.NewServeMux()
	log.Fatal(http.ListenAndServe(":8081",mux))
}