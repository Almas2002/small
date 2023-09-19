package http

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"small/internal/config"
	"small/internal/models"
	"small/internal/modules/product"
	"small/internal/modules/product/delivery/http/dto"
	"small/pkg/httpErrors"
	"small/pkg/tracing"
	"small/pkg/type/logger"
	"strconv"
)

type handlers struct {
	group *echo.Group
	log   logger.Logger
	cfg   *config.Config
	uc    product.IUseCase
	v     *validator.Validate
}

func NewHandlers(group *echo.Group, log logger.Logger, cfg *config.Config, uc product.IUseCase, v *validator.Validate) *handlers {
	return &handlers{
		group: group,
		log:   log,
		cfg:   cfg,
		uc:    uc,
		v:     v,
	}
}

func (h *handlers) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := tracing.StartHttpServerTraceSpan(c, "handlers.CreateProduct")
		defer span.Finish()

		createDto := &dto.CreateProduct{}
		if err := c.Bind(&createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			return c.JSON(400, err)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			return c.JSON(400, err)
		}

		productId, err := h.uc.CreateProduct(ctx, &models.Product{
			Title: createDto.Title,
			Price: createDto.Price,
		})

		if err != nil {
			if errors.Is(err, httpErrors.AlreadyExists) {
				return c.JSON(400, err)
			}
			return c.JSON(500, err)
		}

		return c.JSON(201, productId)

	}
}

func (h *handlers) SubToProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := tracing.StartHttpServerTraceSpan(c, "handlers.SubToProduct")
		defer span.Finish()

		createDto := &dto.SubUnsubProductDto{}
		if err := c.Bind(&createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			return c.JSON(400, err)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			return c.JSON(400, err)
		}

		err := h.uc.SubToProduct(ctx, createDto.UserId, createDto.ProductId)
		if err != nil {
			if errors.Is(err, httpErrors.NotFound) {
				return c.JSON(404, err.Error())
			}
			if errors.Is(err, httpErrors.AlreadyExists) {
				return c.JSON(400, err.Error())
			}
			return c.JSON(500, err.Error())
		}
		return c.JSON(201, "OK")
	}
}

func (h *handlers) UnSubToProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := tracing.StartHttpServerTraceSpan(c, "handlers.UnSubToProduct")
		defer span.Finish()

		createDto := &dto.SubUnsubProductDto{}
		if err := c.Bind(&createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			return c.JSON(400, err)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			return c.JSON(400, err)
		}

		if err := h.uc.UnSubToProduct(ctx, createDto.UserId, createDto.ProductId); err != nil {
			if errors.Is(err, httpErrors.NotFound) {
				return c.JSON(404, err)
			}
			return c.JSON(500, err)
		}

		return c.JSON(201, "OK")
	}
}

func (h *handlers) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := tracing.StartHttpServerTraceSpan(c, "handlers.UpdateProduct")
		defer span.Finish()
		var productId uint32

		if limit, err := strconv.Atoi(c.Param("id")); err == nil && limit != 0 {
			productId = uint32(limit)
		} else {
			h.log.WarnMsg("param", err)
			return c.JSON(400, errors.Wrap(httpErrors.InvalidArguments, "product id param"))
		}
		updateProduct := &dto.UpdateProductDto{}
		if err := c.Bind(&updateProduct); err != nil {
			h.log.WarnMsg("Bind", err)
			return c.JSON(400, err)
		}

		if err := h.v.StructCtx(ctx, updateProduct); err != nil {
			h.log.WarnMsg("validate", err)
			return c.JSON(400, err)
		}

		if err := h.uc.UpdateProductPrice(ctx, &models.Product{Id: productId, Price: updateProduct.Price}); err != nil {
			if errors.Is(err, httpErrors.NotFound) {
				return c.JSON(404, err)
			}
			return c.JSON(500, err)
		}

		return c.JSON(201, "OK")
	}
}
