package service

import (
	"casino_backend/internal/model"
	"context"
)

type SlotsService interface {
	Spin(ctx context.Context, spinReq model.LineSpin) (model.SpinResult, error)
}
