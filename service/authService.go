package service

import (
	"POJECT_UAS/model"
	"POJECT_UAS/repository"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	AuthRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,
	}
}

// Login handler - menggabungkan handler dan service logic
// @Summary User login
// @Description Login dengan username/email dan password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse "Login berhasil"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi input
	if req.Credential == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username atau email harus diisi",
		})
	}
	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password harus diisi",
		})
	}

	// Proses login
	response, err := s.AuthRepo.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login berhasil",
		"data":    response,
	})
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
// RefreshToken - Refresh JWT token
// @Summary Refresh JWT token
// @Description Refresh expired JWT token dengan refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} map[string]interface{} "Token refreshed"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid refresh token"
// @Router /api/v1/auth/refresh [post]
func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh_token is required",
		})
	}

	// TODO: Implement refresh token logic
	// For now, return placeholder response
	return c.JSON(fiber.Map{
		"message": "token refreshed successfully",
		"data": fiber.Map{
			"token": "new-jwt-token",
		},
	})
}

// Logout - User logout
// @Summary User logout
// @Description Logout user dan invalidate token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/auth/logout [post]
func (s *AuthService) Logout(c *fiber.Ctx) error {
	// TODO: Implement token blacklisting
	// For now, return success response
	return c.JSON(fiber.Map{
		"message": "logout successful",
	})
}

// GetProfile - Get user profile
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/auth/profile [get]
func (s *AuthService) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	username := c.Locals("username")
	email := c.Locals("email")
	role := c.Locals("role")

	return c.JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"email":    email,
			"role":     role,
		},
	})
}