package converter

import (
	dto "casino_backend/internal/api/dto/auth"
	"casino_backend/internal/model"
)

func RegisterRequestToUserModel(req *dto.RegisterRequest) *model.User {
	return &model.User{
		Name:     req.Name,
		Login:    req.Login,
		Password: req.Password,
	}
}

func LoginRequestToUserModel(req *dto.LoginRequest) *model.User {
	return &model.User{
		Login:    req.Login,
		Password: req.Password,
	}
}
