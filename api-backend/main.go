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
	tripMetadataRepository := persistence.NewTripMetadataPersistence(conn)
	queryAdapterRepository := persistence.NewQueryAdapter(conn)

	photoImageRepository, err := persistence.NewPhotoImagePersistence(ctx, region, bucket)
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewUserService(userRepository)
	tripService := service.NewTripService(tripMetadataRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, userService)
	photoUsecase := usecase.NewPhotoStoreUsecase(photoMetadataRepository, photoImageRepository)
	tripUsecase := usecase.NewTripStoreUsecase(tripMetadataRepository, tripService)
	listQueryUsecase := usecase.NewDateListQueryUsecase(queryAdapterRepository)
	photoLogQueryUsecase := usecase.NewPhotoLogQueryUsecase(queryAdapterRepository)

	h := handler.NewHandler(jwtSecret, userUsecase, photoUsecase, tripUsecase, listQueryUsecase, photoLogQueryUsecase)

	go func() {
		h.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Println("Caught SIGTERM, shutting down")

	h.Stop(ctx)
}
