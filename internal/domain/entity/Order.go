package entity

type Order struct {
	OrderNumber int    `json: "orderNumber"`
	Owner       string `json: "owner"`
	Status      string `json: "status"`
}
