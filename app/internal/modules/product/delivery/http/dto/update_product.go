package dto

type UpdateProductDto struct {
	Price float64 `json:"price" validate:"required,numeric"`
}
