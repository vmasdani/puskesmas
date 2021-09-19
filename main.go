package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func generateJwtSecret() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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
	jwtSecret := os.Getenv("JWT_SECRET")
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if serverPort == "" {
		serverPort = "8000"
	}

	if jwtSecret == "" {
		jwtSecret = generateJwtSecret()
	}

	if adminUsername == "" {
		panic("Admin username cannot be empty!")
	}

	if adminPassword == "" {
		panic("Admin password cannot be empty!")
	}

	godotenv.Write(map[string]string{
		"DB_HOST":        dbHost,
		"DB_USERNAME":    dbUsername,
		"DB_PASSWORD":    dbPassword,
		"DB_NAME":        dbName,
		"SERVER_PORT":    serverPort,
		"JWT_SECRET":     jwtSecret,
		"ADMIN_USERNAME": adminUsername,
		"ADMIN_PASSWORD": adminPassword,
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

	r.Use(AuthMiddleware)

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

	CheckAdmin := func(auth string) bool {
		fmt.Println("Auth token:", auth)

		token, err := jwt.ParseWithClaims(auth, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return false
		}

		// Cast
		admin, ok := token.Claims.(jwt.MapClaims)["admin"].(bool)

		if !ok {
			return false
		}

		fmt.Println("Is admin:", admin)
		fmt.Println(token.Claims)

		return admin
	}

	r.PathPrefix("/admin").Handler(http.StripPrefix("/admin", http.FileServer(http.Dir("./admin"))))
	r.HandleFunc("/authorize-admin", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")

		admin := CheckAdmin(auth)

		if !admin {
			fmt.Println(err)
			fmt.Fprintf(w, "Error decoding token. Not admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
	r.HandleFunc("/complaints-save", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")

		admin := CheckAdmin(auth)

		if !admin {
			fmt.Println(err)
			fmt.Fprintf(w, "Error decoding token. Not admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var complaint Complaint
		json.NewDecoder(r.Body).Decode(&complaint)

		db.Save(&complaint)

		json.NewEncoder(w).Encode(complaint)
		w.WriteHeader(http.StatusCreated)
	})

	r.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		fmt.Println("Auth token:", auth)

		token, err := jwt.ParseWithClaims(auth, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		fmt.Println(token.Claims)

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Error decoding token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	})
	r.HandleFunc("/complaints-save", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		fmt.Println("Auth token:", auth)

		token, err := jwt.ParseWithClaims(auth, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Error decoding token")
			w.WriteHeader(http.StatusUnauthorized)
		}

		fmt.Println(token.Claims)
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		type LoginStruct struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var loginData LoginStruct
		json.NewDecoder(r.Body).Decode(&loginData)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp":   time.Now().Unix() + 86400*365*30,
			"admin": true,
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

		fmt.Println(tokenString, err)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if loginData.Username != os.Getenv("ADMIN_USERNAME") || loginData.Password != os.Getenv("ADMIN_PASSWORD") {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Username atau password salah!")
			return
		}

		fmt.Fprintf(w, "%s", tokenString)
	}).Methods("POST")

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
