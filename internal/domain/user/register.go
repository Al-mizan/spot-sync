package user

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, handler Handler) {
	api := e.Group("/api/v1/auth")

	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)
}
