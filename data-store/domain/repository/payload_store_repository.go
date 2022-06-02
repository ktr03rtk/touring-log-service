//go:generate mockgen -source=payload_store_repository.go -destination=../../mock/mock_payload_store_repository.go -package=mock
package repository

import (
	"context"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
)

type PayloadStoreRepository interface {
	Store(context.Context, *model.Payload) error
}
