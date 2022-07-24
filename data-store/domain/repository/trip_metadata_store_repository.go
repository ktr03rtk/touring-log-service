//go:generate mockgen -source=trip_metadata_repository.go -destination=../../mock/mock_trip_metadata_repository.go -package=mock
package repository

import "github.com/ktr03rtk/touring-log-service/data-store/domain/model"

type TripMetadataStoreRepository interface {
	Create(*model.Trip) error
}
