package model

import (
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
	AchievementType string      `json:"achievement_type"` // academic, competition, organization, publication, certification, other
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Details         interface{} `json:"details"` // CompetitionDetails atau PublicationDetails
}

// Response model untuk submit prestasi
type SubmitAchievementResponse struct {
	AchievementID          string    `json:"achievement_id"`           // MongoDB ObjectID
	AchievementReferenceID uuid.UUID `json:"achievement_reference_id"` // PostgreSQL UUID
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
