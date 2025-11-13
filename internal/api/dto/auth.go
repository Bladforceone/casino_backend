package dto

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
