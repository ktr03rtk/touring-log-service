//go:generate mockgen -source=photo_metadata_repository.go -destination=../../mock/mock_photo_metadata_repository.go -package=mock
package repository

import "github.com/ktr03rtk/touring-log-service/api-backend/domain/model"

type PhotoMetadataRepository interface {
	Create(*model.Photo) error
}
