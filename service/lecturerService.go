package service

import (
	"POJECT_UAS/middleware"
	"POJECT_UAS/model"
	"POJECT_UAS/repository"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService struct {
	AchievementRepo *repository.AchievementRepository
}

func NewLecturerService(achievementRepo *repository.AchievementRepository) *LecturerService {
	return &LecturerService{
		AchievementRepo: achievementRepo,
	}
}

// ViewStudentAchievements - Dosen view prestasi mahasiswa bimbingan (FR-006)
func (s *LecturerService) ViewStudentAchievements(c *fiber.Ctx) error {
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

	// Ambil lecturer_id dari user_id
	lecturer, err := s.AchievementRepo.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "lecturer not found",
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	status := c.Query("status", "") // Filter by status (optional)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// 1. Get list student IDs dari tabel students where advisor_id
	studentIDs, err := s.AchievementRepo.GetStudentIDsByAdvisor(lecturer.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get students",
		})
	}

	if len(studentIDs) == 0 {
		return c.JSON(fiber.Map{
			"message": "no students found",
			"data": model.PaginatedAchievementsResponse{
				Data: []model.AchievementWithStudent{},
				Pagination: model.PaginationMeta{
					CurrentPage: page,
					PerPage:     perPage,
					TotalPages:  0,
					TotalItems:  0,
				},
			},
		})
	}

	// 2. Get achievements references dengan filter student_ids
	references, totalCount, err := s.AchievementRepo.GetAchievementReferencesWithPagination(
		studentIDs,
		status,
		page,
		perPage,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get achievement references",
		})
	}

	if len(references) == 0 {
		totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
		return c.JSON(fiber.Map{
			"message": "no achievements found",
			"data": model.PaginatedAchievementsResponse{
				Data: []model.AchievementWithStudent{},
				Pagination: model.PaginationMeta{
					CurrentPage: page,
					PerPage:     perPage,
					TotalPages:  totalPages,
					TotalItems:  totalCount,
				},
			},
		})
	}

	// 3. Fetch detail dari MongoDB
	mongoIDs := make([]string, len(references))
	for i, ref := range references {
		mongoIDs[i] = ref.MongoAchievementID
	}

	achievementsMap, err := s.AchievementRepo.GetAchievementsByIDs(mongoIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get achievement details",
		})
	}

	// Get student info for each achievement
	studentsMap := make(map[uuid.UUID]*model.Student)
	usersMap := make(map[uuid.UUID]*model.Users)

	for _, ref := range references {
		if _, exists := studentsMap[ref.StudentID]; !exists {
			student, err := s.AchievementRepo.GetStudentByID(ref.StudentID)
			if err == nil {
				studentsMap[ref.StudentID] = student

				// Get user info
				user, err := s.AchievementRepo.GetUserByID(student.UserID)
				if err == nil {
					usersMap[student.UserID] = user
				}
			}
		}
	}

	// 4. Combine data
	var results []model.AchievementWithStudent
	for _, ref := range references {
		achievement, exists := achievementsMap[ref.MongoAchievementID]
		if !exists {
			continue
		}

		student := studentsMap[ref.StudentID]
		user := usersMap[student.UserID]

		result := model.AchievementWithStudent{
			Achievement:            achievement,
			AchievementReferenceID: ref.ID,
			Status:                 ref.Status,
			SubmittedAt:            ref.SubmittedAt,
			VerifiedAt:             ref.VerifiedAt,
			StudentName:            user.FullName,
			StudentIDNumber:        student.StudentID,
			ProgramStudy:           student.ProgramStudy,
		}

		results = append(results, result)
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))

	// Return with pagination
	return c.JSON(fiber.Map{
		"message": "success",
		"data": model.PaginatedAchievementsResponse{
			Data: results,
			Pagination: model.PaginationMeta{
				CurrentPage: page,
				PerPage:     perPage,
				TotalPages:  totalPages,
				TotalItems:  totalCount,
			},
		},
	})
}


// VerifyAchievement - Dosen verify prestasi mahasiswa (FR-007)
func (s *LecturerService) VerifyAchievement(c *fiber.Ctx) error {
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

	// Parse request body
	var req model.VerifyAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate action
	if req.Action != "approve" && req.Action != "reject" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "action must be 'approve' or 'reject'",
		})
	}

	// Validate note for reject
	if req.Action == "reject" && (req.Note == nil || *req.Note == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rejection note is required when rejecting",
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

	// Ambil lecturer_id dari user_id
	lecturer, err := s.AchievementRepo.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "lecturer not found",
		})
	}

	// Ambil achievement reference
	achievementRef, err := s.AchievementRepo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "achievement reference not found",
		})
	}

	// Verify status is submitted
	if achievementRef.Status != "submitted" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only submitted achievements can be verified",
		})
	}

	// Check apakah lecturer adalah advisor dari student
	isAdvisor, err := s.AchievementRepo.CheckLecturerOwnsStudent(lecturer.ID, achievementRef.StudentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to verify advisor relationship",
		})
	}

	if !isAdvisor {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you can only verify achievements of your own students",
		})
	}

	// Verify achievement
	err = s.AchievementRepo.VerifyAchievement(referenceID, lecturer.ID, req.Action, req.Note)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to verify achievement",
		})
	}

	// Get achievement detail untuk notifikasi
	achievement, err := s.AchievementRepo.GetAchievementByID(achievementRef.MongoAchievementID)
	if err != nil {
		achievement = &model.Achievement{
			Title: "Achievement",
		}
	}

	// Get student info
	student, err := s.AchievementRepo.GetStudentByID(achievementRef.StudentID)
	if err == nil {
		// Create notification untuk mahasiswa
		err = s.createNotificationForStudent(student, achievementRef, achievement, req.Action, req.Note)
		if err != nil {
			// Log error tapi tidak gagalkan request
		}
	}

	// Prepare response message
	var message string
	var status string
	if req.Action == "approve" {
		message = "Prestasi berhasil diverifikasi"
		status = "verified"
	} else {
		message = "Prestasi ditolak"
		status = "rejected"
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"data": model.VerifyAchievementResponse{
			AchievementReferenceID: referenceID,
			Status:                 status,
			VerifiedBy:             lecturer.ID,
			VerifiedAt:             time.Now(),
			RejectionNote:          req.Note,
			Message:                message,
		},
	})
}

// createNotificationForStudent membuat notifikasi untuk mahasiswa
func (s *LecturerService) createNotificationForStudent(student *model.Student, achievementRef *model.AchievementReference, achievement *model.Achievement, action string, note *string) error {
	// Ambil data lecturer untuk nama
	lecturer, err := s.AchievementRepo.GetLecturerByUserID(student.AdvisorID)
	if err != nil {
		return err
	}

	lecturerUser, err := s.AchievementRepo.GetUserByID(lecturer.UserID)
	if err != nil {
		return err
	}

	// Buat notification data
	notifData := model.NotificationData{
		AchievementID:          achievementRef.MongoAchievementID,
		AchievementReferenceID: achievementRef.ID,
		StudentID:              student.ID,
		StudentName:            "",
		AchievementTitle:       achievement.Title,
	}

	dataJSON, err := json.Marshal(notifData)
	if err != nil {
		return err
	}

	// Create notification
	var notifType, title, message string
	if action == "approve" {
		notifType = "achievement_verified"
		title = "Prestasi Diverifikasi"
		message = fmt.Sprintf("Prestasi '%s' Anda telah diverifikasi oleh %s", achievement.Title, lecturerUser.FullName)
	} else {
		notifType = "achievement_rejected"
		title = "Prestasi Ditolak"
		message = fmt.Sprintf("Prestasi '%s' Anda ditolak oleh %s", achievement.Title, lecturerUser.FullName)
		if note != nil && *note != "" {
			message += fmt.Sprintf(". Alasan: %s", *note)
		}
	}

	notification := model.Notification{
		ID:        uuid.New(),
		UserID:    student.UserID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      string(dataJSON),
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.AchievementRepo.CreateNotification(notification)
}


// VerifyAchievement - Dosen approve prestasi (FR-007)
func (s *LecturerService) VerifyAchievement(c *fiber.Ctx) error {
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

	// Ambil lecturer_id dari user_id
	lecturer, err := s.AchievementRepo.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "lecturer not found",
		})
	}

	// Ambil achievement reference
	achievementRef, err := s.AchievementRepo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "achievement reference not found",
		})
	}

	// Verify status is submitted
	if achievementRef.Status != "submitted" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only submitted achievements can be verified",
		})
	}

	// Verify ownership: achievement harus dari mahasiswa bimbingan
	student, err := s.AchievementRepo.GetStudentByID(achievementRef.StudentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	if student.AdvisorID != lecturer.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you can only verify achievements from your own students",
		})
	}

	// Verify achievement
	err = s.AchievementRepo.VerifyAchievement(referenceID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to verify achievement",
		})
	}

	// Ambil achievement detail untuk notifikasi
	achievement, err := s.AchievementRepo.GetAchievementByID(achievementRef.MongoAchievementID)
	if err != nil {
		achievement = &model.Achievement{
			Title: "Achievement",
		}
	}

	// Create notification untuk mahasiswa
	err = s.createNotificationForStudent(student, achievementRef, achievement, "verified", "")
	if err != nil {
		// Log error tapi tidak gagalkan request
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "prestasi berhasil diverifikasi",
		"data": fiber.Map{
			"achievement_reference_id": referenceID,
			"status":                   "verified",
			"verified_by":              userID,
		},
	})
}

// RejectAchievement - Dosen reject prestasi (FR-007)
func (s *LecturerService) RejectAchievement(c *fiber.Ctx) error {
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

	// Parse request body untuk rejection note
	var req struct {
		RejectionNote string `json:"rejection_note"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.RejectionNote == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rejection_note is required",
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

	// Ambil lecturer_id dari user_id
	lecturer, err := s.AchievementRepo.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "lecturer not found",
		})
	}

	// Ambil achievement reference
	achievementRef, err := s.AchievementRepo.GetAchievementReferenceByID(referenceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "achievement reference not found",
		})
	}

	// Verify status is submitted
	if achievementRef.Status != "submitted" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only submitted achievements can be rejected",
		})
	}

	// Verify ownership: achievement harus dari mahasiswa bimbingan
	student, err := s.AchievementRepo.GetStudentByID(achievementRef.StudentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	if student.AdvisorID != lecturer.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you can only reject achievements from your own students",
		})
	}

	// Reject achievement
	err = s.AchievementRepo.RejectAchievement(referenceID, userID, req.RejectionNote)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to reject achievement",
		})
	}

	// Ambil achievement detail untuk notifikasi
	achievement, err := s.AchievementRepo.GetAchievementByID(achievementRef.MongoAchievementID)
	if err != nil {
		achievement = &model.Achievement{
			Title: "Achievement",
		}
	}

	// Create notification untuk mahasiswa
	err = s.createNotificationForStudent(student, achievementRef, achievement, "rejected", req.RejectionNote)
	if err != nil {
		// Log error tapi tidak gagalkan request
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "prestasi ditolak",
		"data": fiber.Map{
			"achievement_reference_id": referenceID,
			"status":                   "rejected",
			"verified_by":              userID,
			"rejection_note":           req.RejectionNote,
		},
	})
}

// createNotificationForStudent membuat notifikasi untuk mahasiswa
func (s *LecturerService) createNotificationForStudent(
	student *model.Student,
	achievementRef *model.AchievementReference,
	achievement *model.Achievement,
	action string,
	rejectionNote string,
) error {
	// Ambil data lecturer user untuk nama
	lecturerUser, err := s.AchievementRepo.GetUserByID(student.UserID)
	if err != nil {
		return err
	}

	var title, message, notifType string

	if action == "verified" {
		notifType = "achievement_verified"
		title = "Prestasi Diverifikasi"
		message = fmt.Sprintf("Prestasi Anda '%s' telah diverifikasi oleh dosen wali", achievement.Title)
	} else {
		notifType = "achievement_rejected"
		title = "Prestasi Ditolak"
		message = fmt.Sprintf("Prestasi Anda '%s' ditolak. Alasan: %s", achievement.Title, rejectionNote)
	}

	// Buat notification data
	notifData := model.NotificationData{
		AchievementID:          achievementRef.MongoAchievementID,
		AchievementReferenceID: achievementRef.ID,
		StudentID:              student.ID,
		StudentName:            lecturerUser.FullName,
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
		UserID:    student.UserID, // Send to student
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      string(dataJSON),
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	return s.AchievementRepo.CreateNotification(notification)
}
