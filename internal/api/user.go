package api

import (
	"casino_backend/internal/api/dto"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandlerDeps struct{}

type UserHandler struct {
	deps UserHandlerDeps
}

func NewUserHandler(r chi.Router, deps UserHandlerDeps) {
	h := &UserHandler{deps: deps}
	r.Route("/user", func(r chi.Router) {
		r.Post("/deposit", h.Deposit)
	})
}

func (uh *UserHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req dto.DepositRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: реализовать реальные операции по пополнению баланса (сервис/репозиторий)
	resp := dto.DepositResponse{Status: 0, Msg: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
