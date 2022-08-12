package usecase

import (
	"io"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type PhotoGetUsecase interface {
	Execute(photoID string) (io.ReadCloser, error)
}

type photoGetUsecase struct {
	queryAdapter         repository.QueryRepository
	photoImageRepository repository.PhotoImageRepository
}

func NewPhotoGetUsecase(qa repository.QueryRepository, pir repository.PhotoImageRepository) PhotoGetUsecase {
	return &photoGetUsecase{
		queryAdapter:         qa,
		photoImageRepository: pir,
	}
}

type photoKey struct {
	S3ObjectKey string `json:"s3_object_key"`
}

func (pu *photoGetUsecase) Execute(photoID string) (io.ReadCloser, error) {
	query := "SELECT * FROM photos WHERE id = ?"

	args := []interface{}{
		photoID,
	}

	var photoMetadata []*photoKey

	res, err := pu.queryAdapter.Fetch(query, args, photoMetadata)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute photo get query usecase")
	}

	metadata, ok := res.([]*photoKey)
	if !ok {
		return nil, errors.New("failed to assertion")
	}

	r, err := pu.photoImageRepository.Get(metadata[0].S3ObjectKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get object")
	}

	return r, nil
}
