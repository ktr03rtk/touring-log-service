package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/repository"
	"github.com/pkg/errors"
)

type httpAdapter struct {
	endpoint string
	*http.Client
}

func NewHTTPAdapter(endpoint string) repository.TripMetadataStoreRepository {
	return &httpAdapter{
		endpoint,
		&http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *httpAdapter) Create(trip *model.Trip) error {
	jsonString, err := json.Marshal(trip)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal trip: %+v", trip)
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewBuffer(jsonString))
	if err != nil {
		return errors.Wrapf(err, "failed to create request. endpoint: %+v, json: %+v", h.endpoint, jsonString)
	}

	req.Header.Set("Content-Type", "application/json")

	if _, err := h.Client.Do(req); err != nil {
		return errors.Wrapf(err, "failed to request")
	}

	return nil
}
