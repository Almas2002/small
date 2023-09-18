package product

import (
	"context"
	"small/internal/models"
)

type IUseCase interface {
	CreateProduct(ctx context.Context, product *models.Product) (uint32, error)
	UpdateProductPrice(ctx context.Context, product *models.Product) error
	UnSubToProduct(ctx context.Context, userId, productId uint32) error
	SubToProduct(ctx context.Context, userId, productId uint32) error
}

type IUserUseCase interface {
	GetProductSubscribers(ctx context.Context, productId uint32) ([]*models.User, error)
	FindOneUserById(ctx context.Context, id uint32) (*models.User, error)
}
