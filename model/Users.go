package model

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Credential string `json:"credential"` // username atau email
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token   string      `json:"token"`
	Profile UserProfile `json:"profile"`
}

type UserProfile struct {
	ID          uuid.UUID    `json:"id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	FullName    string       `json:"full_name"`
	Role        RoleInfo     `json:"role"`
	Permissions []Permission `json:"permissions"`
}

type RoleInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Permission struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Resource string    `json:"resource"`
	Action   string    `json:"action"`
}