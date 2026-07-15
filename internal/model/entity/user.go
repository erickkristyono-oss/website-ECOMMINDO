package entity

import "time"

type User struct {
	ID           int
	FullName     string
	Phone        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
