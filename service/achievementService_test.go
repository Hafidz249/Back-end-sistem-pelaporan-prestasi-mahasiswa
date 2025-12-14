package service

import (
	"POJECT_UAS/model"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAchievementRepository adalah mock untuk AchievementRepository
type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) SubmitAchievement(studentID uuid.UUID, req model.SubmitAchievementRequest) (*model.SubmitAchievementResponse, error) {
	args := m.Called(studentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SubmitAchievementResponse), args.Error(1)
}

func (m *MockAchievementRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockAchievementRepository) GetAchievementReferenceByID(referenceID uuid.UUID) (*model.AchievementReference, error) {
	args := m.Called(referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) SubmitForVerification(referenceID uuid.UUID) error {
	args := m.Called(referenceID)
	return args.Error(0)
}

func (m *MockAchievementRepository) CreateNotification(notification model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockAchievementRepository) GetAdvisorByStudentID(studentID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(studentID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAchievementRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *MockAchievementRepository) DeleteAchievement(referenceID uuid.UUID, mongoAchievementID string) error {
	args := m.Called(referenceID, mongoAchievementID)
	return args.Error(0)
}

func TestAchievementService_SubmitAchievement_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAchievementRepository)
	achievementService := &AchievementService{AchievementRepo: mockRepo}

	// Mock data
	userID := uuid.New()
	studentID := uuid.New()
	
	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "2021001",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2021",
		CreatedAt:    time.Now(),
	}

	submitReq := model.SubmitAchievementRequest{
		AchievementType: "akademik",
		Title:           "Juara 1 Programming Contest",
		Description:     "Kompetisi programming tingkat nasional",
		Details: map[string]interface{}{
			"level":    "nasional",
			"category": "programming",
		},
	}

	expectedResponse := &model.SubmitAchievementResponse{
		AchievementID:          "507f1f77bcf86cd799439011",
		AchievementReferenceID: uuid.New(),
		StudentID:              studentID,
		AchievementType:        "akademik",
		Title:                  "Juara 1 Programming Contest",
		Description:            "Kompetisi programming tingkat nasional",
		Status:                 "draft",
		CreatedAt:              time.Now().Format(time.RFC3339),
	}

	// Mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("SubmitAchievement", studentID, submitReq).Return(expectedResponse, nil)

	// Create Fiber app and request
	app := fiber.New()
	
	// Middleware to set user_id in locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	
	app.Post("/achievements", achievementService.SubmitAchievement)

	reqBody, _ := json.Marshal(submitReq)
	req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitAchievement_StudentNotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockAchievementRepository)
	achievementService := &AchievementService{AchievementRepo: mockRepo}

	userID := uuid.New()

	// Mock expectations - student not found
	mockRepo.On("GetStudentByUserID", userID).Return(nil, assert.AnError)

	// Create Fiber app and request
	app := fiber.New()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	
	app.Post("/achievements", achievementService.SubmitAchievement)

	submitReq := model.SubmitAchievementRequest{
		AchievementType: "akademik",
		Title:           "Test Achievement",
		Description:     "Test Description",
	}
	reqBody, _ := json.Marshal(submitReq)

	req := httptest.NewRequest("POST", "/achievements", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitForVerification_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAchievementRepository)
	achievementService := &AchievementService{AchievementRepo: mockRepo}

	userID := uuid.New()
	studentID := uuid.New()
	referenceID := uuid.New()
	advisorID := uuid.New()

	student := &model.Student{
		ID:        studentID,
		UserID:    userID,
		AdvisorID: &advisorID,
	}

	achievementRef := &model.AchievementReference{
		ID:        referenceID,
		StudentID: studentID,
		Status:    "draft",
	}

	// Mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(achievementRef, nil)
	mockRepo.On("SubmitForVerification", referenceID).Return(nil)
	mockRepo.On("GetAdvisorByStudentID", studentID).Return(advisorID, nil)
	mockRepo.On("CreateNotification", mock.AnythingOfType("model.Notification")).Return(nil)

	// Create Fiber app and request
	app := fiber.New()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	
	app.Post("/achievements/:reference_id/submit", achievementService.SubmitForVerification)

	req := httptest.NewRequest("POST", "/achievements/"+referenceID.String()+"/submit", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_DeleteAchievement_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAchievementRepository)
	achievementService := &AchievementService{AchievementRepo: mockRepo}

	userID := uuid.New()
	studentID := uuid.New()
	referenceID := uuid.New()
	mongoID := "507f1f77bcf86cd799439011"

	student := &model.Student{
		ID:     studentID,
		UserID: userID,
	}

	achievementRef := &model.AchievementReference{
		ID:                 referenceID,
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
	}

	// Mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(achievementRef, nil)
	mockRepo.On("DeleteAchievement", referenceID, mongoID).Return(nil)

	// Create Fiber app and request
	app := fiber.New()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	
	app.Delete("/achievements/:reference_id", achievementService.DeleteAchievement)

	req := httptest.NewRequest("DELETE", "/achievements/"+referenceID.String(), nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}