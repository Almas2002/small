package usecase

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gopkg.in/gomail.v2"
	"small/internal/config"
	"small/internal/models"
	"small/internal/modules/product"
	"small/pkg/httpErrors"
	"small/pkg/tools/mail"
	"small/pkg/tracing"
	"small/pkg/type/logger"
)

type useCase struct {
	log    logger.Logger
	repo   product.IRepository
	userUc product.IUserUseCase
	config *config.Config
	mail   *gomail.Dialer
}

func New(log logger.Logger, repo product.IRepository, userUc product.IUserUseCase, cfg *config.Config) *useCase {
	return &useCase{
		log:    log,
		repo:   repo,
		userUc: userUc,
		config: cfg,
		mail:   mail.New(cfg.Email),
	}
}

func (c *useCase) CreateProduct(ctx context.Context, product *models.Product) (uint32, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.SaveProduct")
	defer span.Finish()

	c.log.Info("useCase.SaveProduct")

	candidate, err := c.repo.FindOneProduct(ctx, product)
	if err != nil {
		return 0, err
	}
	if candidate.Id != 0 {
		err = errors.Wrap(httpErrors.AlreadyExists, "product")
		tracing.TraceError(span, err)
		c.log.WarnMsg("SaveProduct.candidate", err)
		return 0, err
	}

	productId, err := c.repo.SaveProduct(ctx, product)
	if err != nil {
		return 0, err
	}

	return productId, nil

}

func (c *useCase) UpdateProductPrice(ctx context.Context, product *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.UpdateProduct")
	defer span.Finish()

	c.log.Info("useCase.UpdateProduct")

	err := c.repo.UpdateProduct(ctx, product.Id, func(old *models.Product) *models.Product {
		return &models.Product{
			Id:    old.Id,
			Title: old.Title,
			Price: product.Price,
		}
	})
	if err != nil {
		return err
	}

	users, err := c.userUc.GetProductSubscribers(ctx, product.Id)
	if err != nil {
		return err
	}
	emails := make([]string, 0, len(users))

	for _, item := range users {
		if item.Email != "" {
			emails = append(emails, item.Email)
		}

	}
	if len(emails) != 0 {
		wg := errgroup.Group{}

		for _, mailItem := range emails {
			func(mail string) {
				wg.Go(func() error {
					return c.sendEmail(mail, product.Id, product.Price)
				})
			}(mailItem)
		}
		err = wg.Wait()
		if err != nil {
			tracing.TraceError(span, err)
			c.log.WarnMsg("UpdateProduct.Email", err)
			return err
		}

	}
	return nil

}

func (c *useCase) sendEmail(mailItem string, productId uint32, productPrice float64) error {

	m := gomail.NewMessage()
	m.SetHeader("From", c.config.Email.Email)
	m.SetHeader("To", mailItem)
	m.SetHeader("Subject", "Hello From Small!")
	m.SetBody("text/html", fmt.Sprintf("product with id: %d changed price to %f", productId, productPrice))

	if err := c.mail.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (c *useCase) SubToProduct(ctx context.Context, userId, productId uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.SubToProduct")
	defer span.Finish()

	c.log.Info("useCase.SubToProduct")

	user, err := c.userUc.FindOneUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.Id == 0 {
		err = errors.Wrap(httpErrors.NotFound, "user")
		tracing.TraceError(span, err)
		c.log.WarnMsg("SubToProduct.user", err)
		return err
	}

	oneProduct, err := c.repo.FindOneProduct(ctx, &models.Product{Id: productId})
	if err != nil {
		return err
	}

	if oneProduct.Id == 0 {
		err = errors.Wrap(httpErrors.NotFound, "product")
		tracing.TraceError(span, err)
		c.log.WarnMsg("SubToProduct.user", err)
		return err
	}

	err = c.repo.UserSubToProduct(ctx, userId, productId)
	if err != nil {
		return err
	}
	return nil
}

func (c *useCase) UnSubToProduct(ctx context.Context, userId, productId uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "useCase.UnSubToProduct")
	defer span.Finish()

	c.log.Info("useCase.UnSubToProduct")

	user, err := c.userUc.FindOneUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.Id == 0 {
		err = errors.Wrap(httpErrors.NotFound, "user")
		tracing.TraceError(span, err)
		c.log.WarnMsg("UnSubToProduct.user", err)
		return err
	}

	oneProduct, err := c.repo.FindOneProduct(ctx, &models.Product{Id: productId})
	if err != nil {
		return err
	}

	if oneProduct.Id == 0 {
		err = errors.Wrap(httpErrors.NotFound, "product")
		tracing.TraceError(span, err)
		c.log.WarnMsg("UnSubToProduct.user", err)
		return err
	}
	fmt.Println()
	if err = c.repo.UserUnsubToProduct(ctx, userId, productId); err != nil {
		return err
	}
	return nil
}
