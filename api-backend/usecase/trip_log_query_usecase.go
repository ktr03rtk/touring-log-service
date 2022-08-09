package usecase

import (
	"context"
	"strconv"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type TripLogQueryUsecase interface {
	Execute(ctx context.Context, year, month, day int, unit string) ([]*model.WebClientTrip, error)
}

type tripLogQueryUsecase struct {
	athenaQueryAdapter repository.AthenaQueryRepository
}

func NewTripLogQueryUsecase(aa repository.AthenaQueryRepository) TripLogQueryUsecase {
	return &tripLogQueryUsecase{
		athenaQueryAdapter: aa,
	}
}

func (tu *tripLogQueryUsecase) Execute(ctx context.Context, year, month, day int, unit string) ([]*model.WebClientTrip, error) {
	// GPS lat, lon is stored, when mode is 2 or 3.
	// TODO: fetch data with unit key
	query := "SELECT lat, lon FROM %s where year='%02d' and month='%02d' and day='%02d' and (mode=2 or mode=3);"

	args := []interface{}{
		year, month, day,
	}

	res, err := tu.athenaQueryAdapter.Fetch(ctx, query, args)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute trip log query usecase")
	}

	result := make([]*model.WebClientTrip, 0, len(res))

	for _, val := range res {
		lat, err := strconv.ParseFloat(val[0], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse lat")
		}

		lng, err := strconv.ParseFloat(val[1], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse lng")
		}

		array := &model.WebClientTrip{
			Lat: lat,
			Lng: lng,
		}

		result = append(result, array)
	}

	return result, nil
}
