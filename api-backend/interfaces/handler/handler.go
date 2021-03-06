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
	Stop(context.Context)
}

type config struct {
	jwt struct {
		secret string
	}
}

type handler struct {
	config       config
	userUsecase  usecase.UserUsecase
	photoUsecase usecase.PhotoStoreUsecase
	server       *http.Server
}

func NewHandler(secret string, uu usecase.UserUsecase, pu usecase.PhotoStoreUsecase) Handler {
	var cfg config
	cfg.jwt.secret = secret

	h := &handler{
		config:       cfg,
		userUsecase:  uu,
		photoUsecase: pu,
	}

	h.setupServer()

	return h
}

func (h *handler) Start() {
	if err := h.server.ListenAndServe(); err != nil {
		log.Fatalln("Server closed with error:", err)
	}
}

func (h *handler) Stop(ctxParent context.Context) {
	ctx, cancel := context.WithTimeout(ctxParent, 10*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Println("Failed to gracefully shutdown:", err)
	}

	log.Println("Server shutdown")
}

func (h *handler) setupServer() {
	router := httprouter.New()
	secure := alice.New(h.checkToken)

	router.HandlerFunc(http.MethodPost, "/v1/signup", h.signup)
	router.HandlerFunc(http.MethodPost, "/v1/login", h.login)

	router.POST("/v1/upload", h.wrap(secure.ThenFunc(h.upload)))

	h.server = &http.Server{
		Handler: h.enableCORS(router),
		Addr:    ":8080",
	}
}
