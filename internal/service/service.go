package service

import (
	"casino_backend/internal/model"
	"context"
)

type LineService interface {
	Spin(ctx context.Context, spinReq model.LineSpin) (*model.SpinResult, error)
	BuyBonus(amount int) error
	Deposit(amount int) error
	CheckData() (*model.Data, error)
}

type CascadeService interface {
	Spin(ctx context.Context, req model.CascadeSpin) (*model.CascadeSpinResult, error)
	BuyBonus(amount int) error
	Deposit(amount int) error
	CheckData() (*model.CascadeData, error)
}

type AuthService interface {
	Register(ctx context.Context, name, login, password string) (accessToken string, sessionID string, err error)
	Login(ctx context.Context, login, password string) (accessToken string, sessionID string, err error)
	Refresh(ctx context.Context, sessionID string) (newAccessToken string, err error)
	Logout(ctx context.Context, sessionID string) error
}
