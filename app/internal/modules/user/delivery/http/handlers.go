package http

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"small/internal/config"
	"small/internal/models"
	"small/internal/modules/user"
	"small/internal/modules/user/delivery/http/dto"
	"small/pkg/httpErrors"
	"small/pkg/tracing"
	"small/pkg/type/email"
	"small/pkg/type/logger"
)

type handlers struct {
	group *echo.Group
	log   logger.Logger
	cfg   *config.Config
	uc    user.IUseCase
	v     *validator.Validate
}

func NewHandlers(group *echo.Group, log logger.Logger, cfg *config.Config, uc user.IUseCase, v *validator.Validate) *handlers {
	return &handlers{
		group: group,
		log:   log,
		cfg:   cfg,
		uc:    uc,
		v:     v,
	}
}

func (h *handlers) Registration() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := tracing.StartHttpServerTraceSpan(c, "handlers.Registration")
		defer span.Finish()

		createDto := &dto.CreateUserDto{}
		if err := c.Bind(&createDto); err != nil {
			tracing.TraceError(span, err)
			h.log.WarnMsg("Bind", err)
			return c.JSON(400, err)
		}

		if createDto.Email != "" {
			if err := email.New(createDto.Email); err != nil {
				tracing.TraceError(span, err)
				h.log.WarnMsg("Email", err)
				return c.JSON(422, err)
			}
		}

		userId, err := h.uc.Registration(ctx, &models.User{
			Phone: createDto.Phone,
			Email: createDto.Email,
		})
		if err != nil {
			if errors.Is(err, httpErrors.AlreadyExists) {
				return c.JSON(409, err)
			}
			if errors.Is(err, httpErrors.InvalidArguments) {
				return c.JSON(422, err)
			}
			return c.JSON(500, err)

		}
		return c.JSON(201, userId)

	}
}
