package httpresponse

import (
	"github.com/labstack/echo/v5"
)

// SuccessResponse is the standardized success response format.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse is the standardized error response format.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

// Success sends a standardized success JSON response.
func Success(c *echo.Context, code int, message string, data any) error {
	return c.JSON(code, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends a standardized error JSON response.
func Error(c *echo.Context, code int, message string, errors any) error {
	return c.JSON(code, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
