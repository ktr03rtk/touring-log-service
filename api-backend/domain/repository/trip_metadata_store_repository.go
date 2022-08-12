//go:generate mockgen -source=trip_metadata_store_repository.go -destination=../../mock/mock_trip_metadata_store_repository.go -package=mock
package repository

import "github.com/ktr03rtk/touring-log-service/api-backend/domain/model"

type TripMetadataStoreRepository interface {
	Create(*model.Trip) error
	FindByDateAndUnit(year, month, day int, unit string) (*model.Trip, error)
}
