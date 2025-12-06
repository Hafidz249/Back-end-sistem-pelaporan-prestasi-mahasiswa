package service

import (
	"POJECT_UAS/middleware"
	"POJECT_UAS/model"
	"POJECT_UAS/repository"
	"encoding/json"
	"fmt"
	"time"

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

// SubmitForVerification - Mahasiswa submit prestasi draft untuk diverifikasi (FR-004)
func (s *AchievementService) SubmitForVerification(c *fiber.Ctx) error {
	referenceIDStr := c.Params("reference_id")

	if referenceIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "achievement reference id required",
		})
	}

	referenceID, err := uuid.Parse(referenceIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid achievement reference id",
		})
	}

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

	// Ambil achievement reference
	achievementRef, err := s.AchievementRepo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "achievement reference not found",
		})
	}

	// Verify ownership
	if achievementRef.StudentID != student.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you can only submit your own achievements",
		})
	}

	// Verify status is draft
	if achievementRef.Status != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only draft achievements can be submitted for verification",
		})
	}

	// Update status menjadi 'submitted'
	err = s.AchievementRepo.SubmitForVerification(referenceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to submit achievement for verification",
		})
	}

	// Ambil achievement detail untuk notifikasi
	achievement, err := s.AchievementRepo.GetAchievementByID(achievementRef.MongoAchievementID)
	if err != nil {
		// Log error tapi tetap lanjut (notifikasi optional)
		achievement = &model.Achievement{
			Title: "Achievement",
		}
	}

	// Create notification untuk dosen wali
	err = s.createNotificationForAdvisor(student, achievementRef, achievement)
	if err != nil {
		// Log error tapi tidak gagalkan request (notifikasi optional)
		// TODO: Add proper logging
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "prestasi berhasil disubmit untuk verifikasi",
		"data": model.SubmitForVerificationResponse{
			AchievementReferenceID: referenceID,
			Status:                 "submitted",
			SubmittedAt:            time.Now().Format(time.RFC3339),
			Message:                "Prestasi Anda telah disubmit dan menunggu verifikasi dari dosen wali",
		},
	})
}

// createNotificationForAdvisor membuat notifikasi untuk dosen wali
func (s *AchievementService) createNotificationForAdvisor(student *model.Student, achievementRef *model.AchievementReference, achievement *model.Achievement) error {
	// Ambil advisor_id
	advisorID, err := s.AchievementRepo.GetAdvisorByStudentID(student.ID)
	if err != nil {
		return err
	}

	// Ambil data student user untuk nama
	studentUser, err := s.AchievementRepo.GetUserByID(student.UserID)
	if err != nil {
		return err
	}

	// Buat notification data
	notifData := model.NotificationData{
		AchievementID:          achievementRef.MongoAchievementID,
		AchievementReferenceID: achievementRef.ID,
		StudentID:              student.ID,
		StudentName:            studentUser.FullName,
		AchievementTitle:       achievement.Title,
	}

	// Convert to JSON string
	dataJSON, err := json.Marshal(notifData)
	if err != nil {
		return err
	}

	// Create notification
	notification := model.Notification{
		ID:        uuid.New(),
		UserID:    advisorID,
		Type:      "achievement_submitted",
		Title:     "Prestasi Baru Menunggu Verifikasi",
		Message:   fmt.Sprintf("Mahasiswa %s telah mengajukan prestasi '%s' untuk diverifikasi", studentUser.FullName, achievement.Title),
		Data:      string(dataJSON),
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.AchievementRepo.CreateNotification(notification)
}
