package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"` // Penerima notifikasi (dosen wali)
	Type      string    `json:"type"`    // achievement_submitted, achievement_verified, etc
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Data      string    `json:"data"`       // JSON string untuk data tambahan
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationData struct {
	AchievementID          string    `json:"achievement_id"`
	AchievementReferenceID uuid.UUID `json:"achievement_reference_id"`
	StudentID              uuid.UUID `json:"student_id"`
	StudentName            string    `json:"student_name"`
	AchievementTitle       string    `json:"achievement_title"`
}
