package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./landing/")))

	serverPort := os.Getenv("SERVER_PORT")

	fmt.Println("Running on port", serverPort)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+serverPort, r))
}
