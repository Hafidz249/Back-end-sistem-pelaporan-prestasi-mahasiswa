package model

import (
	"time"

	"github.com/google/uuid"
)

type Roles struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}