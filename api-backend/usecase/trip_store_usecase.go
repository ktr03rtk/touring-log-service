package usecase

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/service"
	"github.com/pkg/errors"
)

type TripStoreUsecase interface {
	Execute(year, month, day int, unit string) error
}

type tripStoreUsecase struct {
	tripMetadataRepository repository.TripMetadataStoreRepository
	tripService            service.TripService
}

func NewTripStoreUsecase(tmr repository.TripMetadataStoreRepository, ts service.TripService) TripStoreUsecase {
	return &tripStoreUsecase{
		tripMetadataRepository: tmr,
		tripService:            ts,
	}
}

func (tu *tripStoreUsecase) Execute(year, month, day int, unit string) error {
	ok, err := tu.tripService.IsExists(year, month, day, unit)
	if ok {
		return nil
	} else if err != nil {
		return err
	}

	id := model.CreateUUID()

	metadata := model.NewTrip(model.TripID(id), year, month, day, unit)

	if err := tu.tripMetadataRepository.Create(metadata); err != nil {
		return errors.Wrapf(err, "failed to execute trip store usecase")
	}

	return nil
}
