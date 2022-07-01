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

func (tp *PhotoMetadataPersistence) Create(photo *model.Photo) error {
	if err := tp.conn.Create(&photo).Error; err != nil {
		return errors.Wrapf(err, "failed to create photo. photo: %+v", &photo)
	}

	return nil
}

