package usecase

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type PhotoLogQueryUsecase interface {
	Execute(year, month, day int, user_id string) ([]*model.ClientPhoto, error)
}

type photoLogQueryUsecase struct {
	queryAdapter repository.QueryRepository
}

func NewPhotoLogQueryUsecase(qa repository.QueryRepository) PhotoLogQueryUsecase {
	return &photoLogQueryUsecase{
		queryAdapter: qa,
	}
}

func (lu *photoLogQueryUsecase) Execute(year, month, day int, user_id string) ([]*model.ClientPhoto, error) {
	query := "SELECT id, lat, lon as lng FROM photos WHERE year = ? AND month = ? AND day = ? AND user_id = ?"

	args := []interface{}{
		year, month, day, user_id,
	}

	var touringLog []*model.ClientPhoto

	res, err := lu.queryAdapter.Fetch(query, args, touringLog)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute photo log query usecase")
	}

	result, ok := res.([]*model.ClientPhoto)
	if !ok {
		return nil, errors.New("failed to assertion")
	}

	return result, nil
}
