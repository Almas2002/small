package models

type User struct {
	Id    uint32 `json:"id"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}
