package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Unit     string `json:"unit"`
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
