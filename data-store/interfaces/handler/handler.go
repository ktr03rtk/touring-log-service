package handler

import (
	"context"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/usecase"
)

const (
	concurrency = 100
)

type PayloadHandler interface {
	Handle(ctx context.Context) error
}

type payloadHandler struct {
	payloadSubscribeUsecase usecase.PayloadSubscribeUsecase
	payloadStoreUsecase     usecase.PayloadStoreUsecase
	payloadCh               chan *model.Payload
}

func NewPayloadHandler(sbu usecase.PayloadSubscribeUsecase, stu usecase.PayloadStoreUsecase) PayloadHandler {
	ch := make(chan *model.Payload, concurrency)

	return &payloadHandler{
		payloadSubscribeUsecase: sbu,
		payloadStoreUsecase:     stu,
		payloadCh:               ch,
	}
}

func (ph payloadHandler) Handle(ctx context.Context) error {
	if err := ph.payloadSubscribeUsecase.Execute(ctx, ph.payloadCh); err != nil {
		return err
	}

	if err := ph.payloadStoreUsecase.Execute(ctx, ph.payloadCh); err != nil {
		return err
	}

	return nil
}
