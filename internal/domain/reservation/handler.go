package reservation

import (
	"spotsync/internal/apperror"
	"spotsync/internal/ctxkeys"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	CreateReservation(c *echo.Context) error
	GetMyReservations(c *echo.Context) error
	CancelReservation(c *echo.Context) error
	GetAllReservations(c *echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func getCurrentUserID(c *echo.Context) (uint, bool) {
	userId, ok := c.Get(string(ctxkeys.UserID)).(uint)
	return userId, ok
}

func (h *handler) CreateReservation(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return apperror.NewUnauthorized(nil, "Unauthorized")
	}

	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.CreateReservation(userId, req)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 201, "Reservation confirmed successfully", response)
}

func (h *handler) GetMyReservations(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return apperror.NewUnauthorized(nil, "Unauthorized")
	}

	reservations, err := h.service.GetMyReservations(userId)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "My reservations retrieved successfully", reservations)
}

func (h *handler) CancelReservation(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return apperror.NewUnauthorized(nil, "Unauthorized")
	}

	reservationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.NewBadRequest(err, "Invalid reservation id")
	}

	if err := h.service.CancelReservation(userId, uint(reservationID)); err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "Reservation cancelled successfully", nil)
}

func (h *handler) GetAllReservations(c *echo.Context) error {
	reservations, err := h.service.GetAllReservations()
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "All reservations retrieved successfully", reservations)
}
