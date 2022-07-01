package persistence

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PhotoMetadataPersistence struct {
	conn *gorm.DB
}

func NewPhotoMetadataPersistence(conn *gorm.DB) repository.PhotoMetadataRepository {
	return &PhotoMetadataPersistence{
		conn,
	}
}

func (mp *PhotoMetadataPersistence) Create(photo *model.Photo) error {
	if err := mp.conn.Create(&photo).Error; err != nil {
		return errors.Wrapf(err, "failed to create photo. photo: %+v", &photo)
	}

	return nil
}

// func (mp *PhotoMetadataPersistence) FindByID(id model.PhotoID) (*model.Photo, error) {
// 	t := &model.Photo{ID: id}

// 	if err := mp.conn.First(&t).Error; err != nil {
// 		return nil, errors.Wrapf(err, "failed to find photo. id: %+v", id)
// 	}

// 	return t, nil
// }

// func (mp *PhotoMetadataPersistence) FindAll() ([]*model.Photo, error) {
// 	var tasks []*model.Photo
// 	if err := mp.conn.Find(&tasks).Error; err != nil {
// 		return nil, errors.Wrapf(err, "failed to find all tasks")
// 	}

// 	return tasks, nil
// }

// func (mp *PhotoMetadataPersistence) Update(t *model.Photo) error {
// 	return mp.conn.Save(&t).Error
// }
