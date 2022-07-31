package usecase

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type DateListQueryUsecase interface {
	Execute(year, month int, user_id, unit string) ([]*model.LogDate, error)
}

type dateListQueryUsecase struct {
	queryAdapter repository.QueryRepository
}

func NewDateListQueryUsecase(qa repository.QueryRepository) DateListQueryUsecase {
	return &dateListQueryUsecase{
		queryAdapter: qa,
	}
}

func (lu *dateListQueryUsecase) Execute(year, month int, user_id, unit string) ([]*model.LogDate, error) {
	query := "SELECT year, month, day FROM photos WHERE year = ? AND month = ? AND user_id = ? UNION SELECT year, month, day FROM trips WHERE year = ? AND month = ? AND unit = ?"

	args := []interface{}{
		year, month, user_id, year, month, unit,
	}

	var touringLog []*model.LogDate

	res, err := lu.queryAdapter.Fetch(query, args, touringLog)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute list query usecase")
	}

	result, ok := res.([]*model.LogDate)
	if !ok {
		return nil, errors.New("failed to assertion")
	}

	return result, nil
}
