package reservation

import (
	"errors"
	"spotsync/internal/domain/zone"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrZoneFull            = errors.New("parking zone is full")
)

type Repository interface {
	CreateWithCapacityCheck(reservation *Reservation) error
	GetByID(id uint) (*Reservation, error)
	GetByUserID(userID uint) ([]*Reservation, error)
	GetAll() ([]*Reservation, error)
	Update(reservation *Reservation) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateWithCapacityCheck(reservation *Reservation) error {
	// Start transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		var z zone.ParkingZone

		// 1. Lock the row!
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&z, reservation.ZoneID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return zone.ErrZoneNotFound
			}
			return err
		}

		// 2. Count current 'active' reservations for this zone
		var activeCount int64
		err = tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, "active").
			Count(&activeCount).Error
		if err != nil {
			return err
		}

		// 3. Check if active_count < zone.total_capacity
		if activeCount >= int64(z.TotalCapacity) {
			return ErrZoneFull
		}

		// 4. Create reservation
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil // Commits transaction
	})
}

func (r *repository) GetByID(id uint) (*Reservation, error) {
	var reservation Reservation
	err := r.db.First(&reservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *repository) GetByUserID(userID uint) ([]*Reservation, error) {
	var reservations []*Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) GetAll() ([]*Reservation, error) {
	var reservations []*Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) Update(reservation *Reservation) error {
	return r.db.Save(reservation).Error
}
