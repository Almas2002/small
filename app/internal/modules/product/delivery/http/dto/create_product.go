package dto

type CreateProduct struct {
	Title string  `json:"title" validate:"required"`
	Price float64 `json:"price" validate:"required,numeric"`
}
