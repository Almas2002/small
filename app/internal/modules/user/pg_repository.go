package user

import (
	"context"
	"small/internal/models"
)

type IRepository interface {
	GetSubUsers(ctx context.Context, productId uint32) ([]*models.User, error)
	FinOneUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (uint32, error)
}
