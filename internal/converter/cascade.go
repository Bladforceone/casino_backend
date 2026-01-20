package converter

import (
	"casino_backend/internal/api/dto/cascade"
	"casino_backend/internal/model"
)

func ToCascadeSpin(req cascade.CascadeSpinRequest) model.CascadeSpin {
	return model.CascadeSpin{
		Bet: req.Bet,
	}
}

// Основной конвертер результата спина
func ToCascadeSpinResponse(resp model.CascadeSpinResult) cascade.CascadeSpinResponse {
	return cascade.CascadeSpinResponse{
		InitialBoard:     resp.InitialBoard,
		Board:            resp.Board,
		Cascades:         toCascadeSteps(resp.Cascades),
		TotalPayout:      resp.TotalPayout,
		Balance:          resp.Balance,
		ScatterCount:     resp.ScatterCount,
		AwardedFreeSpins: resp.AwardedFreeSpins,
		FreeSpinsLeft:    resp.FreeSpinsLeft,
		InFreeSpin:       resp.InFreeSpin,
	}
}

// Вспомогательные конвертеры
func toCascadeSteps(steps []model.CascadeStep) []cascade.CascadeStep {
	result := make([]cascade.CascadeStep, len(steps))
	for i, step := range steps {
		result[i] = cascade.CascadeStep{
			CascadeIndex: step.CascadeIndex,
			Clusters:     toClusterInfos(step.Clusters),
			NewSymbols:   toNewSymbols(step.NewSymbols),
		}
	}
	return result
}

func toClusterInfos(clusters []model.ClusterInfo) []cascade.ClusterInfo {
	result := make([]cascade.ClusterInfo, len(clusters))
	for i, cl := range clusters {
		result[i] = cascade.ClusterInfo{
			Symbol:     cl.Symbol,
			Cells:      toPositions(cl.Cells),
			Count:      cl.Count,
			Payout:     cl.Payout,
			Multiplier: cl.Multiplier,
		}
	}
	return result
}

func toPositions(positions []model.Position) []cascade.Position {
	result := make([]cascade.Position, len(positions))
	for i, p := range positions {
		result[i] = cascade.Position{
			Row: p.Row,
			Col: p.Col,
		}
	}
	return result
}

func toNewSymbols(newSyms []struct {
	model.Position
	Symbol int
}) []cascade.NewSymbol {
	result := make([]cascade.NewSymbol, len(newSyms))
	for i, ns := range newSyms {
		result[i] = cascade.NewSymbol{
			Position: cascade.Position{
				Row: ns.Row,
				Col: ns.Col,
			},
			Symbol: ns.Symbol,
		}
	}
	return result
}

// Общий ответ с балансом и фриспинами
func ToCascadeDataResponse(data model.CascadeData) cascade.CascadeDataResponse {
	return cascade.CascadeDataResponse{
		Balance:       data.Balance,
		FreeSpinsLeft: data.FreeSpinCount,
	}
}
