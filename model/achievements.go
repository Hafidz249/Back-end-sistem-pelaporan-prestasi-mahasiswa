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
