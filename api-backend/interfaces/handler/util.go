package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
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

func (h *handler) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			h.errJSON(w, errors.New("invalid auth header"))

			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			h.errJSON(w, errors.New("invalid auth header"))

			return
		}

		if headerParts[0] != "Bearer" {
			h.errJSON(w, errors.New("unauthorized - no bearer"))

			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(h.config.jwt.secret))
		if err != nil {
			h.errJSON(w, errors.New("unauthorized - failed hmac check"), http.StatusForbidden)

			return
		}

		if !claims.Valid(time.Now()) {
			h.errJSON(w, errors.New("unauthorized - token expired"), http.StatusForbidden)

			return
		}

		// TODO: fix
		if !claims.AcceptAudience("my.example.com") {
			h.errJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)

			return
		}

		if _, err := strconv.ParseInt(claims.Subject, 10, 64); err != nil {
			h.errJSON(w, errors.New("unauthorized"), http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}
