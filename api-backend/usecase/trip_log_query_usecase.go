package usecase

import (
	"context"
	"strconv"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type TripLogQueryUsecase interface {
	Execute(ctx context.Context, year, month, day int, unit string) ([]*model.WebClientTrip, *model.WebClientTrip, error)
}

type tripLogQueryUsecase struct {
	athenaQueryAdapter repository.AthenaQueryRepository
}

func NewTripLogQueryUsecase(aa repository.AthenaQueryRepository) TripLogQueryUsecase {
	return &tripLogQueryUsecase{
		athenaQueryAdapter: aa,
	}
}

type latLngRange struct {
	min float64
	max float64
}

func (tu *tripLogQueryUsecase) Execute(ctx context.Context, year, month, day int, unit string) ([]*model.WebClientTrip, *model.WebClientTrip, error) {
	// GPS lat, lon is stored, when mode is 2 or 3.
	// TODO: fetch data with unit key
	query := "SELECT lat, lon FROM %s where year='%02d' and month='%02d' and day='%02d' and (mode=2 or mode=3);"

	args := []interface{}{
		year, month, day,
	}

	res, err := tu.athenaQueryAdapter.Fetch(ctx, query, args)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to execute trip log query usecase")
	}

	if len(res) == 0 {
		return nil, nil, nil
	}

	trip := make([]*model.WebClientTrip, 0, len(res))
	latRange := &latLngRange{min: 180, max: 0}
	lngRange := &latLngRange{min: 180, max: 0}

	for _, val := range res {
		lat, err := strconv.ParseFloat(val[0], 64)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to parse lat")
		}

		lng, err := strconv.ParseFloat(val[1], 64)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to parse lng")
		}

		array := &model.WebClientTrip{
			Lat: lat,
			Lng: lng,
		}

		if lat < latRange.min {
			latRange.min = lat
		} else if lat > latRange.max {
			latRange.max = lat
		}

		if lng < lngRange.min {
			lngRange.min = lng
		} else if lng > lngRange.max {
			lngRange.max = lng
		}

		trip = append(trip, array)
	}

	center := &model.WebClientTrip{}
	center.Lat = (latRange.min + latRange.max) / 2
	center.Lng = (lngRange.min + lngRange.max) / 2

	return trip, center, nil
}
