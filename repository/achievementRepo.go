package repository

import (
	"POJECT_UAS/model"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

// GetAchievementsByStudentID mengambil semua achievements milik student (exclude deleted)
func (r *AchievementRepository) GetAchievementsByStudentID(studentID uuid.UUID) ([]model.Achievement, error) {
	ctx := context.Background()

	var achievements []model.Achievement
	collection := r.MongoDB.Collection("achievements")

	// Filter: hanya ambil yang tidak dihapus
	filter := primitive.M{
		"studentId": studentID,
		"$or": []primitive.M{
			{"isDeleted": primitive.M{"$exists": false}},
			{"isDeleted": false},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// GetAchievementReferenceByID mengambil achievement reference dari PostgreSQL
func (r *AchievementRepository) GetAchievementReferenceByID(referenceID uuid.UUID) (*model.AchievementReference, error) {
	var ref model.AchievementReference

	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`

	err := r.PostgresDB.QueryRow(query, referenceID).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ref, nil
}

// SubmitForVerification update status achievement dari draft ke submitted
func (r *AchievementRepository) SubmitForVerification(referenceID uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE achievement_references
		SET status = 'submitted', submitted_at = $1, updated_at = $2
		WHERE id = $3 AND status = 'draft'
	`

	result, err := r.PostgresDB.Exec(query, now, now, referenceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Tidak ada row yang diupdate (mungkin status bukan draft)
	}

	return nil
}

// CreateNotification membuat notifikasi baru
func (r *AchievementRepository) CreateNotification(notification model.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, data, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.PostgresDB.Exec(
		query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Data,
		notification.IsRead,
		notification.CreatedAt,
	)

	return err
}

// GetAdvisorByStudentID mengambil advisor_id dari student
func (r *AchievementRepository) GetAdvisorByStudentID(studentID uuid.UUID) (uuid.UUID, error) {
	var advisorID uuid.UUID

	query := `SELECT advisor_id FROM students WHERE id = $1`

	err := r.PostgresDB.QueryRow(query, studentID).Scan(&advisorID)
	if err != nil {
		return uuid.Nil, err
	}

	return advisorID, nil
}

// GetUserByID mengambil user data berdasarkan ID
func (r *AchievementRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	var user model.Users

	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.PostgresDB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteAchievement soft delete achievement (FR-005)
func (r *AchievementRepository) DeleteAchievement(referenceID uuid.UUID, mongoAchievementID string) error {
	ctx := context.Background()

	// 1. Soft delete di MongoDB
	objectID, err := primitive.ObjectIDFromHex(mongoAchievementID)
	if err != nil {
		return err
	}

	now := time.Now()
	collection := r.MongoDB.Collection("achievements")

	update := primitive.M{
		"$set": primitive.M{
			"isDeleted": true,
			"deletedAt": now,
		},
	}

	_, err = collection.UpdateOne(ctx, primitive.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	// 2. Update reference di PostgreSQL (soft delete)
	query := `
		UPDATE achievement_references
		SET status = 'deleted', updated_at = $1
		WHERE id = $2 AND status = 'draft'
	`

	result, err := r.PostgresDB.Exec(query, now, referenceID)
	if err != nil {
		// Rollback MongoDB jika gagal update PostgreSQL
		rollbackUpdate := primitive.M{
			"$set": primitive.M{
				"isDeleted": false,
			},
			"$unset": primitive.M{
				"deletedAt": "",
			},
		}
		collection.UpdateOne(ctx, primitive.M{"_id": objectID}, rollbackUpdate)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Rollback MongoDB jika tidak ada row yang diupdate
		rollbackUpdate := primitive.M{
			"$set": primitive.M{
				"isDeleted": false,
			},
			"$unset": primitive.M{
				"deletedAt": "",
			},
		}
		collection.UpdateOne(ctx, primitive.M{"_id": objectID}, rollbackUpdate)
		return sql.ErrNoRows
	}

	return nil
}


// GetLecturerByUserID mengambil data lecturer berdasarkan user_id
func (r *AchievementRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturers, error) {
	var lecturer model.Lecturers

	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	err := r.PostgresDB.QueryRow(query, userID).Scan(
		&lecturer.ID,
		&lecturer.UserID,
		&lecturer.LecturerID,
		&lecturer.Department,
		&lecturer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &lecturer, nil
}

// GetStudentIDsByAdvisor mengambil list student IDs yang dibimbing oleh advisor
func (r *AchievementRepository) GetStudentIDsByAdvisor(advisorID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT id FROM students WHERE advisor_id = $1
	`

	rows, err := r.PostgresDB.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentIDs []uuid.UUID
	for rows.Next() {
		var studentID uuid.UUID
		if err := rows.Scan(&studentID); err != nil {
			return nil, err
		}
		studentIDs = append(studentIDs, studentID)
	}

	return studentIDs, nil
}

// GetAchievementReferencesWithPagination mengambil achievement references dengan pagination
func (r *AchievementRepository) GetAchievementReferencesWithPagination(
	studentIDs []uuid.UUID,
	status string,
	page int,
	perPage int,
) ([]model.AchievementReference, int64, error) {
	// Build query dengan filter
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = ANY($1)
	`

	countQuery := `
		SELECT COUNT(*) FROM achievement_references WHERE student_id = ANY($1)
	`

	// Add status filter if provided
	if status != "" {
		query += ` AND status = $2`
		countQuery += ` AND status = $2`
	}

	// Add ordering
	query += ` ORDER BY created_at DESC`

	// Add pagination
	offset := (page - 1) * perPage
	query += ` LIMIT $3 OFFSET $4`

	// Get total count
	var totalCount int64
	if status != "" {
		err := r.PostgresDB.QueryRow(countQuery, pq.Array(studentIDs), status).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.PostgresDB.QueryRow(countQuery, pq.Array(studentIDs)).Scan(&totalCount)
		if err != nil {
			return nil, 0, err
		}
	}

	// Get paginated data
	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = r.PostgresDB.Query(query, pq.Array(studentIDs), status, perPage, offset)
	} else {
		rows, err = r.PostgresDB.Query(query, pq.Array(studentIDs), perPage, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var references []model.AchievementReference
	for rows.Next() {
		var ref model.AchievementReference
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		references = append(references, ref)
	}

	return references, totalCount, nil
}

// GetStudentByID mengambil student berdasarkan ID
func (r *AchievementRepository) GetStudentByID(studentID uuid.UUID) (*model.Student, error) {
	var student model.Student

	query := `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`

	err := r.PostgresDB.QueryRow(query, studentID).Scan(
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

// GetAchievementsByIDs mengambil multiple achievements dari MongoDB
func (r *AchievementRepository) GetAchievementsByIDs(achievementIDs []string) (map[string]model.Achievement, error) {
	ctx := context.Background()

	// Convert string IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range achievementIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return make(map[string]model.Achievement), nil
	}

	collection := r.MongoDB.Collection("achievements")

	// Find all achievements with given IDs
	cursor, err := collection.Find(ctx, primitive.M{
		"_id": primitive.M{"$in": objectIDs},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Map achievements by ID for easy lookup
	achievementsMap := make(map[string]model.Achievement)
	for cursor.Next(ctx) {
		var achievement model.Achievement
		if err := cursor.Decode(&achievement); err != nil {
			continue
		}
		achievementsMap[achievement.ID.Hex()] = achievement
	}

	return achievementsMap, nil
}
