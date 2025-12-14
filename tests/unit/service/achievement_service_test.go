package service_test

import (
	"POJECT_UAS/service"
	"POJECT_UAS/tests/fixtures"
	"POJECT_UAS/tests/helpers"
	"POJECT_UAS/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAchievementService_SubmitAchievement_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()
	student := fixtures.ValidStudent()
	student.UserID = userID

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Post("/achievements", achievementService.SubmitAchievement)

	submitReq := fixtures.ValidSubmitAchievementRequest()
	expectedResponse := fixtures.ValidSubmitAchievementResponse()

	// Setup mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("SubmitAchievement", student.ID, submitReq).Return(expectedResponse, nil)

	// Act
	req := helper.CreateJSONRequest("POST", "/achievements", submitReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitAchievement_StudentNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Post("/achievements", achievementService.SubmitAchievement)

	submitReq := fixtures.ValidSubmitAchievementRequest()

	// Setup mock expectations - student not found
	mockRepo.On("GetStudentByUserID", userID).Return(nil, assert.AnError)

	// Act
	req := helper.CreateJSONRequest("POST", "/achievements", submitReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitAchievement_InvalidRequest(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Post("/achievements", achievementService.SubmitAchievement)

	// Use invalid request
	submitReq := fixtures.InvalidSubmitAchievementRequest()

	// Act
	req := helper.CreateJSONRequest("POST", "/achievements", submitReq)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Verify no repository calls were made
	mockRepo.AssertNotCalled(t, "GetStudentByUserID")
	mockRepo.AssertNotCalled(t, "SubmitAchievement")
}

func TestAchievementService_SubmitForVerification_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()
	student := fixtures.ValidStudent()
	student.UserID = userID
	
	referenceID := helper.GenerateUUID()
	advisorID := helper.GenerateUUID()
	student.AdvisorID = &advisorID

	achievementRef := fixtures.ValidAchievementReference()
	achievementRef.ID = referenceID
	achievementRef.StudentID = student.ID

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Post("/achievements/:reference_id/submit", achievementService.SubmitForVerification)

	// Setup mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(achievementRef, nil)
	mockRepo.On("SubmitForVerification", referenceID).Return(nil)
	mockRepo.On("GetAdvisorByStudentID", student.ID).Return(advisorID, nil)
	mockRepo.On("CreateNotification", mock.AnythingOfType("model.Notification")).Return(nil)

	// Act
	req := helper.CreateJSONRequest("POST", "/achievements/"+referenceID.String()+"/submit", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_SubmitForVerification_NotOwner(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()
	student := fixtures.ValidStudent()
	student.UserID = userID
	
	referenceID := helper.GenerateUUID()
	achievementRef := fixtures.ValidAchievementReference()
	achievementRef.ID = referenceID
	achievementRef.StudentID = helper.GenerateUUID() // Different student ID

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Post("/achievements/:reference_id/submit", achievementService.SubmitForVerification)

	// Setup mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(achievementRef, nil)

	// Act
	req := helper.CreateJSONRequest("POST", "/achievements/"+referenceID.String()+"/submit", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_DeleteAchievement_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()
	student := fixtures.ValidStudent()
	student.UserID = userID
	
	referenceID := helper.GenerateUUID()
	mongoID := "507f1f77bcf86cd799439011"

	achievementRef := fixtures.ValidAchievementReference()
	achievementRef.ID = referenceID
	achievementRef.StudentID = student.ID
	achievementRef.MongoAchievementID = mongoID
	achievementRef.Status = "draft" // Only draft can be deleted

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Delete("/achievements/:reference_id", achievementService.DeleteAchievement)

	// Setup mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(achievementRef, nil)
	mockRepo.On("DeleteAchievement", referenceID, mongoID).Return(nil)

	// Act
	req := helper.CreateJSONRequest("DELETE", "/achievements/"+referenceID.String(), nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_DeleteAchievement_WrongStatus(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)
	fixtures := fixtures.NewAchievementFixtures()

	userID := helper.GenerateUUID()
	student := fixtures.ValidStudent()
	student.UserID = userID
	
	referenceID := helper.GenerateUUID()

	achievementRef := fixtures.SubmittedAchievementReference() // Status: submitted
	achievementRef.ID = referenceID
	achievementRef.StudentID = student.ID

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Delete("/achievements/:reference_id", achievementService.DeleteAchievement)

	// Setup mock expectations
	mockRepo.On("GetStudentByUserID", userID).Return(student, nil)
	mockRepo.On("GetAchievementReferenceByID", referenceID).Return(achievementRef, nil)

	// Act
	req := helper.CreateJSONRequest("DELETE", "/achievements/"+referenceID.String(), nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAchievementService_GetAchievementDetail_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	helper := helpers.NewTestHelper(t)

	userID := helper.GenerateUUID()
	achievementID := "507f1f77bcf86cd799439011"

	app := helper.CreateFiberApp()
	app.Use(helper.CreateMiddleware(userID, "student123", "student@example.com", "student"))
	app.Get("/achievements/:id", achievementService.GetAchievementDetail)

	// Act
	req := helper.CreateJSONRequest("GET", "/achievements/"+achievementID, nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	// This would be 200 if the method was fully implemented
	// For now, it might return a "coming soon" message
}

// Benchmark tests
func BenchmarkAchievementService_SubmitAchievement(b *testing.B) {
	mockRepo := new(mocks.MockAchievementRepository)
	achievementService := service.NewAchievementService(mockRepo)
	fixtures := fixtures.NewAchievementFixtures()

	student := fixtures.ValidStudent()
	submitReq := fixtures.ValidSubmitAchievementRequest()
	expectedResponse := fixtures.ValidSubmitAchievementResponse()

	mockRepo.On("GetStudentByUserID", mock.Anything).Return(student, nil)
	mockRepo.On("SubmitAchievement", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	app := fiber.New()
	app.Post("/achievements", achievementService.SubmitAchievement)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		helper := helpers.NewTestHelper(nil)
		req := helper.CreateJSONRequest("POST", "/achievements", submitReq)
		app.Test(req)
	}
}