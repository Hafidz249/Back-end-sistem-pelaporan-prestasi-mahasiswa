package service

import (
	"POJECT_UAS/middleware"
	"POJECT_UAS/model"
	"POJECT_UAS/repository"
	"math"
	"strconv"

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
