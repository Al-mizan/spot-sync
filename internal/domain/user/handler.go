package user

import (
	"spotsync/internal/apperror"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

// Handler defines the contract for user HTTP handlers.
type Handler interface {
	Register(c *echo.Context) error
	Login(c *echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Register(c *echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.Register(req)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 201, "User registered successfully", response)
}

func (h *handler) Login(c *echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.Login(req)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "Login successful", response)
}
