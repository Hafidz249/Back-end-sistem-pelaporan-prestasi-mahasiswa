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
