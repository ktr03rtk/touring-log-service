package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
)

type UserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Unit     string `json:"unit"`
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.errJSON(w, err)

		return
	}

	if err := h.userUsecase.SignUp(payload.Email, payload.Password, payload.Unit); err != nil {
		h.errJSON(w, err)

		return
	}

	ok := jsonResp{
		OK: true,
	}

	if err := h.writeJSON(w, http.StatusOK, ok, "response"); err != nil {
		h.errJSON(w, err)
		return
	}
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.errJSON(w, errors.New("unauthorized"))

		return
	}

	user, err := h.userUsecase.Authenticate(payload.Email, payload.Password)
	if err != nil {
		h.errJSON(w, err)

		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(user.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	// TODO: fix
	claims.Issuer = "my.example.com"
	claims.Audiences = []string{"my.example.com"}
	claims.Set = map[string]interface{}{"unit": user.Unit}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(h.config.jwt.secret))
	if err != nil {
		h.errJSON(w, errors.New("error signing"))

		return
	}

	h.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}
