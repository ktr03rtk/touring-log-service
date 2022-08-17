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
	config               config
	userUsecase          usecase.UserUsecase
	photoStoreUsecase    usecase.PhotoStoreUsecase
	photoGetUsecase      usecase.PhotoGetUsecase
	tripUsecase          usecase.TripStoreUsecase
	listQueryUsecase     usecase.DateListQueryUsecase
	photoLogQueryUsecase usecase.PhotoLogQueryUsecase
	tripLogQueryUsecase  usecase.TripLogQueryUsecase
	logger               *log.Logger
	server               *http.Server
}

func NewHandler(secret string, uu usecase.UserUsecase, psu usecase.PhotoStoreUsecase, pgu usecase.PhotoGetUsecase, tu usecase.TripStoreUsecase, du usecase.DateListQueryUsecase, plu usecase.PhotoLogQueryUsecase, tlu usecase.TripLogQueryUsecase, logger *log.Logger) Handler {
	var cfg config
	cfg.jwt.secret = secret

	h := &handler{
		config:               cfg,
		userUsecase:          uu,
		photoStoreUsecase:    psu,
		photoGetUsecase:      pgu,
		tripUsecase:          tu,
		listQueryUsecase:     du,
		photoLogQueryUsecase: plu,
		tripLogQueryUsecase:  tlu,
		logger:               logger,
	}

	h.setupServer()

	return h
}

func (h *handler) Start() {
	if err := h.server.ListenAndServe(); err != nil {
		h.logger.Fatalln("Server closed with error:", err)
	}
}

func (h *handler) Stop(ctxParent context.Context) {
	ctx, cancel := context.WithTimeout(ctxParent, 10*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Println("Failed to gracefully shutdown:", err)
	}

	h.logger.Println("Server shutdown")
}

func (h *handler) setupServer() {
	router := httprouter.New()
	secure := alice.New(h.checkToken)

	router.HandlerFunc(http.MethodPost, "/v1/signup", h.signup)
	router.HandlerFunc(http.MethodPost, "/v1/login", h.login)

	router.POST("/v1/photos", h.wrap(secure.ThenFunc(h.storePhoto)))
	router.HandlerFunc(http.MethodGet, "/v1/photos/:id", h.getPhoto)

	router.HandlerFunc(http.MethodPost, "/v1/trips", h.storeTrip)

	router.POST("/v1/graphql", h.wrap(secure.ThenFunc(h.graphQL)))

	h.server = &http.Server{
		Handler:      h.enableCORS(router),
		Addr:         ":8080",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
