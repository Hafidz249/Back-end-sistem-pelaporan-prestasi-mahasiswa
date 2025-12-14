package service

import (
	"POJECT_UAS/model"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStatisticsRepository adalah mock untuk repository methods yang digunakan StatisticsService
type MockStatisticsRepository struct {
	mock.Mock
}

func (m *MockStatisticsRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStatisticsRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturers, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturers), args.Error(1)
}

func (m *MockStatisticsRepository) GetStudentIDsByAdvisor(advisorID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(advisorID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

func (m *MockStatisticsRepository) GetAchievementStatistics(
	studentIDs []uuid.UUID,
	startDate *time.Time,
	endDate *time.Time,
	achievementType *string,
	status *string,
) (*model.AchievementStatistics, error) {
	args := m.Called(studentIDs, startDate, endDate, achievementType, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementStatistics), args.Error(1)
}

func TestStatisticsService_GetMyStatistics_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockStatisticsRepository)
	statisticsService := &StatisticsService{AchievementRepo: mockRepo}

	// Mock data
	userID := uuid.New()
	studentID := uuid.New()

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "2021001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2021",
	}

	expectedStats := &model.AchievementStatistics{
		TotalByType: []model.TypeStatistic{
			{AchievementType: "akademik", Count: 5, Percentage: 71.4},
			{AchievementType: "non-akademik", Count: 2, Percentage: 28.6},
		},
		TotalByPeriod: []model.PeriodStatistic{
			{Period: "2024-01", Count: 3, Year: 2024, Month: &[]int{1}[0]},
			{Period: "2024-02", Count: 4, Year: 2024, Month: &[]int{2}[0]},
		},
		TopStudents: []model.TopStudent{
			{
				StudentID:       studentID,
				StudentIDNumber: "2021001",
				FullName:        "Test Student",
				TotalCount:      7,
				VerifiedCount:   5,
			},
		},
		CompetitionLevels: []model.LevelStatistic{
			{Level: "nasional", Count: 4, Percentage: 57.1},
			{Level: "regional", Count: 3, Percentage: 42.9},
		},
		Summary: model.StatisticSummary{
			TotalAchievements:    7,
			VerifiedAchievements: 5,
			PendingAchievements:  1,
			RejectedAchievements: 1,
			TotalStudents:        1,
		},
	}

	// Mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementStatistics", []uuid.UUID{studentID}, (*time.Time)(nil), (*time.Time)(nil), (*string)(nil), (*string)(nil)).Return(expectedStats, nil)

	// Create Fiber app and request
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})

	app.Get("/statistics", statisticsService.GetMyStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetMyStatistics_StudentNotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockStatisticsRepository)
	statisticsService := &StatisticsService{AchievementRepo: mockRepo}

	userID := uuid.New()

	// Mock expectations - student not found
	mockRepo.On("GetStudentByUserID", userID).Return(nil, assert.AnError)

	// Create Fiber app and request
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})

	app.Get("/statistics", statisticsService.GetMyStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAdviseeStatistics_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockStatisticsRepository)
	statisticsService := &StatisticsService{AchievementRepo: mockRepo}

	// Mock data
	userID := uuid.New()
	lecturerID := uuid.New()
	studentID1 := uuid.New()
	studentID2 := uuid.New()

	lecturer := &model.Lecturers{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "L001",
		Department: "Teknik Informatika",
	}

	studentIDs := []uuid.UUID{studentID1, studentID2}

	expectedStats := &model.AchievementStatistics{
		TotalByType: []model.TypeStatistic{
			{AchievementType: "akademik", Count: 8, Percentage: 66.7},
			{AchievementType: "non-akademik", Count: 4, Percentage: 33.3},
		},
		TopStudents: []model.TopStudent{
			{
				StudentID:       studentID1,
				StudentIDNumber: "2021001",
				FullName:        "Student 1",
				TotalCount:      7,
				VerifiedCount:   6,
			},
			{
				StudentID:       studentID2,
				StudentIDNumber: "2021002",
				FullName:        "Student 2",
				TotalCount:      5,
				VerifiedCount:   4,
			},
		},
		Summary: model.StatisticSummary{
			TotalAchievements:    12,
			VerifiedAchievements: 10,
			PendingAchievements:  1,
			RejectedAchievements: 1,
			TotalStudents:        2,
		},
	}

	// Mock expectations
	mockRepo.On("GetLecturerByUserID", userID).Return(lecturer, nil)
	mockRepo.On("GetStudentIDsByAdvisor", lecturerID).Return(studentIDs, nil)
	mockRepo.On("GetAchievementStatistics", studentIDs, (*time.Time)(nil), (*time.Time)(nil), (*string)(nil), (*string)(nil)).Return(expectedStats, nil)

	// Create Fiber app and request
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})

	app.Get("/statistics", statisticsService.GetAdviseeStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAdviseeStatistics_NoStudents(t *testing.T) {
	// Setup
	mockRepo := new(MockStatisticsRepository)
	statisticsService := &StatisticsService{AchievementRepo: mockRepo}

	// Mock data
	userID := uuid.New()
	lecturerID := uuid.New()

	lecturer := &model.Lecturers{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "L001",
		Department: "Teknik Informatika",
	}

	// Mock expectations - no students
	mockRepo.On("GetLecturerByUserID", userID).Return(lecturer, nil)
	mockRepo.On("GetStudentIDsByAdvisor", lecturerID).Return([]uuid.UUID{}, nil)

	// Create Fiber app and request
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})

	app.Get("/statistics", statisticsService.GetAdviseeStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAllStatistics_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockStatisticsRepository)
	statisticsService := &StatisticsService{AchievementRepo: mockRepo}

	expectedStats := &model.AchievementStatistics{
		TotalByType: []model.TypeStatistic{
			{AchievementType: "akademik", Count: 50, Percentage: 62.5},
			{AchievementType: "non-akademik", Count: 30, Percentage: 37.5},
		},
		TopStudents: []model.TopStudent{
			{
				StudentID:       uuid.New(),
				StudentIDNumber: "2021001",
				FullName:        "Top Student 1",
				TotalCount:      15,
				VerifiedCount:   12,
			},
			{
				StudentID:       uuid.New(),
				StudentIDNumber: "2021002",
				FullName:        "Top Student 2",
				TotalCount:      12,
				VerifiedCount:   10,
			},
		},
		Summary: model.StatisticSummary{
			TotalAchievements:    80,
			VerifiedAchievements: 65,
			PendingAchievements:  10,
			RejectedAchievements: 5,
			TotalStudents:        25,
		},
	}

	// Mock expectations - nil studentIDs for all students
	mockRepo.On("GetAchievementStatistics", ([]uuid.UUID)(nil), (*time.Time)(nil), (*time.Time)(nil), (*string)(nil), (*string)(nil)).Return(expectedStats, nil)

	// Create Fiber app and request
	app := fiber.New()
	app.Get("/statistics", statisticsService.GetAllStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestStatisticsService_GetAllStatistics_WithFilters(t *testing.T) {
	// Setup
	mockRepo := new(MockStatisticsRepository)
	statisticsService := &StatisticsService{AchievementRepo: mockRepo}

	// Parse expected dates
	startDate, _ := time.Parse("2006-01-02", "2024-01-01")
	endDate, _ := time.Parse("2006-01-02", "2024-12-31")
	achievementType := "akademik"
	status := "verified"

	expectedStats := &model.AchievementStatistics{
		Summary: model.StatisticSummary{
			TotalAchievements:    30,
			VerifiedAchievements: 30,
			TotalStudents:        15,
		},
	}

	// Mock expectations with filters
	mockRepo.On("GetAchievementStatistics", ([]uuid.UUID)(nil), &startDate, &endDate, &achievementType, &status).Return(expectedStats, nil)

	// Create Fiber app and request
	app := fiber.New()
	app.Get("/statistics", statisticsService.GetAllStatistics)

	req := httptest.NewRequest("GET", "/statistics?start_date=2024-01-01&end_date=2024-12-31&achievement_type=akademik&status=verified&top_limit=15", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}