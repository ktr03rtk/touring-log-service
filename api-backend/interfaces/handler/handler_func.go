package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func (h *handler) storePhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
		h.errJSON(w, errors.Wrapf(err, "failed to parse multipartform photo data"))

		return
	}

	id, unit, err := getAccountInfo(r)
	if err != nil {
		h.errJSON(w, err)

		return
	}

	files := r.MultipartForm.File["images"]
	for _, file := range files {
		if err := h.photoUsecase.Execute(file, id, unit); err != nil {
			h.errJSON(w, err)

			return
		}
	}

	res := jsonResp{
		OK: true,
	}

	if err := h.writeJSON(w, http.StatusOK, res, "response"); err != nil {
		h.errJSON(w, err)
		return
	}
}

type TripPayload struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
	Unit  string `json:"unit"`
}

func (h *handler) storeTrip(w http.ResponseWriter, r *http.Request) {
	var payload TripPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.errJSON(w, errors.Wrapf(err, "failed to json decode trip"))

		return
	}

	if err := h.tripUsecase.Execute(payload.Year, payload.Month, payload.Day, payload.Unit); err != nil {
		h.errJSON(w, err)

		return
	}

	res := jsonResp{
		OK: true,
	}

	if err := h.writeJSON(w, http.StatusOK, res, "response"); err != nil {
		h.errJSON(w, err)
		return
	}
}
