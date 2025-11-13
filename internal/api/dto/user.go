package dto

type DepositRequest struct {
	Amount int `json:"amount"`
}

type DepositResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
