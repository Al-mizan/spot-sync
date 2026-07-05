package zone

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, handler Handler, authMw echo.MiddlewareFunc, adminMw echo.MiddlewareFunc) {
	api := e.Group("/api/v1/zones")

	// Public routes
	api.GET("", handler.GetAllZones)
	api.GET("/:id", handler.GetZoneByID)

	// Admin only routes
	api.POST("", handler.CreateZone, authMw, adminMw)
	api.PATCH("/:id", handler.UpdateZone, authMw, adminMw)
	api.DELETE("/:id", handler.DeleteZone, authMw, adminMw)
}
