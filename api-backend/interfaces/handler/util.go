package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

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

		if claims.Subject == "" {
			h.errJSON(w, errors.New("unauthorized"), http.StatusForbidden)

			return
		}

		ctx := context.WithValue(r.Context(), "id", claims.Subject)
		ctx = context.WithValue(ctx, "unit", claims.Set["unit"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getAccountInfo(r *http.Request) (string, string, error) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		return "", "", errors.New("failed to get context id")
	}

	unit, ok := r.Context().Value("unit").(string)
	if !ok {
		return "", "", errors.New("failed to get context unit")
	}

	return id, unit, nil
}
