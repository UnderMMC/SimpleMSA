package repository

import (
	"database/sql"
	"github.com/SimpleMSA/internal/domain/entity"
	"log"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) GetOrderStatus(order entity.Order) (entity.Order, error) {
	var status string
	err := r.db.QueryRow("SELECT status FROM orders WHERE ordernum=$1", order.OrderNumber).Scan(&status)
	if err != nil {
		log.Fatal()
	}
	order.Status = status
	return order, nil
}
