package reservation

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, handler Handler, authMw echo.MiddlewareFunc, adminMw echo.MiddlewareFunc) {
	api := e.Group("/api/v1/reservations", authMw)

	// Auth routes
	api.POST("", handler.CreateReservation)
	api.GET("/my-reservations", handler.GetMyReservations)
	api.DELETE("/:id", handler.CancelReservation)

	// Admin routes
	api.GET("", handler.GetAllReservations, adminMw)
}
