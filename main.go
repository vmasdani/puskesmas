package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
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
	serviceAccountFile := os.Getenv("SERVICE_ACCOUNT_FILE")

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
		"DB_HOST":              dbHost,
		"DB_USERNAME":          dbUsername,
		"DB_PASSWORD":          dbPassword,
		"DB_NAME":              dbName,
		"SERVER_PORT":          serverPort,
		"JWT_SECRET":           jwtSecret,
		"ADMIN_USERNAME":       adminUsername,
		"ADMIN_PASSWORD":       adminPassword,
		"SERVICE_ACCOUNT_FILE": serviceAccountFile,
	}, "./.env")
}

type GormModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UUID      string     `json:"uuid"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type User struct {
	GormModel
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	FcmToken string `json:"fcmToken"`
}

type UserBody struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	ChangePassword bool   `json:"changePassword"`
	NewPassword    string `json:"newPassword"`
}

func userToUserBody(user User) UserBody {
	return UserBody{
		ID:             user.ID,
		Name:           user.Name,
		Username:       user.Username,
		ChangePassword: false,
		NewPassword:    "",
	}
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

type ManpowerCategory struct {
	GormModel
	Name string `json:"name"`
}

type ManpowerStatus struct {
	GormModel
	Name string `json:"name"`
}

type ManpowerStatusAmount struct {
	GormModel
	Value                int              `json:"value"`
	ManpowerCategoryID   *uint            `json:"manpowerCategoryId"`
	ManpowerCategoryUuid string           `json:"manpowerCategoryUuid"`
	ManpowerCategory     ManpowerCategory `json:"manpowerCategory"`
	ManpowerStatusID     *uint            `json:"manpowerStatusId"`
	ManpowerStatusUuid   string           `json:"manpowerStatusUuid"`
	ManpowerStatus       ManpowerStatus   `json:"manpowerStatus"`
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

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile("service-account.json"))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	fmt.Println(app)

	tables :=
		[]interface{}{
			User{},
			UserRole{},
			Role{},
			Article{},
			AdminConfig{},
			Complaint{},
			ManpowerCategory{},
			ManpowerStatus{},
			ManpowerStatusAmount{},
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

	fcmSendMulti := func(tokens []string, title string, body string) {
		// Obtain a messaging.Client from the App.
		ctx := context.Background()
		client, err := app.Messaging(ctx)
		if err != nil {
			log.Fatalf("error getting Messaging client: %v\n", err)
		}

		// See documentation on defining a message payload.
		message := &messaging.MulticastMessage{
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Tokens: tokens,
		}
		br, err := client.SendMulticast(context.Background(), message)
		if err != nil {
			log.Fatalln(err)
		}

		// See the BatchResponse reference documentation
		// for the contents of response.
		fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	}

	r.HandleFunc("/fcm-test", func(w http.ResponseWriter, r *http.Request) {

		// Obtain a messaging.Client from the App.
		ctx := context.Background()
		client, err := app.Messaging(ctx)
		if err != nil {
			log.Fatalf("error getting Messaging client: %v\n", err)
		}

		// See documentation on defining a message payload.
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: "Price drop",
				Body:  "5% off all electronics",
			},
			Token: "eWvUQpXgQem462q6qtJljm:APA91bFH6eDlMqPzLWolT6HwO-4HueRbD-xxIYEE9WQeIc0af4vwdKD78qv4BEo1EJFQHOFp3behskz1opTces-dpMIo5-20iwAOHs8Kjkp8ldtDIoL2L5zaCyzSk8Q-ZnsYEV2wdzIY",
		}
		// Send a message to devices subscribed to the combination of topics
		// specified by the provided condition.
		response, err := client.Send(ctx, message)
		if err != nil {
			log.Fatalln(err)
		}
		// Response is a message ID string.
		fmt.Println("Successfully sent message:", response)
	})

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

			tokens := func() []string {
				var users []User
				db.Find(&users)

				tokens := []string{}

				for _, user := range users {
					if user.FcmToken != "" {
						tokens = append(tokens, user.FcmToken)
					}
				}

				return tokens
			}()

			if len(tokens) > 0 {
				fcmSendMulti(tokens, "Aduan Baru", complaint.Name+": "+complaint.Complaint)
			}
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

	// r.PathPrefix("/admin").Handler(http.StripPrefix("/admin", http.FileServer(http.Dir("./admin"))))
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
	// r.HandleFunc("/complaints-save", func(w http.ResponseWriter, r *http.Request) {
	// 	auth := r.Header.Get("authorization")
	// 	fmt.Println("Auth token:", auth)

	// 	token, err := jwt.ParseWithClaims(auth, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 		return []byte(os.Getenv("JWT_SECRET")), nil
	// 	})

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		fmt.Fprintf(w, "Error decoding token")
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 	}

	// 	fmt.Println(token.Claims)
	// })

	getTokenId := func(auth string) (float64, error) {
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return 0, err
		}

		fmt.Println(claims["jti"])
		// fmt.Println(claims["jti"].(string))
		fmt.Println(claims["jti"].(float64))

		claimsFloat, ok := claims["jti"].(float64)

		if ok {
			return claimsFloat, nil
		}

		return 0, err

	}

	type FcmTokenBody struct {
		Token string `json:"token"`
	}

	r.HandleFunc("/save-fcm-token", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		var fcmToken FcmTokenBody
		json.NewDecoder(r.Body).Decode(&fcmToken)
		fmt.Println("Auth token:", auth)

		id, err := getTokenId(auth)

		fmt.Println(err, id)

		var user User
		if db.Where("id = ?", id).First(&user).Error != nil {
			fmt.Println("Error fetching user")
			return
		}

		fmt.Println(user.Name)
		user.FcmToken = fcmToken.Token

		db.Save(&user)
	})

	r.HandleFunc("/users-jwt", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		var fcmToken FcmTokenBody
		json.NewDecoder(r.Body).Decode(&fcmToken)
		fmt.Println("Auth token:", auth)

		id, err := getTokenId(auth)

		fmt.Println(err, id)

		var user User
		if db.Where("id = ?", id).First(&user).Error != nil {
			fmt.Println("Error fetching user")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "User not found")
			return
		}

		json.NewEncoder(w).Encode(userToUserBody(user))
	})

	type ManpowerData struct {
		ManpowerCategories        []ManpowerCategory `json:"manpowerCategories"`
		ManpowerCategoryDeleteIds []uint             `json:"manpowerCategoryDeleteIds"`

		ManpowerStatuses        []ManpowerStatus `json:"manpowerStatuses"`
		ManpowerStatusDeleteIds []uint           `json:"manpowerStatusDeleteIds"`

		ManpowerStatusAmounts         []ManpowerStatusAmount `json:"manpowerStatusAmounts"`
		ManpowerStatusAmountDeleteIds []uint                 `json:"manpowerStatusAmountDeleteIds"`
	}

	r.HandleFunc("/manpowers-save", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		fmt.Println("Auth token:", auth)

		admin := CheckAdmin(auth)

		if !admin {
			fmt.Println(err)
			fmt.Fprintf(w, "Error decoding token. Not admin")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("content-type", "application/json")

		var manpowerData ManpowerData
		json.NewDecoder(r.Body).Decode(&manpowerData)

		for _, manpowerCategory := range manpowerData.ManpowerCategories {
			db.Save(&manpowerCategory)
		}

		for _, id := range manpowerData.ManpowerCategoryDeleteIds {
			db.Delete(&ManpowerCategory{}, id)
		}

		for _, manpowerStatus := range manpowerData.ManpowerStatuses {
			db.Save(&manpowerStatus)
		}

		for _, id := range manpowerData.ManpowerStatusDeleteIds {
			db.Delete(&ManpowerStatus{}, id)
		}

		for _, manpowerStatusAmount := range manpowerData.ManpowerStatusAmounts {
			db.Save(&manpowerStatusAmount)
		}

		for _, id := range manpowerData.ManpowerStatusAmountDeleteIds {
			db.Delete(&ManpowerStatusAmount{}, id)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(manpowerData)

	}).Methods("POST")
	r.HandleFunc("/manpowercategories", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		switch r.Method {
		case "GET":
			var manpowerCategories []ManpowerCategory
			db.Find(&manpowerCategories)
			json.NewEncoder(w).Encode(manpowerCategories)

		case "POST":
			auth := r.Header.Get("authorization")
			fmt.Println("Auth token:", auth)

			admin := CheckAdmin(auth)

			if !admin {
				fmt.Println(err)
				fmt.Fprintf(w, "Error decoding token. Not admin")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			fmt.Println("POST detect")
			var manpowerCategory []ManpowerCategory
			db.Save(&manpowerCategory)
			json.NewEncoder(w).Encode(&manpowerCategory)
		}
	}).Methods("GET", "POST")
	r.HandleFunc("/manpowerstatuses", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		var manpowerStatuses []ManpowerStatus
		db.Find(&manpowerStatuses)
		json.NewEncoder(w).Encode(manpowerStatuses)
	}).Methods("GET")
	r.HandleFunc("/manpowerstatusamounts", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		var manpowerStatusAmounts []ManpowerStatusAmount
		db.Find(&manpowerStatusAmounts)
		json.NewEncoder(w).Encode(manpowerStatusAmounts)
	}).Methods("GET")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		type LoginStruct struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var loginData LoginStruct
		json.NewDecoder(r.Body).Decode(&loginData)

		tokenString := ""

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if loginData.Username != os.Getenv("ADMIN_USERNAME") {
			var foundUser User

			if db.Where("lower(username) = ?", strings.ToLower(loginData.Username)).First(&foundUser).Error != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Username tidak ditemukan!")
				return
			}

			err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginData.Password))

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Password salah!")
				return
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp":   time.Now().Unix() + 86400*365*30,
				"admin": false,
				"jti":   foundUser.ID,
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))

			fmt.Println(tokenString, err)
		} else {
			if loginData.Password != os.Getenv("ADMIN_PASSWORD") {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Password salah!")
				return
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp":   time.Now().Unix() + 86400*365*30,
				"admin": true,
				"jti":   0,
			})

			// Sign and get the complete encoded token as a string using the secret
			var err error
			tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))

			fmt.Println(tokenString, err)

		}

		fmt.Fprintf(w, "%s", tokenString)
	}).Methods("POST")
	r.HandleFunc("/users-view", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		var users []User
		db.Find(&users)

		usersMapped := []UserBody{}

		for _, user := range users {
			usersMapped = append(usersMapped, userToUserBody(user))
		}

		json.NewEncoder(w).Encode(usersMapped)
	}).Methods("GET")

	type UserSave struct {
		UserBody      []UserBody `json:"userBody"`
		UserDeleteIds []uint     `json:"userDeleteIds"`
	}

	r.HandleFunc("/duplicate-remove", func(w http.ResponseWriter, r *http.Request) {
		var manpowerCategories []ManpowerCategory
		db.Find(&manpowerCategories)

		var wg sync.WaitGroup

		for _, category := range manpowerCategories {
			// fmt.Println(category)

			wg.Add(1)

			go func(category ManpowerCategory) {
				defer wg.Done()

				var foundManpowerAmounts []ManpowerStatusAmount
				db.Where("manpower_category_uuid = ?", category.UUID).Find(&foundManpowerAmounts)

				if len(foundManpowerAmounts) > 1 {
					for i, manpowerStatusAmount := range foundManpowerAmounts {
						if i > 0 {
							fmt.Println(category.Name, manpowerStatusAmount.ID, manpowerStatusAmount.Value)
							db.Delete(&ManpowerStatusAmount{}, manpowerStatusAmount.ID)
						}
					}
				}

			}(category)

			wg.Wait()

		}
	})

	r.HandleFunc("/users-save", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		var userSave UserSave
		json.NewDecoder(r.Body).Decode(&userSave)

		fmt.Println(userSave.UserDeleteIds)

		for _, user := range userSave.UserBody {
			var foundUser User
			if db.Where("username = ?", user.Username).Find(&foundUser).Error != nil {
				newUser := User{
					Username: user.Username,
					Name:     user.Name,
				}
				newUser.ID = user.ID

				if user.ChangePassword {
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), bcrypt.DefaultCost)

					if err != nil {
						fmt.Println("Error gen password")
					}

					newUser.Password = string(hashedPassword)
				}

				db.Save(&newUser)

			} else {
				fmt.Println(user.Name, user.Username, foundUser.Username)

				foundUser.ID = user.ID
				foundUser.Name = user.Name
				foundUser.Username = user.Username

				if user.ChangePassword {
					hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), bcrypt.DefaultCost)

					if err != nil {
						fmt.Println("Error gen password")
					}

					foundUser.Password = string(hashedPassword)
				}

				fmt.Println(foundUser)

				db.Save(&foundUser)
			}
		}

		for _, userDeleteId := range userSave.UserDeleteIds {
			fmt.Println(userDeleteId)
			db.Delete(&User{}, userDeleteId)
		}

		json.NewEncoder(w).Encode(userSave)
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
