package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ktr03rtk/touring-log-service/api-backend/usecase"

	"github.com/julienschmidt/httprouter"
)

type Handler interface {
	Start()
	Stop()
}

type handler struct {
	userUsecase usecase.UserUsecase
	server *http.Server
}

func NewHandler(uu usecase.UserUsecase) Handler {
	h := &handler{
		userUsecase: uu,
	}

	h.setupServer()

	return h
}

func (h *handler) Start() {
	if err := h.server.ListenAndServe(); err != nil {
		log.Fatalln("Server closed with error:", err)
	}
}

func (h *handler) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Println("Failed to gracefully shutdown:", err)
	}

	log.Println("Server shutdown")
}

func (h *handler) setupServer() {
	router := httprouter.New()

	router.POST("/v1/signup", h.signup)


	h.server = &http.Server{
		Handler: h.enableCORS(router),
		Addr:    ":8080",
	}
}
