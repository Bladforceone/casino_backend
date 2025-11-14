package service

import "casino_backend/internal/model"

type SlotsService interface {
	Spin(spin model.LineSpin) (model.SpinResult, error)
}
