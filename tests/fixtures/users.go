package fixtures

import (
	"POJECT_UAS/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserFixtures provides test data for users
type UserFixtures struct{}

// NewUserFixtures creates a new UserFixtures instance
func NewUserFixtures() *UserFixtures {
	return &UserFixtures{}
}

// ValidUser returns a valid user for testing
func (f *UserFixtures) ValidUser() *model.Users {
	return &model.Users{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$12$hashedpassword",
		FullName:     "Test User",
		RoleID:       uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// AdminUser returns an admin user for testing
func (f *UserFixtures) AdminUser() *model.Users {
	return &model.Users{
		ID:           uuid.New(),
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: "$2a$12$hashedpassword",
		FullName:     "Admin User",
		RoleID:       uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// StudentUser returns a student user for testing
func (f *UserFixtures) StudentUser() *model.Users {
	return &model.Users{
		ID:           uuid.New(),
		Username:     "student123",
		Email:        "student@example.com",
		PasswordHash: "$2a$12$hashedpassword",
		FullName:     "Student User",
		RoleID:       uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// LecturerUser returns a lecturer user for testing
func (f *UserFixtures) LecturerUser() *model.Users {
	return &model.Users{
		ID:           uuid.New(),
		Username:     "lecturer123",
		Email:        "lecturer@example.com",
		PasswordHash: "$2a$12$hashedpassword",
		FullName:     "Lecturer User",
		RoleID:       uuid.New(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// InactiveUser returns an inactive user for testing
func (f *UserFixtures) InactiveUser() *model.Users {
	return &model.Users{
		ID:           uuid.New(),
		Username:     "inactive",
		Email:        "inactive@example.com",
		PasswordHash: "$2a$12$hashedpassword",
		FullName:     "Inactive User",
		RoleID:       uuid.New(),
		IsActive:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// ValidCreateUserRequest returns a valid create user request
func (f *UserFixtures) ValidCreateUserRequest() model.CreateUserRequest {
	return model.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleID:   uuid.New(),
	}
}

// InvalidCreateUserRequest returns an invalid create user request
func (f *UserFixtures) InvalidCreateUserRequest() model.CreateUserRequest {
	return model.CreateUserRequest{
		Username: "", // Invalid: empty username
		Email:    "invalid-email", // Invalid: bad email format
		Password: "123", // Invalid: too short
		FullName: "",
		RoleID:   uuid.Nil, // Invalid: nil UUID
	}
}

// ValidUpdateUserRequest returns a valid update user request
func (f *UserFixtures) ValidUpdateUserRequest() model.UpdateUserRequest {
	return model.UpdateUserRequest{
		FullName: "Updated User Name",
		RoleID:   uuid.New(),
		IsActive: true,
	}
}

// ValidLoginRequest returns a valid login request
func (f *UserFixtures) ValidLoginRequest() model.LoginRequest {
	return model.LoginRequest{
		Credential: "testuser",
		Password:   "password123",
	}
}

// ValidLoginRequestWithEmail returns a valid login request with email
func (f *UserFixtures) ValidLoginRequestWithEmail() model.LoginRequest {
	return model.LoginRequest{
		Credential: "test@example.com",
		Password:   "password123",
	}
}

// InvalidLoginRequest returns an invalid login request
func (f *UserFixtures) InvalidLoginRequest() model.LoginRequest {
	return model.LoginRequest{
		Credential: "", // Invalid: empty credential
		Password:   "", // Invalid: empty password
	}
}

// ValidLoginResponse returns a valid login response
func (f *UserFixtures) ValidLoginResponse() *model.LoginResponse {
	return &model.LoginResponse{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token",
		User: model.UserProfile{
			ID:       uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			FullName: "Test User",
			Role:     "student",
			Permissions: []string{
				"achievements:create",
				"achievements:read",
			},
		},
	}
}

// MultipleUsers returns multiple users for testing pagination
func (f *UserFixtures) MultipleUsers(count int) []model.Users {
	users := make([]model.Users, count)
	for i := 0; i < count; i++ {
		users[i] = model.Users{
			ID:           uuid.New(),
			Username:     fmt.Sprintf("user%d", i+1),
			Email:        fmt.Sprintf("user%d@example.com", i+1),
			PasswordHash: "$2a$12$hashedpassword",
			FullName:     fmt.Sprintf("User %d", i+1),
			RoleID:       uuid.New(),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}
	return users
}

// ValidRole returns a valid role for testing
func (f *UserFixtures) ValidRole() model.Roles {
	return model.Roles{
		ID:          uuid.New(),
		Name:        "student",
		Description: "Student role",
		CreatedAt:   time.Now(),
	}
}

// MultipleRoles returns multiple roles for testing
func (f *UserFixtures) MultipleRoles() []model.Roles {
	return []model.Roles{
		{
			ID:          uuid.New(),
			Name:        "admin",
			Description: "Administrator role",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "student",
			Description: "Student role",
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "lecturer",
			Description: "Lecturer role",
			CreatedAt:   time.Now(),
		},
	}
}