//go:generate mockgen -source=photo_image_repository.go -destination=../../mock/mock_photo_image_repository.go -package=mock
package repository

import (
	"io"
)

type PhotoImageRepository interface {
	Store(io.Reader, string) error
	Get(key string) (io.ReadCloser, error)
}
