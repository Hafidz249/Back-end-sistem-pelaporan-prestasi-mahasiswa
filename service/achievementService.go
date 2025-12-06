package service

import (
	"POJECT_UAS/middleware"
	"POJECT_UAS/model"
	"POJECT_UAS/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementService struct {
	AchievementRepo *repository.AchievementRepository
}

func NewAchievementService(achievementRepo *repository.AchievementRepository) *AchievementService {
	return &AchievementService{
		AchievementRepo: achievementRepo,
	}
}

// SubmitAchievement - Mahasiswa submit prestasi
func (s *AchievementService) SubmitAchievement(c *fiber.Ctx) error {
	var req model.SubmitAchievementRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi input
	if err := s.validateSubmitRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil user_id dari context (dari JWT)
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	// Ambil student_id dari user_id
	student, err := s.AchievementRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	// Submit achievement
	response, err := s.AchievementRepo.SubmitAchievement(student.ID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to submit achievement",
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "prestasi berhasil disubmit",
		"data":    response,
	})
}

// validateSubmitRequest validasi input submit prestasi
func (s *AchievementService) validateSubmitRequest(req model.SubmitAchievementRequest) error {
	if req.AchievementType == "" {
		return &ValidationError{Message: "achievement_type harus diisi"}
	}

	// Validasi achievement_type
	validTypes := map[string]bool{
		"academic":      true,
		"competition":   true,
		"organization":  true,
		"publication":   true,
		"certification": true,
		"other":         true,
	}

	if !validTypes[req.AchievementType] {
		return &ValidationError{Message: "achievement_type tidak valid"}
	}

	if req.Title == "" {
		return &ValidationError{Message: "title harus diisi"}
	}

	if req.Description == "" {
		return &ValidationError{Message: "description harus diisi"}
	}

	return nil
}

// GetMyAchievements - Mahasiswa melihat prestasi sendiri
func (s *AchievementService) GetMyAchievements(c *fiber.Ctx) error {
	// Ambil user_id dari context
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	// Ambil student_id dari user_id
	student, err := s.AchievementRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	// Ambil semua achievements
	achievements, err := s.AchievementRepo.GetAchievementsByStudentID(student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get achievements",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    achievements,
	})
}

// GetAchievementDetail - Melihat detail prestasi
func (s *AchievementService) GetAchievementDetail(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	if achievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "achievement id required",
		})
	}

	// Ambil achievement dari MongoDB
	achievement, err := s.AchievementRepo.GetAchievementByID(achievementID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "achievement not found",
		})
	}

	// Verify ownership: pastikan achievement milik user yang login
	userIDStr := middleware.GetUserID(c)
	userID, _ := uuid.Parse(userIDStr)

	student, err := s.AchievementRepo.GetStudentByUserID(userID)
	if err == nil && achievement.StudentID != student.ID {
		// Jika bukan milik sendiri, cek apakah user punya permission read_all
		if !middleware.HasPermission(c, "achievements", "read_all") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you can only view your own achievements",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    achievement,
	})
}
