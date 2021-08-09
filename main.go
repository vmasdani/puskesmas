package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Setup database
	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Opening DB error", err)
		fmt.Println(db)

		return
	}

	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./landing/")))

	serverPort := os.Getenv("SERVER_PORT")

	fmt.Println("Running on port", serverPort)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+serverPort, r))
}
