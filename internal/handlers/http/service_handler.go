package httphandler

import (
	"ecommindo/internal/domain"
	"net/http"
)

type ServiceHandler struct {
	uc domain.ServiceUsecase
}

func NewServiceHandler(uc domain.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{uc: uc}
}

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	services, err := h.uc.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil daftar layanan.")
		return
	}

	writeJSON(w, http.StatusOK, services)
}
