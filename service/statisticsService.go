package service

import (
	"POJECT_UAS/model"
	"POJECT_UAS/repository"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StatisticsService struct {
	AchievementRepo *repository.AchievementRepository
}

func NewStatisticsService(achievementRepo *repository.AchievementRepository) *StatisticsService {
	return &StatisticsService{
		AchievementRepo: achievementRepo,
	}
}

// GetMyStatistics - Mahasiswa melihat statistik prestasi sendiri (FR-011)
func (s *StatisticsService) GetMyStatistics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	// Get student data
	student, err := s.AchievementRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "student not found",
		})
	}

	// Parse query parameters
	startDate := parseDate(c.Query("start_date"))
	endDate := parseDate(c.Query("end_date"))
	achievementType := parseString(c.Query("achievement_type"))
	status := parseString(c.Query("status"))

	// Get statistics for this student only
	studentIDs := []uuid.UUID{student.ID}
	statistics, err := s.AchievementRepo.GetAchievementStatistics(
		studentIDs,
		startDate,
		endDate,
		achievementType,
		status,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get statistics",
		})
	}

	// Sort and limit results
	s.sortAndLimitStatistics(statistics, 10)

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    statistics,
	})
}

// GetAdviseeStatistics - Dosen wali melihat statistik mahasiswa bimbingan (FR-011)
func (s *StatisticsService) GetAdviseeStatistics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	// Get lecturer data
	lecturer, err := s.AchievementRepo.GetLecturerByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "lecturer not found",
		})
	}

	// Get student IDs yang dibimbing
	studentIDs, err := s.AchievementRepo.GetStudentIDsByAdvisor(lecturer.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get advisee students",
		})
	}

	if len(studentIDs) == 0 {
		return c.JSON(fiber.Map{
			"message": "success",
			"data": &model.AchievementStatistics{
				TotalByType:       []model.TypeStatistic{},
				TotalByPeriod:     []model.PeriodStatistic{},
				TopStudents:       []model.TopStudent{},
				CompetitionLevels: []model.LevelStatistic{},
				Summary: model.StatisticSummary{
					TotalAchievements:    0,
					VerifiedAchievements: 0,
					PendingAchievements:  0,
					RejectedAchievements: 0,
					TotalStudents:        0,
				},
			},
		})
	}

	// Parse query parameters
	startDate := parseDate(c.Query("start_date"))
	endDate := parseDate(c.Query("end_date"))
	achievementType := parseString(c.Query("achievement_type"))
	status := parseString(c.Query("status"))
	topLimit, _ := strconv.Atoi(c.Query("top_limit", "10"))

	// Get statistics for advisee students
	statistics, err := s.AchievementRepo.GetAchievementStatistics(
		studentIDs,
		startDate,
		endDate,
		achievementType,
		status,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get statistics",
		})
	}

	// Sort and limit results
	s.sortAndLimitStatistics(statistics, topLimit)

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    statistics,
	})
}

// GetAllStatistics - Admin melihat statistik semua prestasi (FR-011)
func (s *StatisticsService) GetAllStatistics(c *fiber.Ctx) error {
	// Parse query parameters
	startDate := parseDate(c.Query("start_date"))
	endDate := parseDate(c.Query("end_date"))
	achievementType := parseString(c.Query("achievement_type"))
	status := parseString(c.Query("status"))
	topLimit, _ := strconv.Atoi(c.Query("top_limit", "20"))

	// Get statistics for all students (empty studentIDs)
	statistics, err := s.AchievementRepo.GetAchievementStatistics(
		nil, // all students
		startDate,
		endDate,
		achievementType,
		status,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get statistics",
		})
	}

	// Sort and limit results
	s.sortAndLimitStatistics(statistics, topLimit)

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    statistics,
	})
}

// Helper functions
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	if date, err := time.Parse("2006-01-02", dateStr); err == nil {
		return &date
	}
	return nil
}

func parseString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func (s *StatisticsService) sortAndLimitStatistics(statistics *model.AchievementStatistics, topLimit int) {
	// Sort top students by total count (descending)
	sort.Slice(statistics.TopStudents, func(i, j int) bool {
		if statistics.TopStudents[i].TotalCount == statistics.TopStudents[j].TotalCount {
			return statistics.TopStudents[i].VerifiedCount > statistics.TopStudents[j].VerifiedCount
		}
		return statistics.TopStudents[i].TotalCount > statistics.TopStudents[j].TotalCount
	})

	// Limit top students
	if len(statistics.TopStudents) > topLimit {
		statistics.TopStudents = statistics.TopStudents[:topLimit]
	}

	// Sort type statistics by count (descending)
	sort.Slice(statistics.TotalByType, func(i, j int) bool {
		return statistics.TotalByType[i].Count > statistics.TotalByType[j].Count
	})

	// Sort period statistics by period (ascending)
	sort.Slice(statistics.TotalByPeriod, func(i, j int) bool {
		return statistics.TotalByPeriod[i].Period < statistics.TotalByPeriod[j].Period
	})

	// Sort level statistics by count (descending)
	sort.Slice(statistics.CompetitionLevels, func(i, j int) bool {
		return statistics.CompetitionLevels[i].Count > statistics.CompetitionLevels[j].Count
	})
}