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
	region   string
	bucket   string
	endpoint string
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

	e, ok := os.LookupEnv("API_ENDPOINT")
	if !ok {
		return errors.New("env API_ENDPOINT is not found")
	}

	endpoint = e

	return nil
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	sbr, err := adapter.NewMqttAdapter(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	str, err := persistence.NewS3Persistence(ctx, region, bucket)
	if err != nil {
		logger.Fatal(err)
	}

	tms := adapter.NewHTTPAdapter(endpoint)

	sbu := usecase.NewPayloadSubscribeUsecase(sbr)
	stu := usecase.NewPayloadStoreUsecase(str, tms)

	h := handler.NewPayloadHandler(sbu, stu, logger)

	if err := h.Handle(ctx); err != nil {
		logger.Fatal(err)
	}
}
