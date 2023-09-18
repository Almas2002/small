package product

import (
	"context"
	"small/internal/models"
)

type IRepository interface {
	SaveProduct(ctx context.Context, product *models.Product) (uint32, error)
	FindOneProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, id uint32, updateFn func(old *models.Product) *models.Product) error
	UserUnsubToProduct(ctx context.Context, userId uint32, productId uint32) error
	UserSubToProduct(ctx context.Context, userId uint32, productId uint32) error
}
