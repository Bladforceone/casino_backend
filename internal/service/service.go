package service

import "casino_backend/internal/model"

type SlotsService interface {
	Spin() (model.SpinResult, error)
}
