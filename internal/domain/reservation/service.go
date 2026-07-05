package reservation

import (
	"errors"
	"spotsync/internal/apperror"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/domain/zone"
)

type Service interface {
	CreateReservation(userID uint, req dto.CreateRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]*dto.MyReservationResponse, error)
	CancelReservation(userID uint, reservationID uint) error
	GetAllReservations() ([]*dto.AdminReservationResponse, error)
}

type service struct {
	reservationRepo Repository
	zoneRepo        zone.Repository
}

func NewService(reservationRepo Repository, zoneRepo zone.Repository) Service {
	return &service{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *service) CreateReservation(userID uint, req dto.CreateRequest) (*dto.ReservationResponse, error) {
	reservation := &Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	err := s.reservationRepo.CreateWithCapacityCheck(reservation)
	if err != nil {
		if errors.Is(err, zone.ErrZoneNotFound) {
			return nil, apperror.NewNotFound(err, "parking zone not found")
		}
		if errors.Is(err, ErrZoneFull) {
			return nil, apperror.NewConflict(err, "parking zone is full")
		}
		return nil, apperror.NewInternal(err, "failed to create reservation")
	}

	response := &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}

	return response, nil
}

func (s *service) GetMyReservations(userID uint) ([]*dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.GetByUserID(userID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch reservations")
	}

	responses := make([]*dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, &dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneMiniResponse{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, nil
}

func (s *service) CancelReservation(userID uint, reservationID uint) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		if errors.Is(err, ErrReservationNotFound) {
			return apperror.NewNotFound(err, "reservation not found")
		}
		return apperror.NewInternal(err, "failed to fetch reservation")
	}

	if reservation.UserID != userID {
		return apperror.NewForbidden(nil, "you do not own this reservation")
	}

	if reservation.Status == "cancelled" {
		// already cancelled
		return nil
	}

	reservation.Status = "cancelled"
	if err := s.reservationRepo.Update(reservation); err != nil {
		return apperror.NewInternal(err, "failed to cancel reservation")
	}

	return nil
}

func (s *service) GetAllReservations() ([]*dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.GetAll()
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch reservations")
	}

	responses := make([]*dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, &dto.AdminReservationResponse{
			ID: r.ID,
			User: dto.UserMiniResponse{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneMiniResponse{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}

	return responses, nil
}
