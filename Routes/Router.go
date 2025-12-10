package route

import (
	"POJECT_UAS/middleware"
	"POJECT_UAS/service"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes register routes to the provided Fiber app
func SetupRoutes(
	app *fiber.App,
	authService *service.AuthService,
	achievementService *service.AchievementService,
	lecturerService *service.LecturerService,
	adminService *service.AdminService,
	permMiddleware *middleware.PermissionMiddleware,
	roleMiddleware *middleware.RoleMiddleware,
) {
	// Health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "pong"})
	})

	// Auth routes (public)
	auth := app.Group("/api/auth")
	auth.Post("/login", authService.Login)

	// Protected routes - require authentication
	api := app.Group("/api", middleware.JWTAuth())

	// Example: User profile (authenticated users only)
	api.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id":  c.Locals("user_id"),
			"username": c.Locals("username"),
			"email":    c.Locals("email"),
		})
	})

	// Achievements routes
	achievements := api.Group("/achievements")

	// FR-003: Submit Prestasi (create draft)
	achievements.Post("/",
		permMiddleware.RequirePermission("achievements", "create"),
		achievementService.SubmitAchievement,
	)

	// FR-004: Submit untuk Verifikasi (draft -> submitted)
	achievements.Post("/:reference_id/submit",
		achievementService.SubmitForVerification,
	)

	// Mahasiswa melihat prestasi sendiri
	achievements.Get("/my",
		achievementService.GetMyAchievements,
	)

	// Melihat detail prestasi
	achievements.Get("/:id",
		achievementService.GetAchievementDetail,
	)

	// FR-005: Hapus Prestasi (soft delete)
	achievements.Delete("/:reference_id",
		achievementService.DeleteAchievement,
	)

	// TODO: Update achievement
	achievements.Put("/:id",
		permMiddleware.RequirePermission("achievements", "update"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "update achievement - coming soon"})
		},
	)

	// FR-009: Admin routes - Manage Users
	admin := api.Group("/admin", roleMiddleware.RequireRole("admin", "super_admin"))

	// User management
	admin.Get("/users", adminService.GetAllUsers)
	admin.Post("/users", adminService.CreateUser)
	admin.Put("/users/:user_id", adminService.UpdateUser)
	admin.Delete("/users/:user_id", adminService.DeleteUser)

	// Profile management
	admin.Post("/students/profile", adminService.CreateStudentProfile)
	admin.Post("/lecturers/profile", adminService.CreateLecturerProfile)

	// Advisor management
	admin.Put("/students/:student_id/advisor", adminService.UpdateStudentAdvisor)

	// Role management
	admin.Get("/roles", adminService.GetAllRoles)

	// FR-010: View All Achievements
	admin.Get("/achievements", adminService.ViewAllAchievements)

	// Example: Student routes
	students := api.Group("/students")
	students.Get("/",
		permMiddleware.RequirePermission("students", "read"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list students"})
		},
	)

	// Lecturer routes
	lecturer := api.Group("/lecturer", roleMiddleware.RequireRole("lecturer", "dosen"))

	// FR-006: View Prestasi Mahasiswa Bimbingan
	lecturer.Get("/students/achievements",
		lecturerService.ViewStudentAchievements,
	)

	// FR-007 & FR-008: Verify/Reject Prestasi
	lecturer.Post("/achievements/:reference_id/verify",
		lecturerService.VerifyAchievement,
	)
	lecturer.Post("/achievements/:reference_id/reject",
		lecturerService.RejectAchievement,
	)
}