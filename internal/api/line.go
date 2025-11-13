package api

import (
	"casino_backend/internal/api/dto"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LineHandlerDeps struct{}

type LineHandler struct {
	deps LineHandlerDeps
}

func NewLineHandler(r chi.Router, deps LineHandlerDeps) {
	handler := &LineHandler{deps: deps}
	r.Route("/line", func(r chi.Router) {
		r.Post("/spin", handler.Spin)
	})
}

func (lh *LineHandler) Spin(w http.ResponseWriter, r *http.Request) {
	data := dto.LineSpinRequest{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: реализовать бизнес-логику спина и формирование ответа
	w.WriteHeader(http.StatusOK)
	return
}
