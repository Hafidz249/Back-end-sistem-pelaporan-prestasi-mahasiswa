package repository

import (
	"POJECT_UAS/model"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthRepository_Login_Success_WithUsername(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authRepo := &AuthRepository{
		DB:        db,
		JWTSecret: "test-secret",
	}

	// Mock data
	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	// Mock query expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "testuser", "test@example.com", string(hashedPassword), "Test User", roleID, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE \(username = \$1 OR email = \$1\) AND is_active = true`).
		WithArgs("testuser").
		WillReturnRows(userRows)

	// Mock role query
	roleRows := sqlmock.NewRows([]string{"name"}).AddRow("student")
	mock.ExpectQuery(`SELECT name FROM roles WHERE id = \$1`).
		WithArgs(roleID).
		WillReturnRows(roleRows)

	// Mock permissions query
	permRows := sqlmock.NewRows([]string{"resource", "action"}).
		AddRow("achievements", "create").
		AddRow("achievements", "read")
	mock.ExpectQuery(`SELECT p.resource, p.action FROM permissions p`).
		WithArgs(roleID).
		WillReturnRows(permRows)

	// Execute
	result, err := authRepo.Login("testuser", "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.User.Username)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.Equal(t, "Test User", result.User.FullName)
	assert.Equal(t, "student", result.User.Role)
	assert.NotEmpty(t, result.Token)
	assert.Len(t, result.User.Permissions, 2)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepository_Login_Success_WithEmail(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authRepo := &AuthRepository{
		DB:        db,
		JWTSecret: "test-secret",
	}

	// Mock data
	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	// Mock query expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "testuser", "test@example.com", string(hashedPassword), "Test User", roleID, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE \(username = \$1 OR email = \$1\) AND is_active = true`).
		WithArgs("test@example.com").
		WillReturnRows(userRows)

	// Mock role query
	roleRows := sqlmock.NewRows([]string{"name"}).AddRow("student")
	mock.ExpectQuery(`SELECT name FROM roles WHERE id = \$1`).
		WithArgs(roleID).
		WillReturnRows(roleRows)

	// Mock permissions query
	permRows := sqlmock.NewRows([]string{"resource", "action"})
	mock.ExpectQuery(`SELECT p.resource, p.action FROM permissions p`).
		WithArgs(roleID).
		WillReturnRows(permRows)

	// Execute
	result, err := authRepo.Login("test@example.com", "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.User.Username)
	assert.Equal(t, "test@example.com", result.User.Email)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepository_Login_UserNotFound(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authRepo := &AuthRepository{
		DB:        db,
		JWTSecret: "test-secret",
	}

	// Mock query expectations - no user found
	mock.ExpectQuery(`SELECT (.+) FROM users WHERE \(username = \$1 OR email = \$1\) AND is_active = true`).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	// Execute
	result, err := authRepo.Login("nonexistent", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepository_Login_WrongPassword(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authRepo := &AuthRepository{
		DB:        db,
		JWTSecret: "test-secret",
	}

	// Mock data
	userID := uuid.New()
	roleID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	
	// Mock query expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "testuser", "test@example.com", string(hashedPassword), "Test User", roleID, true, time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE \(username = \$1 OR email = \$1\) AND is_active = true`).
		WithArgs("testuser").
		WillReturnRows(userRows)

	// Execute with wrong password
	result, err := authRepo.Login("testuser", "wrongpassword")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepository_Login_InactiveUser(t *testing.T) {
	// Setup mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	authRepo := &AuthRepository{
		DB:        db,
		JWTSecret: "test-secret",
	}

	// Mock query expectations - no active user found
	mock.ExpectQuery(`SELECT (.+) FROM users WHERE \(username = \$1 OR email = \$1\) AND is_active = true`).
		WithArgs("inactiveuser").
		WillReturnError(sql.ErrNoRows)

	// Execute
	result, err := authRepo.Login("inactiveuser", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}