package entity

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}
