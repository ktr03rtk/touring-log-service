package handler

import (
	"context"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/usecase"
	"golang.org/x/sync/errgroup"
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

func (ph payloadHandler) Handle(ctxParent context.Context) error {
	eg, ctx := errgroup.WithContext(ctxParent)

	eg.Go(func() error { return ph.payloadSubscribeUsecase.Execute(ctx, ph.payloadCh) })
	eg.Go(func() error { return ph.payloadStoreUsecase.Execute(ctx, ph.payloadCh) })

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
