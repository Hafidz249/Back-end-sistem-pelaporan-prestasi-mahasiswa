package repository

import (
	"POJECT_UAS/model"
	"context"
	"database/sql"
	"strconv"
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


// VerifyAchievement approve atau reject prestasi (FR-007)
func (r *AchievementRepository) VerifyAchievement(referenceID uuid.UUID, verifierID uuid.UUID, action string, note *string) error {
	now := time.Now()

	var query string
	var args []interface{}

	if action == "approve" {
		// Update status ke 'verified'
		query = `
			UPDATE achievement_references
			SET status = 'verified', verified_by = $1, verified_at = $2, updated_at = $3
			WHERE id = $4 AND status = 'submitted'
		`
		args = []interface{}{verifierID, now, now, referenceID}
	} else if action == "reject" {
		// Update status ke 'rejected'
		query = `
			UPDATE achievement_references
			SET status = 'rejected', verified_by = $1, verified_at = $2, rejection_note = $3, updated_at = $4
			WHERE id = $5 AND status = 'submitted'
		`
		args = []interface{}{verifierID, now, note, now, referenceID}
	} else {
		return sql.ErrNoRows
	}

	result, err := r.PostgresDB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Tidak ada row yang diupdate (mungkin status bukan submitted)
	}

	return nil
}

// CheckLecturerOwnsStudent mengecek apakah lecturer adalah advisor dari student
func (r *AchievementRepository) CheckLecturerOwnsStudent(lecturerID uuid.UUID, studentID uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM students WHERE id = $1 AND advisor_id = $2`

	err := r.PostgresDB.QueryRow(query, studentID, lecturerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}


// VerifyAchievement approve prestasi (FR-007)
func (r *AchievementRepository) VerifyAchievement(referenceID uuid.UUID, verifiedBy uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE achievement_references
		SET status = 'verified', verified_at = $1, verified_by = $2, updated_at = $3
		WHERE id = $4 AND status = 'submitted'
	`

	result, err := r.PostgresDB.Exec(query, now, verifiedBy, now, referenceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// RejectAchievement reject prestasi dengan rejection note (FR-007)
func (r *AchievementRepository) RejectAchievement(referenceID uuid.UUID, verifiedBy uuid.UUID, rejectionNote string) error {
	now := time.Now()

	query := `
		UPDATE achievement_references
		SET status = 'rejected', verified_at = $1, verified_by = $2, rejection_note = $3, updated_at = $4
		WHERE id = $5 AND status = 'submitted'
	`

	result, err := r.PostgresDB.Exec(query, now, verifiedBy, rejectionNote, now, referenceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetAllAchievementReferencesWithPagination mengambil semua achievement references untuk admin (FR-010)
func (r *AchievementRepository) GetAllAchievementReferencesWithPagination(
	status string,
	achievementType string,
	page int,
	perPage int,
) ([]model.AchievementReference, int64, error) {
	// Build base query
	baseWhere := "WHERE 1=1"
	var args []interface{}
	argIndex := 1

	// Add status filter if provided
	if status != "" {
		baseWhere += " AND ar.status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM achievement_references ar
		LEFT JOIN students s ON ar.student_id = s.id
		LEFT JOIN users u ON s.user_id = u.id
	` + baseWhere

	// Select query
	selectQuery := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status, 
		       ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note, 
		       ar.created_at, ar.updated_at
		FROM achievement_references ar
		LEFT JOIN students s ON ar.student_id = s.id
		LEFT JOIN users u ON s.user_id = u.id
	` + baseWhere + `
		ORDER BY ar.created_at DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	// Add pagination args
	offset := (page - 1) * perPage
	paginationArgs := append(args, perPage, offset)

	// Get total count
	var totalCount int64
	err := r.PostgresDB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	rows, err := r.PostgresDB.Query(selectQuery, paginationArgs...)
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
// GetAchievementStatistics mengambil statistik prestasi (FR-011)
func (r *AchievementRepository) GetAchievementStatistics(
	studentIDs []uuid.UUID, // kosong jika admin (all students)
	startDate *time.Time,
	endDate *time.Time,
	achievementType *string,
	status *string,
) (*model.AchievementStatistics, error) {
	ctx := context.Background()

	// Build base filter untuk PostgreSQL
	baseWhere := "WHERE 1=1"
	var args []interface{}
	argIndex := 1

	// Filter by student IDs (untuk dosen wali atau mahasiswa)
	if len(studentIDs) > 0 {
		baseWhere += " AND ar.student_id = ANY($" + strconv.Itoa(argIndex) + ")"
		args = append(args, pq.Array(studentIDs))
		argIndex++
	}

	// Filter by date range
	if startDate != nil {
		baseWhere += " AND ar.created_at >= $" + strconv.Itoa(argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		baseWhere += " AND ar.created_at <= $" + strconv.Itoa(argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	// Filter by status
	if status != nil {
		baseWhere += " AND ar.status = $" + strconv.Itoa(argIndex)
		args = append(args, *status)
		argIndex++
	}

	// 1. Get summary statistics
	summaryQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN ar.status = 'verified' THEN 1 END) as verified,
			COUNT(CASE WHEN ar.status = 'submitted' THEN 1 END) as pending,
			COUNT(CASE WHEN ar.status = 'rejected' THEN 1 END) as rejected,
			COUNT(DISTINCT ar.student_id) as total_students,
			MIN(ar.created_at) as min_date,
			MAX(ar.created_at) as max_date
		FROM achievement_references ar
	` + baseWhere

	var summary model.StatisticSummary
	var minDate, maxDate sql.NullTime
	err := r.PostgresDB.QueryRow(summaryQuery, args...).Scan(
		&summary.TotalAchievements,
		&summary.VerifiedAchievements,
		&summary.PendingAchievements,
		&summary.RejectedAchievements,
		&summary.TotalStudents,
		&minDate,
		&maxDate,
	)
	if err != nil {
		return nil, err
	}

	if minDate.Valid && maxDate.Valid {
		summary.DateRange = model.DateRange{
			StartDate: minDate.Time,
			EndDate:   maxDate.Time,
		}
	}

	// 2. Get achievement references untuk analisis lebih lanjut
	referencesQuery := `
		SELECT ar.mongo_achievement_id, ar.student_id, ar.status, ar.created_at,
		       s.student_id as student_number, u.full_name, s.program_study, s.academic_year
		FROM achievement_references ar
		LEFT JOIN students s ON ar.student_id = s.id
		LEFT JOIN users u ON s.user_id = u.id
	` + baseWhere + ` ORDER BY ar.created_at DESC`

	rows, err := r.PostgresDB.Query(referencesQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mongoIDs []string
	var studentMap = make(map[uuid.UUID]*model.TopStudent)
	var periodMap = make(map[string]int64)

	for rows.Next() {
		var mongoID string
		var studentID uuid.UUID
		var status string
		var createdAt time.Time
		var studentNumber, fullName, programStudy, academicYear sql.NullString

		err := rows.Scan(&mongoID, &studentID, &status, &createdAt,
			&studentNumber, &fullName, &programStudy, &academicYear)
		if err != nil {
			continue
		}

		mongoIDs = append(mongoIDs, mongoID)

		// Build student statistics
		if studentNumber.Valid && fullName.Valid {
			if student, exists := studentMap[studentID]; exists {
				student.TotalCount++
				if status == "verified" {
					student.VerifiedCount++
				}
			} else {
				studentMap[studentID] = &model.TopStudent{
					StudentID:       studentID,
					StudentIDNumber: studentNumber.String,
					FullName:        fullName.String,
					ProgramStudy:    programStudy.String,
					AcademicYear:    academicYear.String,
					TotalCount:      1,
					VerifiedCount:   0,
				}
				if status == "verified" {
					studentMap[studentID].VerifiedCount = 1
				}
			}
		}

		// Build period statistics (monthly)
		period := createdAt.Format("2006-01")
		periodMap[period]++
	}

	// 3. Get achievements from MongoDB untuk type dan level statistics
	var typeMap = make(map[string]int64)
	var levelMap = make(map[string]int64)

	if len(mongoIDs) > 0 {
		// Convert to ObjectIDs
		var objectIDs []primitive.ObjectID
		for _, id := range mongoIDs {
			if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
				objectIDs = append(objectIDs, objectID)
			}
		}

		if len(objectIDs) > 0 {
			collection := r.MongoDB.Collection("achievements")
			cursor, err := collection.Find(ctx, primitive.M{
				"_id": primitive.M{"$in": objectIDs},
			})
			if err == nil {
				defer cursor.Close(ctx)

				for cursor.Next(ctx) {
					var achievement model.Achievement
					if err := cursor.Decode(&achievement); err == nil {
						// Filter by achievement type if specified
						if achievementType == nil || *achievementType == "" || achievement.AchievementType == *achievementType {
							typeMap[achievement.AchievementType]++

							// Extract level from details (assuming it's stored in details)
							if level, ok := achievement.Details["level"].(string); ok {
								levelMap[level]++
							} else {
								levelMap["unknown"]++
							}
						}
					}
				}
			}
		}
	}

	// 4. Build response
	statistics := &model.AchievementStatistics{
		Summary: summary,
	}

	// Convert type map to slice
	var totalTypes int64
	for _, count := range typeMap {
		totalTypes += count
	}
	for achType, count := range typeMap {
		percentage := float64(0)
		if totalTypes > 0 {
			percentage = float64(count) / float64(totalTypes) * 100
		}
		statistics.TotalByType = append(statistics.TotalByType, model.TypeStatistic{
			AchievementType: achType,
			Count:           count,
			Percentage:      percentage,
		})
	}

	// Convert period map to slice
	for period, count := range periodMap {
		year, _ := strconv.Atoi(period[:4])
		month, _ := strconv.Atoi(period[5:7])
		statistics.TotalByPeriod = append(statistics.TotalByPeriod, model.PeriodStatistic{
			Period: period,
			Count:  count,
			Year:   year,
			Month:  &month,
		})
	}

	// Convert student map to top students slice (sorted by total count)
	for _, student := range studentMap {
		statistics.TopStudents = append(statistics.TopStudents, *student)
	}

	// Convert level map to slice
	var totalLevels int64
	for _, count := range levelMap {
		totalLevels += count
	}
	for level, count := range levelMap {
		percentage := float64(0)
		if totalLevels > 0 {
			percentage = float64(count) / float64(totalLevels) * 100
		}
		statistics.CompetitionLevels = append(statistics.CompetitionLevels, model.LevelStatistic{
			Level:      level,
			Count:      count,
			Percentage: percentage,
		})
	}

	return statistics, nil
}