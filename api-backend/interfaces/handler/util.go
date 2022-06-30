package handler

import (
	"encoding/json"
	"net/http"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (h *handler) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		next.ServeHTTP(w, r)
	})
}

func (h *handler) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	wrapper[wrap] = data

	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (h *handler) errJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{Message: err.Error()}

	h.writeJSON(w, statusCode, theError, "error")
}
