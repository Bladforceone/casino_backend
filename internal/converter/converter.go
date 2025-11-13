package converter

import (
	"casino_backend/internal/api/dto"
	"casino_backend/internal/model"
)

func ToLineSpinResponse(resp model.SpinResult) dto.LineSpinResponse {
	return dto.LineSpinResponse{
		Board:            resp.Board,
		LineWins:         toLineWins(resp.LineWins),
		ScatterCount:     resp.ScatterCount,
		ScatterPayout:    resp.ScatterPayout,
		AwardedFreeSpins: resp.AwardedFreeSpins,
		TotalPayout:      resp.TotalPayout,
		Balance:          resp.Balance,
		InFreeSpin:       resp.InFreeSpin,
	}
}

func ToLineSpinRequest(bet int, buyBonus bool) dto.LineSpinRequest {
	return dto.LineSpinRequest{
		Bet:      bet,
		BuyBonus: buyBonus,
	}
}

func toLineWins(lineWins []model.LineWin) []dto.LineWin {
	result := make([]dto.LineWin, len(lineWins))
	for i, line := range lineWins {
		result[i] = dto.LineWin{
			Line:   line.Line,
			Symbol: line.Symbol,
			Count:  line.Count,
			Payout: line.Payout,
		}
	}
	return result
}
