package usecase

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type ListQueryUsecase interface {
	Execute(year, month int, user_id, unit string) ([]*model.TouringLog, error)
}

type listQueryUsecase struct {
	queryAdapter repository.QueryRepository
}

func NewListQueryUsecase(qa repository.QueryRepository) ListQueryUsecase {
	return &listQueryUsecase{
		queryAdapter: qa,
	}
}

func (lu *listQueryUsecase) Execute(year, month int, user_id, unit string) ([]*model.TouringLog, error) {
	query := "SELECT year, month, day FROM photos WHERE year = ? AND month = ? AND user_id = ? UNION SELECT year, month, day FROM trips WHERE year = ? AND month = ? AND unit = ?"

	args := []interface{}{
		year, month, user_id, year, month, unit,
	}

	var touringLog []*model.TouringLog

	res, err := lu.queryAdapter.Fetch(query, args, touringLog)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute list query usecase")
	}

	result, ok := res.([]*model.TouringLog)
	if !ok {
		return nil, errors.New("failed to assertion")
	}

	return result, nil
}
