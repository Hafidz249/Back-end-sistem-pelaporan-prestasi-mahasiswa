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
	UserRepo *repository.UserRepository
}

func NewAdminService(userRepo *repository.UserRepository) *AdminService {
	return &AdminService{
		UserRepo: userRepo,
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
