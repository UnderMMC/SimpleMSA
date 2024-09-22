package app

import (
	"database/sql"
	"encoding/json"
	"github.com/SimpleMSA/internal/domain/entity"
	"github.com/SimpleMSA/internal/domain/repository"
	"github.com/SimpleMSA/internal/domain/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync"
	"time"
)

type Service interface {
	Registration(user entity.User) error
	Authorization(user entity.User) error
}

type App struct {
	serv Service
}

var db *sql.DB
var secretKey = []byte("jwt_token_example")

type AuthResponse struct {
	Token string `json:"token"`
}

func generateToken(user entity.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   user.Login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
func (a *App) registrHandler(w http.ResponseWriter, r *http.Request) {
	var regUser entity.User
	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.serv.Registration(regUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	return
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.serv.Authorization(user)
	if err != nil {
		log.Fatal(err)
	}
	token, err := generateToken(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
	return
}

func (a *App) validatorHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.Token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	tokenString := request.Token

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user := entity.User{
		Login: claims.Subject,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func New() *App {
	return &App{}
}

func (a *App) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	connStr := "user=postgres password=pgpwd4habr dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresUserRepository(db)
	serv := service.NewUserService(repo)
	a.serv = serv

	r := mux.NewRouter()

	r.HandleFunc("/reg", a.registrHandler).Methods("POST")
	r.HandleFunc("/login", a.loginHandler).Methods("POST")
	r.HandleFunc("/validate", a.validatorHandler).Methods("POST")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
