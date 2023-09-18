package dto

type SubUnsubProductDto struct {
	UserId    uint32 `json:"user_id" validate:"required,numeric"`
	ProductId uint32 `json:"product_id" validate:"required,numeric"`
}
