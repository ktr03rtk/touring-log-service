package usecase

import (
	"context"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/repository"
	"github.com/pkg/errors"
)

type PayloadStoreUsecase interface {
	Execute(context.Context, <-chan *model.Payload) error
}

type payloadStoreUsecase struct {
	payloadStoreRepository repository.PayloadStoreRepository
}

func NewPayloadStoreUsecase(pr repository.PayloadStoreRepository) PayloadStoreUsecase {
	return &payloadStoreUsecase{
		payloadStoreRepository: pr,
	}
}

func (pu *payloadStoreUsecase) Execute(ctx context.Context, ch <-chan *model.Payload) error {
	for {
		select {
		case <-ctx.Done():

			return nil
		case p := <-ch:
			if err := pu.payloadStoreRepository.Store(ctx, p); err != nil {
				return errors.Wrapf(err, "failed to execute payload store usecase")
			}
		}
	}
}
