package repository

import (
	"POJECT_UAS/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	PostgresDB *sql.DB
	MongoDB    *mongo.Database
}

func NewAchievementRepository(postgresDB *sql.DB, mongoDB *mongo.Database) *AchievementRepository {
	return &AchievementRepository{
		PostgresDB: postgresDB,
		MongoDB:    mongoDB,
	}
}

// SubmitAchievement menyimpan prestasi ke MongoDB dan reference ke PostgreSQL
func (r *AchievementRepository) SubmitAchievement(studentID uuid.UUID, req model.SubmitAchievementRequest) (*model.SubmitAchievementResponse, error) {
	ctx := context.Background()

	// 1. Simpan ke MongoDB
	achievement := model.Achievement{
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
	}

	collection := r.MongoDB.Collection("achievements")
	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	mongoID := result.InsertedID.(primitive.ObjectID)

	// 2. Simpan reference ke PostgreSQL
	referenceID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.PostgresDB.Exec(
		query,
		referenceID,
		studentID,
		mongoID.Hex(),
		"draft", // Status awal: draft
		now,
		now,
	)
	if err != nil {
		// Rollback: hapus dari MongoDB jika gagal insert ke PostgreSQL
		collection.DeleteOne(ctx, primitive.M{"_id": mongoID})
		return nil, err
	}

	// 3. Return response
	return &model.SubmitAchievementResponse{
		AchievementID:          mongoID.Hex(),
		AchievementReferenceID: referenceID,
		StudentID:              studentID,
		AchievementType:        req.AchievementType,
		Title:                  req.Title,
		Description:            req.Description,
		Status:                 "draft",
		CreatedAt:              now.Format(time.RFC3339),
	}, nil
}

// GetStudentByUserID mengambil data student berdasarkan user_id
func (r *AchievementRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	var student model.Student

	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1
	`

	err := r.PostgresDB.QueryRow(query, userID).Scan(
		&student.ID,
		&student.UserID,
		&student.StudentID,
		&student.ProgramStudy,
		&student.AcademicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

// GetAchievementByID mengambil achievement dari MongoDB berdasarkan ID
func (r *AchievementRepository) GetAchievementByID(achievementID string) (*model.Achievement, error) {
	ctx := context.Background()

	objectID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	collection := r.MongoDB.Collection("achievements")

	err = collection.FindOne(ctx, primitive.M{"_id": objectID}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

// GetAchievementsByStudentID mengambil semua achievements milik student
func (r *AchievementRepository) GetAchievementsByStudentID(studentID uuid.UUID) ([]model.Achievement, error) {
	ctx := context.Background()

	var achievements []model.Achievement
	collection := r.MongoDB.Collection("achievements")

	cursor, err := collection.Find(ctx, primitive.M{"studentId": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}
