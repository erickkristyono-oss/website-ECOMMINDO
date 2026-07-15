package httphandler

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
)

type CartHandler struct {
	uc domain.CartUsecase
}

func NewCartHandler(uc domain.CartUsecase) *CartHandler {
	return &CartHandler{uc: uc}
}

func (h *CartHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := userIDFromContext(r.Context())

	summary, err := h.uc.Get(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil keranjang.")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

type addCartRequest struct {
	ServiceID string `json:"service_id"`
}

func (h *CartHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID, _ := userIDFromContext(r.Context())

	var req addCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ServiceID == "" {
		writeError(w, http.StatusBadRequest, "service_id wajib diisi.")
		return
	}

	summary, err := h.uc.Add(r.Context(), userID, req.ServiceID)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Gagal menambahkan layanan ke keranjang.")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func (h *CartHandler) Remove(w http.ResponseWriter, r *http.Request) {
	userID, _ := userIDFromContext(r.Context())
	serviceID := r.PathValue("serviceID")

	summary, err := h.uc.Remove(r.Context(), userID, serviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal menghapus layanan dari keranjang.")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}
