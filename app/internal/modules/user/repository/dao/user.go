package dao

type User struct {
	Id    uint32 `db:"user_id"`
	Phone string `db:"phone"`
	Email string `db:"email"`
}

const UserTableName = "users"

const UserSubProductsTableName = "user_sub_products"

var UsersSelect = []string{
	"user_id",
	"phone",
	"email",
}

var UserInsert = []string{
	"phone",
	"email",
}

var UserSubToProductInsert = []string{
	"user_id",
	"product_id",
}
