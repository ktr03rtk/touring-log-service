package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ktr03rtk/touring-log-service/api-backend/config"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/service"
	"github.com/ktr03rtk/touring-log-service/api-backend/infrastructure/persistence"
	"github.com/ktr03rtk/touring-log-service/api-backend/interfaces/handler"
	"github.com/ktr03rtk/touring-log-service/api-backend/usecase"
)

var (
	region    string
	bucket    string
	jwtSecret string
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

	j, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return errors.New("env JWT_SECRET is not found")
	}

	jwtSecret = j

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn := config.NewDBConn()
	userRepository := persistence.NewUserPersistence(conn)
	photoMetadataRepository := persistence.NewPhotoMetadataPersistence(conn)

	photoImageRepository, err := persistence.NewPhotoImagePersistence(ctx, region, bucket)
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewUService(userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, userService)
	photoUsecase := usecase.NewPhotoStoreUsecase(photoMetadataRepository, photoImageRepository)

	h := handler.NewHandler(jwtSecret, userUsecase, photoUsecase)

	go func() {
		h.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Println("Caught SIGTERM, shutting down")

	h.Stop(ctx)
}
