package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CheckEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8000"
	}

	godotenv.Write(map[string]string{
		"DB_HOST":     dbHost,
		"DB_USERNAME": dbUsername,
		"DB_PASSWORD": dbPassword,
		"DB_NAME":     dbName,
		"SERVER_PORT": serverPort,
	}, "./.env")
}

type GormModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type User struct {
	GormModel
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserRole struct {
	GormModel

	UserID uint `json:"userId"`
	User   User `json:"user"`
	RoleID uint `json:"roleId"`
	Role   Role `json:"role"`
}

type Role struct {
	GormModel
	Name string `json:"name"`
}

type Article struct {
	GormModel
	Date    time.Time `json:"date"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type AdminConfig struct {
	GormModel
	FacebookUrl    string `json:"facebookURL"`
	InstagramUrl   string `json:"instagramURL"`
	WhatsappNumber string `json:"whatsappNumber"`
}

type Complaint struct {
	GormModel
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Answer    string `json:"answer"`
	Complaint string `json:"complaint"`
}

func main() {
	CheckEnv()

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

	tables :=
		[]interface{}{
			User{},
			UserRole{},
			Role{},
			Article{},
			AdminConfig{},
			Complaint{},
		}

	for _, table := range tables {
		db.AutoMigrate(table)
	}

	if err != nil {
		fmt.Println("Opening DB error", err)
		fmt.Println(db)

		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/complaints", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("content-type", "application/json")
			fmt.Println("GET")

			var complaints []Complaint
			db.Find(&complaints)

			json.NewEncoder(w).Encode(&complaints)
		case http.MethodPost:
			var complaint Complaint
			json.NewDecoder(r.Body).Decode(&complaint)

			db.Save(&complaint)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(&complaint)
		default:
			fmt.Println("[complaints] method irrelevant")
		}
	}).Methods("GET", "POST")

	r.PathPrefix("/admin").Handler(http.StripPrefix("/admin", http.FileServer(http.Dir("./admin"))))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./landing/")))

	serverPort := os.Getenv("SERVER_PORT")

	fmt.Println("Running on ", "http://localhost:"+serverPort)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":"+serverPort, cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	}).Handler(r)))
}
