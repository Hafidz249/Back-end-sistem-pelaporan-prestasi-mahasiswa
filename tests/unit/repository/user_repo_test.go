package repository_test

import (
	"POJECT_UAS/repository"
	"POJECT_UAS/tests/fixtures"
	"POJECT_UAS/tests/mocks"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewUserFixtures()

	req := fixtures.ValidCreateUserRequest()
	expectedUser := fixtures.ValidUser()

	// Setup mock expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.FullName,
		expectedUser.RoleID, expectedUser.IsActive, expectedUser.CreatedAt, expectedUser.UpdatedAt,
	)

	mockDB.PostgresMock.ExpectQuery(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), req.Username, req.Email, sqlmock.AnyArg(), req.FullName, req.RoleID, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(userRows)

	// Act
	result, err := userRepo.CreateUser(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.FullName, result.FullName)
	assert.Equal(t, expectedUser.RoleID, result.RoleID)
	assert.True(t, result.IsActive)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_DatabaseError(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewUserFixtures()

	req := fixtures.ValidCreateUserRequest()

	// Setup mock expectations - database error
	mockDB.PostgresMock.ExpectQuery(`INSERT INTO users`).
		WithArgs(sqlmock.AnyArg(), req.Username, req.Email, sqlmock.AnyArg(), req.FullName, req.RoleID, true, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Act
	result, err := userRepo.CreateUser(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewUserFixtures()

	userID := uuid.New()
	req := fixtures.ValidUpdateUserRequest()
	expectedUser := fixtures.ValidUser()
	expectedUser.ID = userID
	expectedUser.FullName = req.FullName
	expectedUser.RoleID = req.RoleID

	// Setup mock expectations
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.FullName,
		expectedUser.RoleID, expectedUser.IsActive, expectedUser.CreatedAt, expectedUser.UpdatedAt,
	)

	mockDB.PostgresMock.ExpectQuery(`UPDATE users SET full_name = \$1, role_id = \$2, is_active = \$3, updated_at = \$4 WHERE id = \$5`).
		WithArgs(req.FullName, req.RoleID, req.IsActive, sqlmock.AnyArg(), userID).
		WillReturnRows(userRows)

	// Act
	result, err := userRepo.UpdateUser(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.FullName, result.FullName)
	assert.Equal(t, expectedUser.RoleID, result.RoleID)
	assert.Equal(t, req.IsActive, result.IsActive)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	userID := uuid.New()

	// Setup mock expectations
	mockDB.PostgresMock.ExpectExec(`UPDATE users SET is_active = false, updated_at = \$1 WHERE id = \$2`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Act
	err = userRepo.DeleteUser(userID)

	// Assert
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_DeleteUser_NotFound(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	userID := uuid.New()

	// Setup mock expectations - no rows affected
	mockDB.PostgresMock.ExpectExec(`UPDATE users SET is_active = false, updated_at = \$1 WHERE id = \$2`).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	// Act
	err = userRepo.DeleteUser(userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_GetAllUsers_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewUserFixtures()

	page := 1
	perPage := 10
	totalCount := int64(25)
	users := fixtures.MultipleUsers(2)

	// Setup mock expectations for count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(totalCount)
	mockDB.PostgresMock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
		WillReturnRows(countRows)

	// Setup mock expectations for data query
	userRows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	})
	for _, user := range users {
		userRows.AddRow(
			user.ID, user.Username, user.Email, user.FullName,
			user.RoleID, user.IsActive, user.CreatedAt, user.UpdatedAt,
		)
	}

	mockDB.PostgresMock.ExpectQuery(`SELECT (.+) FROM users ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
		WithArgs(perPage, 0).
		WillReturnRows(userRows)

	// Act
	result, total, err := userRepo.GetAllUsers(page, perPage)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, totalCount, total)
	assert.Equal(t, users[0].Username, result[0].Username)
	assert.Equal(t, users[1].Username, result[1].Username)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_CreateStudentProfile_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewAchievementFixtures()

	userID := uuid.New()
	advisorID := uuid.New()
	req := model.CreateStudentProfileRequest{
		UserID:          userID,
		StudentIDNumber: "2021001",
		ProgramStudy:    "Teknik Informatika",
		AcademicYear:    "2021",
		AdvisorID:       &advisorID,
	}

	expectedStudent := fixtures.ValidStudent()
	expectedStudent.UserID = userID
	expectedStudent.StudentID = req.StudentIDNumber
	expectedStudent.AdvisorID = &advisorID

	// Setup mock expectations
	studentRows := sqlmock.NewRows([]string{
		"id", "user_id", "student_id", "program_study", "academic_year", "advisor_id", "created_at",
	}).AddRow(
		expectedStudent.ID, expectedStudent.UserID, expectedStudent.StudentID,
		expectedStudent.ProgramStudy, expectedStudent.AcademicYear, expectedStudent.AdvisorID, expectedStudent.CreatedAt,
	)

	mockDB.PostgresMock.ExpectQuery(`INSERT INTO students`).
		WithArgs(sqlmock.AnyArg(), userID, req.StudentIDNumber, req.ProgramStudy, req.AcademicYear, advisorID, sqlmock.AnyArg()).
		WillReturnRows(studentRows)

	// Act
	result, err := userRepo.CreateStudentProfile(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, req.StudentIDNumber, result.StudentID)
	assert.Equal(t, req.ProgramStudy, result.ProgramStudy)
	assert.Equal(t, req.AcademicYear, result.AcademicYear)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_GetAllRoles_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewUserFixtures()

	expectedRoles := fixtures.MultipleRoles()

	// Setup mock expectations
	roleRows := sqlmock.NewRows([]string{
		"id", "name", "description", "created_at",
	})
	for _, role := range expectedRoles {
		roleRows.AddRow(role.ID, role.Name, role.Description, role.CreatedAt)
	}

	mockDB.PostgresMock.ExpectQuery(`SELECT (.+) FROM roles ORDER BY name`).
		WillReturnRows(roleRows)

	// Act
	result, err := userRepo.GetAllRoles()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, len(expectedRoles))
	assert.Equal(t, expectedRoles[0].Name, result[0].Name)
	assert.Equal(t, expectedRoles[1].Name, result[1].Name)
	assert.Equal(t, expectedRoles[2].Name, result[2].Name)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserRepository_UpdateStudentAdvisor_Success(t *testing.T) {
	// Arrange
	mockDB, err := mocks.NewMockDatabase()
	assert.NoError(t, err)
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	studentID := uuid.New()
	advisorID := uuid.New()

	// Setup mock expectations
	mockDB.PostgresMock.ExpectExec(`UPDATE students SET advisor_id = \$1 WHERE id = \$2`).
		WithArgs(advisorID, studentID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Act
	err = userRepo.UpdateStudentAdvisor(studentID, advisorID)

	// Assert
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkUserRepository_CreateUser(b *testing.B) {
	mockDB, _ := mocks.NewMockDatabase()
	defer mockDB.Close()

	userRepo := repository.NewUserRepository(mockDB.PostgresDB)
	fixtures := fixtures.NewUserFixtures()
	req := fixtures.ValidCreateUserRequest()

	// Setup mock for benchmark
	mockDB.PostgresMock.ExpectQuery(`INSERT INTO users`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())).
		Times(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userRepo.CreateUser(req)
	}
}

// Table-driven tests
func TestUserRepository_CreateUser_ValidationCases(t *testing.T) {
	testCases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		request     model.CreateUserRequest
		expectError bool
		errorType   error
	}{
		{
			name: "Success case",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
			request: model.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
				RoleID:   uuid.New(),
			},
			expectError: false,
		},
		{
			name: "Database constraint violation",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users`).
					WillReturnError(sql.ErrConnDone)
			},
			request: model.CreateUserRequest{
				Username: "duplicate",
				Email:    "duplicate@example.com",
				Password: "password123",
				FullName: "Duplicate User",
				RoleID:   uuid.New(),
			},
			expectError: true,
			errorType:   sql.ErrConnDone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockDB, err := mocks.NewMockDatabase()
			assert.NoError(t, err)
			defer mockDB.Close()

			userRepo := repository.NewUserRepository(mockDB.PostgresDB)
			tc.setupMock(mockDB.PostgresMock)

			// Act
			result, err := userRepo.CreateUser(tc.request)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tc.errorType != nil {
					assert.Equal(t, tc.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			// Verify expectations
			assert.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}