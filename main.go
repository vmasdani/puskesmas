package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./landing/")))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
