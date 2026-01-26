package cascade

import (
	dto "casino_backend/internal/api/dto/cascade"
	"casino_backend/internal/converter"
	"casino_backend/internal/service"
	"casino_backend/pkg/req"
	"casino_backend/pkg/resp"
	"net/http"
)

type HandlerDeps struct {
	Serv service.CascadeService
}

type Handler struct {
	serv service.CascadeService
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{serv: deps.Serv}
}

func (h *Handler) Spin(w http.ResponseWriter, r *http.Request) {
	payload, err := req.Decode[dto.CascadeSpinRequest](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.serv.Spin(r.Context(), converter.ToCascadeSpin(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := converter.ToCascadeSpinResponse(*result)
	resp.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) BuyBonus(w http.ResponseWriter, r *http.Request) {
	payload, err := req.Decode[dto.BuyCascadeBonusRequest](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.serv.BuyBonus(r.Context(), payload.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteJSONResponse(w, http.StatusOK, map[string]string{"result": "ok"})
}
