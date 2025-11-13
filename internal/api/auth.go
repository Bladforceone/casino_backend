package api

import (
	"casino_backend/internal/api/dto"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AuthHandlerDeps struct{}

type AuthHandler struct {
	deps AuthHandlerDeps
}

func NewAuthHandler(r chi.Router, deps AuthHandlerDeps) {
	handler := &AuthHandler{deps: deps}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	var req dto.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
		return
	}

	// TODO: заменить заглушку реальной регистрацией пользователя и генерацией токена
	resp := dto.RegistrationResponse{Token: "mock-registration-token"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
		return
	}

	// TODO: заменить заглушку реальной аутентификацией и генерацией токена
	resp := dto.LoginResponse{Token: "mock-login-token"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
