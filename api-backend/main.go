package main

import (
	"fmt"
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

func main() {
	conn := config.NewDBConn()
	userRepository := persistence.NewUserPersistence(conn)
	userService := service.NewUService(userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, userService)

	h := handler.NewHandler(userUsecase)

	go func() {
		h.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Println("Caught SIGTERM, shutting down")

	h.Stop()
}
