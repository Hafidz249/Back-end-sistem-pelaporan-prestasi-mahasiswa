package service

import (
	"POJECT_UAS/model"
	"POJECT_UAS/repository"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminService struct {
	UserRepo        *repository.UserRepository
	AchievementRepo *repository.AchievementRepository
}

func NewAdminService(userRepo *repository.UserRepository, achievementRepo *repository.AchievementRepository) *AdminService {
	return &AdminService{
		UserRepo:        userRepo,
		AchievementRepo: achievementRepo,
	}
}

// CreateUser - Admin create user (FR-009)
func (s *AdminService) CreateUser(c *fiber.Ctx) error {
	var req model.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi
	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username, email, password, and full_name are required",
		})
	}

	// Create user
	user, err := s.UserRepo.CreateUser(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"data":    user,
	})
}

// UpdateUser - Admin update user (FR-009)
func (s *AdminService) UpdateUser(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Update user
	user, err := s.UserRepo.UpdateUser(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "user updated successfully",
		"data":    user,
	})
}

// DeleteUser - Admin delete user (FR-009)
func (s *AdminService) DeleteUser(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	// Delete user (soft delete)
	err = s.UserRepo.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}

// GetAllUsers - Admin get all users (FR-009)
func (s *AdminService) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, totalCount, err := s.UserRepo.GetAllUsers(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get users",
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))

	return c.JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"users": users,
			"pagination": fiber.Map{
				"current_page": page,
				"per_page":     perPage,
				"total_pages":  totalPages,
				"total_items":  totalCount,
			},
		},
	})
}

// CreateStudentProfile - Admin create student profile (FR-009)
func (s *AdminService) CreateStudentProfile(c *fiber.Ctx) error {
	var req model.CreateStudentProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi
	if req.StudentIDNumber == "" || req.ProgramStudy == "" || req.AcademicYear == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "student_id_number, program_study, and academic_year are required",
		})
	}

	// Create student profile
	student, err := s.UserRepo.CreateStudentProfile(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create student profile",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "student profile created successfully",
		"data":    student,
	})
}

// CreateLecturerProfile - Admin create lecturer profile (FR-009)
func (s *AdminService) CreateLecturerProfile(c *fiber.Ctx) error {
	var req model.CreateLecturerProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi
	if req.LecturerIDNumber == "" || req.Department == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "lecturer_id_number and department are required",
		})
	}

	// Create lecturer profile
	lecturer, err := s.UserRepo.CreateLecturerProfile(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create lecturer profile",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "lecturer profile created successfully",
		"data":    lecturer,
	})
}

// UpdateStudentAdvisor - Admin set advisor untuk student (FR-009)
func (s *AdminService) UpdateStudentAdvisor(c *fiber.Ctx) error {
	studentIDStr := c.Params("student_id")

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid student id",
		})
	}

	var req model.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Update advisor
	err = s.UserRepo.UpdateStudentAdvisor(studentID, req.AdvisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update advisor",
		})
	}

	return c.JSON(fiber.Map{
		"message": "advisor updated successfully",
	})
}

// GetAllRoles - Admin get all roles (FR-009)
func (s *AdminService) GetAllRoles(c *fiber.Ctx) error {
	roles, err := s.UserRepo.GetAllRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get roles",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    roles,
	})
}

// ViewAllAchievements - Admin view all achievements dengan filter dan pagination (FR-010)
func (s *AdminService) ViewAllAchievements(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	status := c.Query("status", "")
	achievementType := c.Query("achievement_type", "")

	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// 1. Get achievement references dengan filter dan pagination
	references, totalCount, err := s.AchievementRepo.GetAllAchievementReferencesWithPagination(
		status,
		achievementType,
		page,
		perPage,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get achievement references",
		})
	}

	if len(references) == 0 {
		return c.JSON(fiber.Map{
			"message": "success",
			"data": fiber.Map{
				"achievements": []interface{}{},
				"pagination": fiber.Map{
					"current_page": page,
					"per_page":     perPage,
					"total_pages":  0,
					"total_items":  totalCount,
				},
			},
		})
	}

	// 2. Collect MongoDB achievement IDs
	var mongoIDs []string
	var studentIDs []uuid.UUID
	for _, ref := range references {
		mongoIDs = append(mongoIDs, ref.MongoAchievementID)
		studentIDs = append(studentIDs, ref.StudentID)
	}

	// 3. Batch fetch achievements dari MongoDB
	achievementsMap, err := s.AchievementRepo.GetAchievementsByIDs(mongoIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get achievements from MongoDB",
		})
	}

	// 4. Batch fetch student info
	studentsMap := make(map[uuid.UUID]*model.Student)
	usersMap := make(map[uuid.UUID]*model.Users)

	for _, studentID := range studentIDs {
		if _, exists := studentsMap[studentID]; !exists {
			student, err := s.AchievementRepo.GetStudentByID(studentID)
			if err != nil {
				continue // Skip jika error
			}
			studentsMap[studentID] = student

			// Get user info
			user, err := s.AchievementRepo.GetUserByID(student.UserID)
			if err != nil {
				continue
			}
			usersMap[student.UserID] = user
		}
	}

	// 5. Combine data
	var result []fiber.Map
	for _, ref := range references {
		achievement, exists := achievementsMap[ref.MongoAchievementID]
		if !exists {
			continue
		}

		student := studentsMap[ref.StudentID]
		var user *model.Users
		if student != nil {
			user = usersMap[student.UserID]
		}

		// Filter by achievement type if specified
		if achievementType != "" && achievement.AchievementType != achievementType {
			continue
		}

		item := fiber.Map{
			"reference_id":       ref.ID,
			"achievement_id":     ref.MongoAchievementID,
			"student_id":         ref.StudentID,
			"achievement_type":   achievement.AchievementType,
			"title":              achievement.Title,
			"description":        achievement.Description,
			"status":             ref.Status,
			"submitted_at":       ref.SubmittedAt,
			"verified_at":        ref.VerifiedAt,
			"verified_by":        ref.VerifiedBy,
			"rejection_note":     ref.RejectionNote,
			"created_at":         ref.CreatedAt,
			"updated_at":         ref.UpdatedAt,
		}

		// Add student info if available
		if student != nil && user != nil {
			item["student_info"] = fiber.Map{
				"student_id_number": student.StudentID,
				"full_name":         user.FullName,
				"program_study":     student.ProgramStudy,
				"academic_year":     student.AcademicYear,
			}
		}

		result = append(result, item)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))

	return c.JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"achievements": result,
			"pagination": fiber.Map{
				"current_page": page,
				"per_page":     perPage,
				"total_pages":  totalPages,
				"total_items":  totalCount,
			},
		},
	})
}
