package repository

import (
	"POJECT_UAS/model"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock data
	roleID := uuid.New()
	req := model.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleID:   roleID,
	}

	// Mock query expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), "newuser", "newuser@example.com", "New User", roleID, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), "newuser", "newuser@example.com", sqlmock.AnyArg(), "New User", roleID, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(userRows)

	// Execute
	result, err := userRepo.CreateUser(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newuser", result.Username)
	assert.Equal(t, "newuser@example.com", result.Email)
	assert.Equal(t, "New User", result.FullName)
	assert.Equal(t, roleID, result.RoleID)
	assert.True(t, result.IsActive)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock data
	userID := uuid.New()
	roleID := uuid.New()
	req := model.UpdateUserRequest{
		FullName: "Updated User",
		RoleID:   roleID,
		IsActive: true,
	}

	// Mock query expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "testuser", "test@example.com", "Updated User", roleID, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery(`UPDATE users SET full_name = \$1, role_id = \$2, is_active = \$3, updated_at = \$4 WHERE id = \$5`).
		WithArgs("Updated User", roleID, true, sqlmock.AnyArg(), userID).
		WillReturnRows(userRows)

	// Execute
	result, err := userRepo.UpdateUser(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated User", result.FullName)
	assert.Equal(t, roleID, result.RoleID)
	assert.True(t, result.IsActive)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock data
	userID := uuid.New()

	// Mock query expectations
	mock.ExpectExec(`UPDATE users SET is_active = false, updated_at = \$1 WHERE id = \$2`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Execute
	err = userRepo.DeleteUser(userID)

	// Assert
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_NotFound(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock data
	userID := uuid.New()

	// Mock query expectations - no rows affected
	mock.ExpectExec(`UPDATE users SET is_active = false, updated_at = \$1 WHERE id = \$2`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	// Execute
	err = userRepo.DeleteUser(userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAllUsers_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
		WillReturnRows(countRows)

	// Mock data query
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		uuid.New(), "user1", "user1@example.com", "User 1", uuid.New(), true, time.Now(), time.Now(),
	).AddRow(
		uuid.New(), "user2", "user2@example.com", "User 2", uuid.New(), true, time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM users ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(10, 0).
		WillReturnRows(userRows)

	// Execute
	users, totalCount, err := userRepo.GetAllUsers(1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, int64(2), totalCount)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateStudentProfile_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock data
	userID := uuid.New()
	advisorID := uuid.New()
	req := model.CreateStudentProfileRequest{
		UserID:          userID,
		StudentIDNumber: "2021001",
		ProgramStudy:    "Teknik Informatika",
		AcademicYear:    "2021",
		AdvisorID:       &advisorID,
	}

	// Mock query expectations
	studentRows := sqlmock.NewRows([]string{
		"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "created_at",
	}).AddRow(
		uuid.New(), userID, "2021001", "Teknik Informatika", "2021", advisorID, time.Now(),
	)

	mock.ExpectQuery(`INSERT INTO students`).
		WithArgs(sqlmock.AnyArg(), userID, "2021001", "Teknik Informatika", "2021", advisorID, sqlmock.AnyArg()).
		WillReturnRows(studentRows)

	// Execute
	result, err := userRepo.CreateStudentProfile(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "2021001", result.StudentID)
	assert.Equal(t, "Teknik Informatika", result.ProgramStudy)
	assert.Equal(t, "2021", result.AcademicYear)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetAllRoles_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := NewUserRepository(db)

	// Mock query expectations
	roleRows := sqlmock.NewRows([]string{
		"id", "name", "description", "created_at",
	}).AddRow(
		uuid.New(), "admin", "Administrator", time.Now(),
	).AddRow(
		uuid.New(), "student", "Student", time.Now(),
	).AddRow(
		uuid.New(), "lecturer", "Lecturer", time.Now(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM roles ORDER BY name`).
		WillReturnRows(roleRows)

	// Execute
	roles, err := userRepo.GetAllRoles()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, roles, 3)
	assert.Equal(t, "admin", roles[0].Name)
	assert.Equal(t, "student", roles[1].Name)
	assert.Equal(t, "lecturer", roles[2].Name)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}