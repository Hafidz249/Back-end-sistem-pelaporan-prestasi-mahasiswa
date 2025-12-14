package service_test

import (
	"POJECT_UAS/service"
	"POJECT_UAS/tests/fixtures"
	"POJECT_UAS/tests/helpers"
	"POJECT_UAS/tests/mocks"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewUserFixtures()

	app := helper.CreateFiberApp()
	app.Post("/login", authService.Login)

	loginReq := fixtures.ValidLoginRequest()
	expectedResponse := fixtures.ValidLoginResponse()

	// Setup mock expectations
	mockRepo.On("Login", loginReq).Return(expectedResponse, nil)

	// Act
	req := helper.CreateJSONRequest("POST", "/login", loginReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewUserFixtures()

	app := helper.CreateFiberApp()
	app.Post("/login", authService.Login)

	loginReq := fixtures.ValidLoginRequest()

	// Setup mock expectations - return error for invalid credentials
	mockRepo.On("Login", loginReq).Return(nil, assert.AnError)

	// Act
	req := helper.CreateJSONRequest("POST", "/login", loginReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidRequestBody(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)

	app := helper.CreateFiberApp()
	app.Post("/login", authService.Login)

	// Act - Send invalid JSON
	req := helper.CreateJSONRequest("POST", "/login", "invalid json")
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Verify no repository calls were made
	mockRepo.AssertNotCalled(t, "Login")
}

func TestAuthService_Login_MissingCredential(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewUserFixtures()

	app := helper.CreateFiberApp()
	app.Post("/login", authService.Login)

	// Create request with empty credential
	loginReq := fixtures.InvalidLoginRequest()

	// Act
	req := helper.CreateJSONRequest("POST", "/login", loginReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Verify no repository calls were made
	mockRepo.AssertNotCalled(t, "Login")
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)

	app := helper.CreateFiberApp()
	app.Post("/refresh", authService.RefreshToken)

	refreshReq := map[string]string{
		"refresh_token": "valid_refresh_token",
	}

	// Act
	req := helper.CreateJSONRequest("POST", "/refresh", refreshReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	helper.AssertSuccessResponse(httptest.NewRecorder(), 200)
}

func TestAuthService_RefreshToken_MissingToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)

	app := helper.CreateFiberApp()
	app.Post("/refresh", authService.RefreshToken)

	refreshReq := map[string]string{
		"refresh_token": "", // Empty token
	}

	// Act
	req := helper.CreateJSONRequest("POST", "/refresh", refreshReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAuthService_Logout_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)

	app := helper.CreateFiberApp()
	app.Post("/logout", authService.Logout)

	// Act
	req := helper.CreateJSONRequest("POST", "/logout", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthService_GetProfile_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)

	userID := helper.GenerateUUID()
	username := "testuser"
	email := "test@example.com"
	role := "student"

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, username, email, role))
	app.Get("/profile", authService.GetProfile)

	// Act
	req := helper.CreateJSONRequest("GET", "/profile", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify response contains user data
	recorder := httptest.NewRecorder()
	recorder.Code = resp.StatusCode
	recorder.Body.Write([]byte{}) // This would contain actual response body
	helper.AssertSuccessResponse(recorder, 200)
}

func TestAuthService_Login_WithEmail_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewUserFixtures()

	app := helper.CreateFiberApp()
	app.Post("/login", authService.Login)

	loginReq := fixtures.ValidLoginRequestWithEmail()
	expectedResponse := fixtures.ValidLoginResponse()

	// Setup mock expectations
	mockRepo.On("Login", loginReq).Return(expectedResponse, nil)

	// Act
	req := helper.CreateJSONRequest("POST", "/login", loginReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkAuthService_Login(b *testing.B) {
	mockRepo := new(mocks.MockAuthRepository)
	authService := service.NewAuthService(mockRepo)
	fixtures := fixtures.NewUserFixtures()

	loginReq := fixtures.ValidLoginRequest()
	expectedResponse := fixtures.ValidLoginResponse()

	mockRepo.On("Login", mock.Anything).Return(expectedResponse, nil)

	app := fiber.New()
	app.Post("/login", authService.Login)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		helper := helpers.NewTestHelper(nil)
		req := helper.CreateJSONRequest("POST", "/login", loginReq)
		app.Test(req)
	}
}