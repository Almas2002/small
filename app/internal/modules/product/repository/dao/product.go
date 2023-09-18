package dao

type Product struct {
	Id    uint32  `db:"product_id"`
	Title string  `db:"title"`
	Price float64 `db:"price"`
}

const ProductTableName = "products"

var ProductsSelect = []string{
	"product_id",
	"title",
	"price",
}

var ProductInsert = []string{
	"title",
	"price",
}
