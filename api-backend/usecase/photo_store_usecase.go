package usecase

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type PhotoStoreUsecase interface {
	Execute(file *multipart.FileHeader, userID, unit string) error
}

type photoStoreUsecase struct {
	photoMetadataRepository repository.PhotoMetadataRepository
	photoImageRepository    repository.PhotoImageRepository
}

func NewPhotoStoreUsecase(pmr repository.PhotoMetadataRepository, pir repository.PhotoImageRepository) PhotoStoreUsecase {
	return &photoStoreUsecase{
		photoMetadataRepository: pmr,
		photoImageRepository:    pir,
	}
}

func (pu *photoStoreUsecase) Execute(file *multipart.FileHeader, userID, unit string) error {
	src, err := file.Open()
	if err != nil {
		return errors.Wrapf(err, "failed to open file")
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	reader := io.TeeReader(src, buf)

	// consume reader to pass data to TeeReader
	s3ObjectKey, err := storeMetadata(reader, userID, unit, pu)
	if err != nil {
		return err
	}

	if err := pu.photoImageRepository.Store(buf, s3ObjectKey); err != nil {
		return err
	}

	return nil
}

func storeMetadata(reader io.Reader, userID string, unit string, pu *photoStoreUsecase) (s3ObjectKey string, err error) {
	lat, lon, time, err := model.ExtractMetadata(reader)
	if err != nil {
		return "", err
	}

	id := model.CreateUUID()

	metadata, err := model.NewPhoto(model.PhotoID(id), time, lat, lon, model.UserID(userID), unit)
	if err != nil {
		return "", err
	}

	if err := pu.photoMetadataRepository.Create(metadata); err != nil {
		return "", errors.Wrapf(err, "failed to execute photo store usecase")
	}

	if _, err := io.ReadAll(reader); err != nil {
		return "", errors.Wrapf(err, "failed to read remaining photo data")
	}

	return metadata.S3ObjectKey, nil
}
