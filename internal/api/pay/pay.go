package pay

import (
	dto "casino_backend/internal/api/dto/pay"
	"casino_backend/internal/service"
	"casino_backend/pkg/req"
	"casino_backend/pkg/resp"
	"net/http"
)

type HandlerDeps struct {
	serv service.PaymentService
}

type Handler struct {
	serv service.PaymentService
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		serv: deps.serv,
	}
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	requestBody, err := req.Decode[dto.DepositRequest](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.serv.Deposit(r.Context(), requestBody.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := h.serv.GetBalance(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"balance": balance})
}
