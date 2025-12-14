package mocks

import (
	"POJECT_UAS/model"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAchievementRepository adalah mock implementation untuk AchievementRepository
type MockAchievementRepository struct {
	mock.Mock
}

// SubmitAchievement mocks achievement submission
func (m *MockAchievementRepository) SubmitAchievement(studentID uuid.UUID, req model.SubmitAchievementRequest) (*model.SubmitAchievementResponse, error) {
	args := m.Called(studentID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SubmitAchievementResponse), args.Error(1)
}

// GetStudentByUserID mocks getting student by user ID
func (m *MockAchievementRepository) GetStudentByUserID(userID uuid.UUID) (*model.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

// GetAchievementByID mocks getting achievement by ID
func (m *MockAchievementRepository) GetAchievementByID(achievementID string) (*model.Achievement, error) {
	args := m.Called(achievementID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Achievement), args.Error(1)
}

// GetAchievementReferenceByID mocks getting achievement reference
func (m *MockAchievementRepository) GetAchievementReferenceByID(referenceID uuid.UUID) (*model.AchievementReference, error) {
	args := m.Called(referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AchievementReference), args.Error(1)
}

// SubmitForVerification mocks submission for verification
func (m *MockAchievementRepository) SubmitForVerification(referenceID uuid.UUID) error {
	args := m.Called(referenceID)
	return args.Error(0)
}

// CreateNotification mocks notification creation
func (m *MockAchievementRepository) CreateNotification(notification model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

// GetAdvisorByStudentID mocks getting advisor by student ID
func (m *MockAchievementRepository) GetAdvisorByStudentID(studentID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(studentID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// GetUserByID mocks getting user by ID
func (m *MockAchievementRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

// DeleteAchievement mocks achievement deletion
func (m *MockAchievementRepository) DeleteAchievement(referenceID uuid.UUID, mongoAchievementID string) error {
	args := m.Called(referenceID, mongoAchievementID)
	return args.Error(0)
}

// GetLecturerByUserID mocks getting lecturer by user ID
func (m *MockAchievementRepository) GetLecturerByUserID(userID uuid.UUID) (*model.Lecturers, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturers), args.Error(1)
}

// GetStudentIDsByAdvisor mocks getting student IDs by advisor
func (m *MockAchievementRepository) GetStudentIDsByAdvisor(advisorID uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(advisorID)
	return args.Get(0).([]uuid.UUID), args.Error(1)
}

// GetAchievementReferencesWithPagination mocks paginated achievement references
func (m *MockAchievementRepository) GetAchievementReferencesWithPagination(
	studentIDs []uuid.UUID,
	status string,
	page int,
	perPage int,
) ([]model.AchievementReference, int64, error) {
	args := m.Called(studentIDs, status, page, perPage)
	return args.Get(0).([]model.AchievementReference), args.Get(1).(int64), args.Error(2)
}

// GetStudentByID mocks getting student by ID
func (m *MockAchievementRepository) GetStudentByID(studentID uuid.UUID) (*model.Student, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

// GetAchievementsByIDs mocks getting multiple achievements by IDs
func (m *MockAchievementRepository) GetAchievementsByIDs(achievementIDs []string) (map[string]model.Achievement, error) {
	args := m.Called(achievementIDs)
	return args.Get(0).(map[string]model.Achievement), args.Error(1)
}

// VerifyAchievement mocks achievement verification
func (m *MockAchievementRepository) VerifyAchievement(referenceID uuid.UUID, verifiedBy uuid.UUID) error {
	args := m.Called(referenceID, verifiedBy)
	return args.Error(0)
}

// RejectAchievement mocks achievement rejection
func (m *MockAchievementRepository) RejectAchievement(referenceID uuid.UUID, verifiedBy uuid.UUID, rejectionNote string) error {
	args := m.Called(referenceID, verifiedBy, rejectionNote)
	return args.Error(0)
}

// CheckLecturerOwnsStudent mocks checking lecturer ownership
func (m *MockAchievementRepository) CheckLecturerOwnsStudent(lecturerID uuid.UUID, studentID uuid.UUID) (bool, error) {
	args := m.Called(lecturerID, studentID)
	return args.Bool(0), args.Error(1)
}

// GetAllAchievementReferencesWithPagination mocks getting all achievement references
func (m *MockAchievementRepository) GetAllAchievementReferencesWithPagination(
	status string,
	achievementType string,
	page int,
	perPage int,
) ([]model.AchievementReference, int64, error) {
	args := m.Called(status, achievementType, page, perPage)
	return args.Get(0).([]model.AchievementReference), args.Get(1).(int64), args.Error(2)
}

// GetAchievementStatistics mocks getting achievement statistics
func (m *MockAchievementRepository) GetAchievementStatistics(
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