package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ktr03rtk/touring-log-service/data-store/infrastructure/adapter"
	"github.com/ktr03rtk/touring-log-service/data-store/infrastructure/persistence"
	"github.com/ktr03rtk/touring-log-service/data-store/interfaces/handler"
	"github.com/ktr03rtk/touring-log-service/data-store/usecase"
)

var (
	region string
	bucket string
)

func init() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
}

func getEnv() error {
	r, ok := os.LookupEnv("REGION")
	if !ok {
		return errors.New("env REGION is not found")
	}

	region = r

	b, ok := os.LookupEnv("BUCKET")
	if !ok {
		return errors.New("env BUCKET is not found")
	}

	bucket = b

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	sbr, err := adapter.NewMqttAdapter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	str, err := persistence.NewS3Persistence(ctx, region, bucket)
	if err != nil {
		log.Fatal(err)
	}

	sbu := usecase.NewPayloadSubscribeUsecase(sbr)
	stu := usecase.NewPayloadStoreUsecase(str)

	h := handler.NewPayloadHandler(sbu, stu)

	if err := h.Handle(ctx); err != nil {
		log.Fatal(err)
	}
}
