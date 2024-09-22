package service

import (
	"github.com/SimpleMSA/internal/domain/entity"
	"log"
)

type OrderRepository interface {
	GetOrderStatus(order entity.Order) (entity.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) OrderStatus(order entity.Order) (entity.Order, error) {
	var err error
	order, err = s.repo.GetOrderStatus(order)
	if err != nil {
		log.Fatal()
	}
	return order, err
}
