package user

import (
	"context"
	"small/internal/models"
)

type IUseCase interface {
	GetProductSubscribers(ctx context.Context, productId uint32) ([]*models.User, error)
	FindOneUserById(ctx context.Context, id uint32) (*models.User, error)
	Registration(ctx context.Context, user *models.User) (uint32, error)
}
