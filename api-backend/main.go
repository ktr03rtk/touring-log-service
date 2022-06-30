package main

import (
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

var jwtSecret string

func init() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
}

func getEnv() error {
	j, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return errors.New("env JWT_SECRET is not found")
	}

	jwtSecret = j

	return nil
}

func main() {
	conn := config.NewDBConn()
	userRepository := persistence.NewUserPersistence(conn)
	userService := service.NewUService(userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, userService)

	h := handler.NewHandler(jwtSecret, userUsecase)

	go func() {
		h.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Println("Caught SIGTERM, shutting down")

	h.Stop()
}
