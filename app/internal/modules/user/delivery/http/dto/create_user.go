package dto

type CreateUserDto struct {
	Phone string `json:"phone"`
	Email string `json:"email" validate:""`
}
