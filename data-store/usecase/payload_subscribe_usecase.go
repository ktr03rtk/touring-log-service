package usecase

import (
	"context"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/repository"
	"github.com/pkg/errors"
)

const (
	concurrency = 100
)

type PayloadSubscribeUsecase interface {
	Execute(context.Context, chan<- *model.Payload) error
}

type payloadSubscribeUsecase struct {
	payloadSubscribeRepository repository.PayloadSubscribeRepository
}

func NewPayloadSubscribeUsecase(pr repository.PayloadSubscribeRepository) PayloadSubscribeUsecase {
	return &payloadSubscribeUsecase{
		payloadSubscribeRepository: pr,
	}
}

func (pu *payloadSubscribeUsecase) Execute(ctx context.Context, ch chan<- *model.Payload) error {
	if err := pu.payloadSubscribeRepository.Subscribe(ctx, ch); err != nil {
		return errors.Wrapf(err, "failed to execute payload subscribe usecase")
	}

	return nil
}
