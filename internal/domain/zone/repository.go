package zone

import (
	"errors"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	Create(zone *ParkingZone) error
	GetAll() ([]*ParkingZone, error)
	GetByID(id uint) (*ParkingZone, error)
	Update(zone *ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAll() ([]*ParkingZone, error) {
	var zones []*ParkingZone

	err := r.db.Find(&zones).Error
	if err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *repository) GetByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone

	err := r.db.First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	return &zone, nil
}

func (r *repository) Update(zone *ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&ParkingZone{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrZoneNotFound
	}
	return nil
}

// CountActiveReservations counts the number of active reservations for a zone.
func (r *repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return count, err
}
