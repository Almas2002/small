package usecase

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"small/internal/models"
	"small/internal/modules/user"
	"small/pkg/httpErrors"
	"small/pkg/tracing"
	"small/pkg/type/logger"
)

type useCase struct {
	log  logger.Logger
	repo user.IRepository
}

func New(log logger.Logger, repo user.IRepository) *useCase {
	return &useCase{
		log:  log,
		repo: repo,
	}
}

func (c *useCase) Registration(ctx context.Context, user *models.User) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.Registration")
	defer span.Finish()

	c.log.Info("useCase.Registration")

	if user.Phone == "" && user.Email == "" {
		err := errors.Wrap(httpErrors.InvalidArguments, "phone and email not must be empty")
		tracing.TraceError(span, err)
		c.log.WarnMsg("Registration.Valid", err)
		return 0, httpErrors.InvalidArguments
	}

	candidate, err := c.repo.FinOneUser(ctx, user)
	if err != nil {
		return 0, err
	}

	if candidate.Id != 0 {
		err = errors.Wrap(httpErrors.AlreadyExists, "user")
		tracing.TraceError(span, err)
		c.log.WarnMsg("Registration.candidate", err)
		return 0, err
	}

	userId, err := c.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (c *useCase) FindOneUserById(ctx context.Context, id uint32) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.FindOneById")
	defer span.Finish()

	c.log.Info("useCase.FindOneById")

	oneUser, err := c.repo.FinOneUser(ctx, &models.User{Id: id})
	if err != nil {
		return nil, err
	}
	return oneUser, nil
}

func (c *useCase) GetProductSubscribers(ctx context.Context, productId uint32) ([]*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.GetProductSubscribers")
	defer span.Finish()

	c.log.Info("useCase.GetProductSubscribers")

	users, err := c.repo.GetSubUsers(ctx, productId)
	if err != nil {
		return nil, err
	}
	return users, err
}
