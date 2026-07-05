package zone

import (
	"errors"
	"spotsync/internal/apperror"
	"spotsync/internal/domain/zone/dto"
)

type Service interface {
	CreateZone(req dto.CreateRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]*dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req dto.UpdateRequest) (*dto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateZone(req dto.CreateRequest) (*dto.ZoneResponse, error) {
	zone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(&zone); err != nil {
		return nil, apperror.NewInternal(err, "failed to create parking zone")
	}

	return s.toResponse(&zone, 0), nil
}

func (s *service) GetAllZones() ([]*dto.ZoneResponse, error) {
	zones, err := s.repo.GetAll()
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch parking zones")
	}

	responses := make([]*dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		activeCount, err := s.repo.CountActiveReservations(z.ID)
		if err != nil {
			return nil, apperror.NewInternal(err, "failed to count active reservations")
		}
		responses = append(responses, s.toResponse(z, activeCount))
	}

	return responses, nil
}

func (s *service) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return nil, apperror.NewNotFound(err, "parking zone not found")
		}
		return nil, apperror.NewInternal(err, "failed to fetch parking zone")
	}

	activeCount, err := s.repo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to count active reservations")
	}

	return s.toResponse(zone, activeCount), nil
}

func (s *service) UpdateZone(id uint, req dto.UpdateRequest) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return nil, apperror.NewNotFound(err, "parking zone not found")
		}
		return nil, apperror.NewInternal(err, "failed to fetch parking zone")
	}

	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity != 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour != 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.repo.Update(zone); err != nil {
		return nil, apperror.NewInternal(err, "failed to update parking zone")
	}

	activeCount, err := s.repo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to count active reservations")
	}

	return s.toResponse(zone, activeCount), nil
}

func (s *service) DeleteZone(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return apperror.NewNotFound(err, "parking zone not found")
		}
		return apperror.NewInternal(err, "failed to delete parking zone")
	}
	return nil
}

// toResponse converts a ParkingZone entity to a ZoneResponse DTO.
func (s *service) toResponse(zone *ParkingZone, activeCount int64) *dto.ZoneResponse {
	availableSpots := zone.TotalCapacity - int(activeCount)
	if availableSpots < 0 {
		availableSpots = 0
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}
}
