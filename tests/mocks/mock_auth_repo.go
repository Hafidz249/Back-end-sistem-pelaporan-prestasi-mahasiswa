package mocks

import (
	"POJECT_UAS/model"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepository adalah mock implementation untuk AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

// Login mocks the Login method
func (m *MockAuthRepository) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoginResponse), args.Error(1)
}

// ValidateToken mocks token validation
func (m *MockAuthRepository) ValidateToken(token string) (*model.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenClaims), args.Error(1)
}

// RefreshToken mocks token refresh
func (m *MockAuthRepository) RefreshToken(refreshToken string) (*model.LoginResponse, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoginResponse), args.Error(1)
}

// InvalidateToken mocks token invalidation
func (m *MockAuthRepository) InvalidateToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}