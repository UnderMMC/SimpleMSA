package app

import (
	"database/sql"
	"encoding/json"
	"github.com/SimpleMSA/internal/domain/entity"
	"github.com/SimpleMSA/internal/domain/repository"
	"github.com/SimpleMSA/internal/domain/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

type OrderService interface {
	OrderStatus(order entity.Order) (entity.Order, error)
}

type OrderApp struct {
	OderServ OrderService
}

func getUserFromToken(token string) (entity.User, error) {
	var user entity.User
	resp, err := http.Get("http://localhost:8080/validate?token=" + token)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (o *OrderApp) OrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}
	var order entity.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	order, err = o.OderServ.OrderStatus(order)

	var user entity.User
	user, err = getUserFromToken(token)
	if err != nil {
		http.Error(w, "Could not retrieve user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order.Status)
}

func NewOrder() *OrderApp {
	return &OrderApp{}
}

func (o *OrderApp) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	connStr := "user=postgres password=pgpwd4habr dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(db)
	orderServ := service.NewOrderService(orderRepo)
	o.OderServ = orderServ

	r := mux.NewRouter()
	r.HandleFunc("/order", o.OrderStatusHandler)

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
