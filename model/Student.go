package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	StudentID     string    `json:"student_id"`
	Program_Study string    `json:"program_study"`
	Academic_Year string    `json:"academic_year"`
	AdvisorID     uuid.UUID `json:"advisor_id"`
	Created_at    time.Time `json:"created_at"`
}