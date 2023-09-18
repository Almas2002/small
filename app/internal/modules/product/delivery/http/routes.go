package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *handlers) MapRoutes() {
	h.group.POST("/", h.CreateProduct())
	h.group.PUT("/:id", h.UpdateProduct())
	h.group.POST("/sub", h.SubToProduct())
	h.group.DELETE("/unsub", h.UnSubToProduct())
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
