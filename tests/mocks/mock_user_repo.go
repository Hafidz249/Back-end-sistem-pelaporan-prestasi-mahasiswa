package mocks

import (
	"POJECT_UAS/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah mock implementation untuk UserRepository
type MockUserRepository struct {
	mock.Mock
}

// CreateUser mocks user creation
func (m *MockUserRepository) CreateUser(req model.CreateUserRequest) (*model.Users, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

// UpdateUser mocks user update
func (m *MockUserRepository) UpdateUser(userID uuid.UUID, req model.UpdateUserRequest) (*model.Users, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

// DeleteUser mocks user deletion
func (m *MockUserRepository) DeleteUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// GetAllUsers mocks getting all users with pagination
func (m *MockUserRepository) GetAllUsers(page, perPage int) ([]model.Users, int64, error) {
	args := m.Called(page, perPage)
	return args.Get(0).([]model.Users), args.Get(1).(int64), args.Error(2)
}

// GetUserByID mocks getting user by ID
func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*model.Users, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Users), args.Error(1)
}

// CreateStudentProfile mocks student profile creation
func (m *MockUserRepository) CreateStudentProfile(req model.CreateStudentProfileRequest) (*model.Student, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Student), args.Error(1)
}

// CreateLecturerProfile mocks lecturer profile creation
func (m *MockUserRepository) CreateLecturerProfile(req model.CreateLecturerProfileRequest) (*model.Lecturers, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Lecturers), args.Error(1)
}

// UpdateStudentAdvisor mocks updating student advisor
func (m *MockUserRepository) UpdateStudentAdvisor(studentID, advisorID uuid.UUID) error {
	args := m.Called(studentID, advisorID)
	return args.Error(0)
}

// GetAllRoles mocks getting all roles
func (m *MockUserRepository) GetAllRoles() ([]model.Roles, error) {
	args := m.Called()
	return args.Get(0).([]model.Roles), args.Error(1)
}

// UpdateUserRole mocks updating user role
func (m *MockUserRepository) UpdateUserRole(userID, roleID uuid.UUID) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

// GetAllStudents mocks getting all students
func (m *MockUserRepository) GetAllStudents(page, perPage int) ([]model.Student, int64, error) {
	args := m.Called(page, perPage)
	return args.Get(0).([]model.Student), args.Get(1).(int64), args.Error(2)
}

// GetAllLecturers mocks getting all lecturers
func (m *MockUserRepository) GetAllLecturers(page, perPage int) ([]model.Lecturers, int64, error) {
	args := m.Called(page, perPage)
	return args.Get(0).([]model.Lecturers), args.Get(1).(int64), args.Error(2)
}