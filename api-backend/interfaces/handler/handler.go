package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ktr03rtk/touring-log-service/api-backend/usecase"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Handler interface {
	Start()
	Stop()
}

type config struct {
	jwt struct {
		secret string
	}
}

type handler struct {
	config config
	userUsecase usecase.UserUsecase
	server *http.Server
}

func NewHandler(secret string, uu usecase.UserUsecase) Handler {
	var cfg config
	cfg.jwt.secret = secret

	h := &handler{
		config: cfg,
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
	_ = alice.New(h.checkToken)

	router.HandlerFunc(http.MethodPost, "/v1/signup", h.signup)
	router.HandlerFunc(http.MethodPost, "/v1/login", h.login)


	h.server = &http.Server{
		Handler: h.enableCORS(router),
		Addr:    ":8080",
	}
}
