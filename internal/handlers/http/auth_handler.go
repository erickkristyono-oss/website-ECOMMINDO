package httphandler

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"ecommindo/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type AuthHandler struct {
	uc domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type userResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func toUserResponse(u *entity.User) userResponse {
	return userResponse{ID: u.ID, FullName: u.FullName, Phone: u.Phone, Email: u.Email}
}

type registerRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Data yang dikirim tidak valid.")
		return
	}

	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.FullName == "" || req.Phone == "" || req.Email == "" || len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "Nama, nomor HP, email wajib diisi dan password minimal 6 karakter.")
		return
	}

	user, token, err := h.uc.Register(r.Context(), req.FullName, req.Phone, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailTaken) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Gagal mendaftarkan akun.")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  toUserResponse(user),
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Data yang dikirim tidak valid.")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, token, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Gagal login.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  toUserResponse(user),
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Sesi tidak valid.")
		return
	}

	user, err := h.uc.Me(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Pengguna tidak ditemukan.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"user": toUserResponse(user)})
}
