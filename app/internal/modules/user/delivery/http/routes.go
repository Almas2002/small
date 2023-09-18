package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *handlers) MapRoutes() {
	h.group.POST("/", h.Registration())
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
