package service

import (
	"POJECT_UAS/model"
	"POJECT_UAS/repository"
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

// MockAuthRepository adalah mock untuk AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Login(username, password string) (*model.LoginResponse, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoginResponse), args.Error(1)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	authService := &AuthService{AuthRepo: mockRepo}

	// Mock data
	expectedResponse := &model.LoginResponse{
		Token: "mock-jwt-token",
		User: model.UserProfile{
			ID:       uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     "student",
		},
	}

	// Mock expectations
	mockRepo.On("Login", "testuser", "password123").Return(expectedResponse, nil)

	// Create Fiber app and request
	app := fiber.New()
	app.Post("/login", authService.Login)

	loginReq := model.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	reqBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	authService := &AuthService{AuthRepo: mockRepo}

	// Mock expectations - return error for invalid credentials
	mockRepo.On("Login", "wronguser", "wrongpass").Return(nil, assert.AnError)

	// Create Fiber app and request
	app := fiber.New()
	app.Post("/login", authService.Login)

	loginReq := model.LoginRequest{
		Username: "wronguser",
		Password: "wrongpass",
	}
	reqBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidRequestBody(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	authService := &AuthService{AuthRepo: mockRepo}

	// Create Fiber app and request
	app := fiber.New()
	app.Post("/login", authService.Login)

	// Invalid JSON
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Verify no repository calls were made
	mockRepo.AssertNotCalled(t, "Login")
}

func TestAuthService_Login_MissingFields(t *testing.T) {
	// Setup
	mockRepo := new(MockAuthRepository)
	authService := &AuthService{AuthRepo: mockRepo}

	// Create Fiber app and request
	app := fiber.New()
	app.Post("/login", authService.Login)

	// Missing password
	loginReq := model.LoginRequest{
		Username: "testuser",
		Password: "", // empty password
	}
	reqBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Verify no repository calls were made
	mockRepo.AssertNotCalled(t, "Login")
}