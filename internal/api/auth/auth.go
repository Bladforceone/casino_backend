package auth

import "net/http"

type AuthHandlerDeps struct {
}

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

}
