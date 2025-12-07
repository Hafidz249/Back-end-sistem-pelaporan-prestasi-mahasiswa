package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       uuid.UUID          `bson:"studentId" json:"student_id"`
	AchievementType string             `bson:"achievementType" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         interface{}        `bson:"details" json:"details"`
	IsDeleted       bool               `bson:"isDeleted,omitempty" json:"is_deleted,omitempty"`
	DeletedAt       *time.Time         `bson:"deletedAt,omitempty" json:"deleted_at,omitempty"`
}

type CompetitionDetails struct {
	CompetitionName  *string `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel *string `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
	Rank             *int    `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        *string `bson:"medalType,omitempty" json:"medal_type,omitempty"`
}

type PublicationDetails struct {
	PublicationType  *string   `bson:"publicationType,omitempty" json:"publication_type,omitempty"`
	PublicationTitle *string   `bson:"publicationTitle,omitempty" json:"publication_title,omitempty"`
	Authors          *[]string `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher        *string   `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN             *string   `bson:"issn,omitempty" json:"issn,omitempty"`
}

// Request model untuk submit prestasi
type SubmitAchievementRequest struct {
	AchievementType string      `json:"achievement_type"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Details         interface{} `json:"details"`
}

// Response model untuk submit prestasi
type SubmitAchievementResponse struct {
	AchievementID          string    `json:"achievement_id"`
	AchievementReferenceID uuid.UUID `json:"achievement_reference_id"`
	StudentID              uuid.UUID `json:"student_id"`
	AchievementType        string    `json:"achievement_type"`
	Title                  string    `json:"title"`
	Description            string    `json:"description"`
	Status                 string    `json:"status"`
	CreatedAt              string    `json:"created_at"`
}

// Response model untuk submit verification
type SubmitForVerificationResponse struct {
	AchievementReferenceID uuid.UUID `json:"achievement_reference_id"`
	Status                 string    `json:"status"`
	SubmittedAt            string    `json:"submitted_at"`
	Message                string    `json:"message"`
}

// Achievement with student info untuk dosen
type AchievementWithStudent struct {
	Achievement
	AchievementReferenceID uuid.UUID  `json:"achievement_reference_id"`
	Status                 string     `json:"status"`
	SubmittedAt            *time.Time `json:"submitted_at,omitempty"`
	VerifiedAt             *time.Time `json:"verified_at,omitempty"`
	StudentName            string     `json:"student_name"`
	StudentIDNumber        string     `json:"student_id_number"`
	ProgramStudy           string     `json:"program_study"`
}

// Pagination response
type PaginatedAchievementsResponse struct {
	Data       []AchievementWithStudent `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}
