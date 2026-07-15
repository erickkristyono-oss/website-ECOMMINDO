package httphandler

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"ecommindo/internal/usecase"
	"errors"
	"net/http"
)

type OrderHandler struct {
	uc domain.OrderUsecase
}

func NewOrderHandler(uc domain.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

var bankAccounts = []map[string]string{
	{"bank": "Bank BCA", "holder": "PT Ecommindo Jaya Persada", "number": "1234567890"},
	{"bank": "Bank Mandiri", "holder": "PT Ecommindo Jaya Persada", "number": "9876543210"},
}

type orderItemResponse struct {
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Price       float64 `json:"price"`
}

type orderResponse struct {
	OrderCode string               `json:"order_code"`
	Total     float64              `json:"total"`
	Status    string               `json:"status"`
	Items     []orderItemResponse  `json:"items"`
	Banks     []map[string]string  `json:"bank_accounts,omitempty"`
}

func toOrderResponse(o *entity.Order, withBanks bool) orderResponse {
	resp := orderResponse{
		OrderCode: o.OrderCode,
		Total:     o.Total,
		Status:    o.Status,
	}
	for _, it := range o.Items {
		resp.Items = append(resp.Items, orderItemResponse{
			ServiceID:   it.ServiceID,
			ServiceName: it.ServiceName,
			Price:       it.Price,
		})
	}
	if withBanks {
		resp.Banks = bankAccounts
	}
	return resp
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, _ := userIDFromContext(r.Context())

	order, err := h.uc.Checkout(r.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrCartEmpty) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Gagal membuat pesanan.")
		return
	}

	writeJSON(w, http.StatusCreated, toOrderResponse(order, true))
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := userIDFromContext(r.Context())

	orders, err := h.uc.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil riwayat pesanan.")
		return
	}

	resp := make([]orderResponse, 0, len(orders))
	for _, o := range orders {
		resp = append(resp, toOrderResponse(&o, false))
	}

	writeJSON(w, http.StatusOK, resp)
}
