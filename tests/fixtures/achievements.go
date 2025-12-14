package fixtures

import (
	"POJECT_UAS/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AchievementFixtures provides test data for achievements
type AchievementFixtures struct{}

// NewAchievementFixtures creates a new AchievementFixtures instance
func NewAchievementFixtures() *AchievementFixtures {
	return &AchievementFixtures{}
}

// ValidAchievement returns a valid achievement for testing
func (f *AchievementFixtures) ValidAchievement() *model.Achievement {
	return &model.Achievement{
		ID:              primitive.NewObjectID(),
		StudentID:       uuid.New(),
		AchievementType: "akademik",
		Title:           "Juara 1 Programming Contest",
		Description:     "Kompetisi programming tingkat nasional",
		Details: map[string]interface{}{
			"level":    "nasional",
			"category": "programming",
			"date":     "2024-01-15",
			"location": "Jakarta",
		},
	}
}

// ValidAchievementReference returns a valid achievement reference
func (f *AchievementFixtures) ValidAchievementReference() *model.AchievementReference {
	return &model.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.New(),
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// SubmittedAchievementReference returns a submitted achievement reference
func (f *AchievementFixtures) SubmittedAchievementReference() *model.AchievementReference {
	now := time.Now()
	return &model.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.New(),
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             "submitted",
		SubmittedAt:        &now,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// VerifiedAchievementReference returns a verified achievement reference
func (f *AchievementFixtures) VerifiedAchievementReference() *model.AchievementReference {
	now := time.Now()
	verifierID := uuid.New()
	return &model.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.New(),
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             "verified",
		SubmittedAt:        &now,
		VerifiedAt:         &now,
		VerifiedBy:         &verifierID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// RejectedAchievementReference returns a rejected achievement reference
func (f *AchievementFixtures) RejectedAchievementReference() *model.AchievementReference {
	now := time.Now()
	verifierID := uuid.New()
	rejectionNote := "Dokumen tidak lengkap"
	return &model.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.New(),
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             "rejected",
		SubmittedAt:        &now,
		VerifiedAt:         &now,
		VerifiedBy:         &verifierID,
		RejectionNote:      &rejectionNote,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// ValidSubmitAchievementRequest returns a valid submit achievement request
func (f *AchievementFixtures) ValidSubmitAchievementRequest() model.SubmitAchievementRequest {
	return model.SubmitAchievementRequest{
		AchievementType: "akademik",
		Title:           "Juara 1 Programming Contest",
		Description:     "Kompetisi programming tingkat nasional",
		Details: map[string]interface{}{
			"level":    "nasional",
			"category": "programming",
			"date":     "2024-01-15",
		},
	}
}

// InvalidSubmitAchievementRequest returns an invalid submit achievement request
func (f *AchievementFixtures) InvalidSubmitAchievementRequest() model.SubmitAchievementRequest {
	return model.SubmitAchievementRequest{
		AchievementType: "", // Invalid: empty type
		Title:           "", // Invalid: empty title
		Description:     "", // Invalid: empty description
		Details:         nil,
	}
}

// ValidSubmitAchievementResponse returns a valid submit achievement response
func (f *AchievementFixtures) ValidSubmitAchievementResponse() *model.SubmitAchievementResponse {
	return &model.SubmitAchievementResponse{
		AchievementID:          primitive.NewObjectID().Hex(),
		AchievementReferenceID: uuid.New(),
		StudentID:              uuid.New(),
		AchievementType:        "akademik",
		Title:                  "Juara 1 Programming Contest",
		Description:            "Kompetisi programming tingkat nasional",
		Status:                 "draft",
		CreatedAt:              time.Now().Format(time.RFC3339),
	}
}

// ValidStudent returns a valid student for testing
func (f *AchievementFixtures) ValidStudent() *model.Student {
	advisorID := uuid.New()
	return &model.Student{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		StudentID:    "2021001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2021",
		AdvisorID:    &advisorID,
		CreatedAt:    time.Now(),
	}
}

// ValidLecturer returns a valid lecturer for testing
func (f *AchievementFixtures) ValidLecturer() *model.Lecturers {
	return &model.Lecturers{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		LecturerID: "L001",
		Department: "Teknik Informatika",
		CreatedAt:  time.Now(),
	}
}

// ValidNotification returns a valid notification for testing
func (f *AchievementFixtures) ValidNotification() model.Notification {
	return model.Notification{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Type:      "achievement_submitted",
		Title:     "Prestasi Baru Disubmit",
		Message:   "Mahasiswa telah submit prestasi untuk verifikasi",
		Data:      `{"achievement_id": "` + uuid.New().String() + `"}`,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
}

// MultipleAchievements returns multiple achievements for testing
func (f *AchievementFixtures) MultipleAchievements(count int) []model.Achievement {
	achievements := make([]model.Achievement, count)
	for i := 0; i < count; i++ {
		achievements[i] = model.Achievement{
			ID:              primitive.NewObjectID(),
			StudentID:       uuid.New(),
			AchievementType: "akademik",
			Title:           fmt.Sprintf("Achievement %d", i+1),
			Description:     fmt.Sprintf("Description for achievement %d", i+1),
			Details: map[string]interface{}{
				"level": "nasional",
				"rank":  i + 1,
			},
		}
	}
	return achievements
}

// MultipleAchievementReferences returns multiple achievement references for testing
func (f *AchievementFixtures) MultipleAchievementReferences(count int) []model.AchievementReference {
	references := make([]model.AchievementReference, count)
	statuses := []string{"draft", "submitted", "verified", "rejected"}
	
	for i := 0; i < count; i++ {
		status := statuses[i%len(statuses)]
		ref := model.AchievementReference{
			ID:                 uuid.New(),
			StudentID:          uuid.New(),
			MongoAchievementID: primitive.NewObjectID().Hex(),
			Status:             status,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		
		if status == "submitted" || status == "verified" || status == "rejected" {
			now := time.Now()
			ref.SubmittedAt = &now
		}
		
		if status == "verified" || status == "rejected" {
			now := time.Now()
			verifierID := uuid.New()
			ref.VerifiedAt = &now
			ref.VerifiedBy = &verifierID
		}
		
		if status == "rejected" {
			note := "Test rejection note"
			ref.RejectionNote = &note
		}
		
		references[i] = ref
	}
	return references
}

// ValidAchievementStatistics returns valid achievement statistics
func (f *AchievementFixtures) ValidAchievementStatistics() *model.AchievementStatistics {
	month := 1
	return &model.AchievementStatistics{
		TotalByType: []model.TypeStatistic{
			{
				AchievementType: "akademik",
				Count:           15,
				Percentage:      75.0,
			},
			{
				AchievementType: "non-akademik",
				Count:           5,
				Percentage:      25.0,
			},
		},
		TotalByPeriod: []model.PeriodStatistic{
			{
				Period: "2024-01",
				Count:  10,
				Year:   2024,
				Month:  &month,
			},
		},
		TopStudents: []model.TopStudent{
			{
				StudentID:       uuid.New(),
				StudentIDNumber: "2021001",
				FullName:        "Top Student",
				ProgramStudy:    "Teknik Informatika",
				AcademicYear:    "2021",
				TotalCount:      10,
				VerifiedCount:   8,
			},
		},
		CompetitionLevels: []model.LevelStatistic{
			{
				Level:      "nasional",
				Count:      12,
				Percentage: 60.0,
			},
			{
				Level:      "regional",
				Count:      8,
				Percentage: 40.0,
			},
		},
		Summary: model.StatisticSummary{
			TotalAchievements:    20,
			VerifiedAchievements: 15,
			PendingAchievements:  3,
			RejectedAchievements: 2,
			TotalStudents:        5,
			DateRange: model.DateRange{
				StartDate: time.Now().AddDate(0, -6, 0),
				EndDate:   time.Now(),
			},
		},
	}
}