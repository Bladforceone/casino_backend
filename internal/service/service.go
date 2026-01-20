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
