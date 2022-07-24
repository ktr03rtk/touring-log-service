package service

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"

	"github.com/pkg/errors"
)

type TripService interface {
	IsExists(year, month, day int, unit string) (bool, error)
}

type tripService struct {
	tripMetadataRepository repository.TripMetadataStoreRepository
}

func NewTripService(tr repository.TripMetadataStoreRepository) TripService {
	return &tripService{tripMetadataRepository: tr}
}

func (s *tripService) IsExists(year, month, day int, unit string) (bool, error) {
	t, err := s.tripMetadataRepository.FindByDateAndUnit(year, month, day, unit)
	if err != nil {
		return false, errors.Wrapf(err, "failed to find trip")
	} else if t == nil {
		return false, nil
	}

	return true, nil
}
