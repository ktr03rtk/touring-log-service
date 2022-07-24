package persistence

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TripMetadataPersistence struct {
	conn *gorm.DB
}

func NewTripMetadataPersistence(conn *gorm.DB) repository.TripMetadataStoreRepository {
	return &TripMetadataPersistence{
		conn,
	}
}

func (tp *TripMetadataPersistence) Create(trip *model.Trip) error {
	if err := tp.conn.Create(&trip).Error; err != nil {
		return errors.Wrapf(err, "failed to create trip. trip: %+v", &trip)
	}

	return nil
}

func (tp *TripMetadataPersistence) FindByDateAndUnit(year, month, day int, unit string) (*model.Trip, error) {
	t := &model.Trip{
		Year:  year,
		Month: month,
		Day:   day,
		Unit:  unit,
	}

	if err := tp.conn.First(&t).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to find Trip")
	}

	return t, nil
}

// func (mp *TripMetadataPersistence) FindAll() ([]*model.Trip, error) {
// 	var tasks []*model.Trip
// 	if err := mp.conn.Find(&tasks).Error; err != nil {
// 		return nil, errors.Wrapf(err, "failed to find all tasks")
// 	}

// 	return tasks, nil
// }

// func (mp *TripMetadataPersistence) Update(t *model.Trip) error {
// 	return mp.conn.Save(&t).Error
// }
