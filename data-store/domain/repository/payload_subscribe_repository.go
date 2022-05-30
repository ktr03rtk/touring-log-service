//go:generate mockgen -source=payload_subscribe_repository.go -destination=../../mock/mock_payload_subscribe_repository.go -package=mock
package repository

import (
	"context"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
)

type PayloadSubscribeRepository interface {
	Subscribe(context.Context, <-chan *model.Payload) error
}
