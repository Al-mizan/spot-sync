package zone

import (
	"spotsync/internal/apperror"
	"spotsync/internal/domain/zone/dto"
	"spotsync/internal/httpresponse"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	CreateZone(c *echo.Context) error
	GetAllZones(c *echo.Context) error
	GetZoneByID(c *echo.Context) error
	UpdateZone(c *echo.Context) error
	DeleteZone(c *echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &handler{service: s}
}

func (h *handler) CreateZone(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.CreateZone(req)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 201, "Parking zone created successfully", response)
}

func (h *handler) GetAllZones(c *echo.Context) error {
	zones, err := h.service.GetAllZones()
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "Parking zones retrieved successfully", zones)
}

func (h *handler) GetZoneByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.NewBadRequest(err, "Invalid zone id")
	}

	response, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "Parking zone retrieved successfully", response)
}

func (h *handler) UpdateZone(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.NewBadRequest(err, "Invalid zone id")
	}

	var req dto.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.UpdateZone(uint(id), req)
	if err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "Parking zone updated successfully", response)
}

func (h *handler) DeleteZone(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return apperror.NewBadRequest(err, "Invalid zone id")
	}

	if err := h.service.DeleteZone(uint(id)); err != nil {
		return err
	}

	return httpresponse.Success(c, 200, "Parking zone deleted successfully", nil)
}
