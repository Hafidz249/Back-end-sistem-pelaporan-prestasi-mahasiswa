package repository

import (
	"POJECT_UAS/model"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser membuat user baru (FR-009)
func (r *UserRepository) CreateUser(req model.CreateUserRequest) (*model.Users, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, username, email, full_name, role_id, is_active, created_at, updated_at
	`

	var user model.Users
	err = r.DB.QueryRow(
		query,
		userID,
		req.Username,
		req.Email,
		string(hashedPassword),
		req.FullName,
		req.RoleID,
		true, // is_active default true
		now,
		now,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

// UpdateUser update user data (FR-009)
func (r *UserRepository) UpdateUser(userID uuid.UUID, req model.UpdateUserRequest) (*model.Users, error) {
	now := time.Now()

	query := `
		UPDATE users
		SET full_name = $1, role_id = $2, is_active = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, username, email, full_name, role_id, is_active, created_at, updated_at
	`

	var user model.Users
	err := r.DB.QueryRow(
		query,
		req.FullName,
		req.RoleID,
		req.IsActive,
		now,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

// DeleteUser soft delete user (FR-009)
func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	now := time.Now()

	query := `
		UPDATE users
		SET is_active = false, updated_at = $1
		WHERE id = $2
	`

	result, err := r.DB.Exec(query, now, userID)
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

// GetAllUsers get all users with pagination (FR-009)
func (r *UserRepository) GetAllUsers(page, perPage int) ([]model.Users, int64, error) {
	// Count total
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * perPage
	query := `
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.Users
	for rows.Next() {
		var user model.Users
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FullName,
			&user.RoleID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalCount, nil
}

// CreateStudentProfile membuat profile student (FR-009)
func (r *UserRepository) CreateStudentProfile(req model.CreateStudentProfileRequest) (*model.Student, error) {
	studentID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, student_id, program_study, academic_year, advisor_id, created_at
	`

	var student model.Student
	err := r.DB.QueryRow(
		query,
		studentID,
		req.UserID,
		req.StudentIDNumber,
		req.ProgramStudy,
		req.AcademicYear,
		req.AdvisorID,
		now,
	).Scan(
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

// CreateLecturerProfile membuat profile lecturer (FR-009)
func (r *UserRepository) CreateLecturerProfile(req model.CreateLecturerProfileRequest) (*model.Lecturers, error) {
	lecturerID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, lecturer_id, department, created_at
	`

	var lecturer model.Lecturers
	err := r.DB.QueryRow(
		query,
		lecturerID,
		req.UserID,
		req.LecturerIDNumber,
		req.Department,
		now,
	).Scan(
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

// UpdateStudentAdvisor update advisor untuk student (FR-009)
func (r *UserRepository) UpdateStudentAdvisor(studentID, advisorID uuid.UUID) error {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`

	result, err := r.DB.Exec(query, advisorID, studentID)
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

// GetAllRoles get all roles (FR-009)
func (r *UserRepository) GetAllRoles() ([]model.Roles, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		ORDER BY name
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Roles
	for rows.Next() {
		var role model.Roles
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
