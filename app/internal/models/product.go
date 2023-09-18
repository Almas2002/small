package models

type Product struct {
	Id    uint32  `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}
