package usecase

import "errors"

var (
	ErrEmailTaken         = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrServiceNotFound    = errors.New("layanan tidak ditemukan")
	ErrCartEmpty          = errors.New("keranjang masih kosong")
)
