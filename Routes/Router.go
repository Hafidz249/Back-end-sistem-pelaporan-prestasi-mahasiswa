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
	statisticsService *service.StatisticsService,
	permMiddleware *middleware.PermissionMiddleware,
	roleMiddleware *middleware.RoleMiddleware,
) {
	// Health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "pong"})
	})

	// API v1 routes
	v1 := app.Group("/api/v1")

	// 5.1 Authentication routes (public)
	auth := v1.Group("/auth")
	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.RefreshToken)
	auth.Post("/logout", authService.Logout)

	// Protected routes - require authentication
	authProtected := v1.Group("/auth", middleware.JWTAuth())
	authProtected.Get("/profile", authService.GetProfile)

	// Protected API routes
	api := v1.Group("", middleware.JWTAuth())

	// 5.2 Users (Admin only)
	users := api.Group("/users", roleMiddleware.RequireRole("admin", "super_admin"))
	users.Get("/", adminService.GetAllUsers)
	users.Get("/:id", adminService.GetUserByID)
	users.Post("/", adminService.CreateUser)
	users.Put("/:id", adminService.UpdateUser)
	users.Delete("/:id", adminService.DeleteUser)
	users.Put("/:id/role", adminService.UpdateUserRole)

	// 5.4 Achievements routes
	achievements := api.Group("/achievements")

	// List achievements (filtered by role)
	achievements.Get("/", achievementService.GetAchievements)

	// Detail achievement
	achievements.Get("/:id", achievementService.GetAchievementDetail)

	// Create achievement (Mahasiswa)
	achievements.Post("/",
		permMiddleware.RequirePermission("achievements", "create"),
		achievementService.SubmitAchievement,
	)

	// Update achievement (Mahasiswa)
	achievements.Put("/:id",
		permMiddleware.RequirePermission("achievements", "update"),
		achievementService.UpdateAchievement,
	)

	// Delete achievement (Mahasiswa)
	achievements.Delete("/:id",
		achievementService.DeleteAchievement,
	)

	// Submit for verification
	achievements.Post("/:id/submit",
		achievementService.SubmitForVerification,
	)

	// Verify achievement (Dosen Wali)
	achievements.Post("/:id/verify",
		roleMiddleware.RequireRole("lecturer", "dosen"),
		lecturerService.VerifyAchievement,
	)

	// Reject achievement (Dosen Wali)
	achievements.Post("/:id/reject",
		roleMiddleware.RequireRole("lecturer", "dosen"),
		lecturerService.RejectAchievement,
	)

	// Status history
	achievements.Get("/:id/history",
		achievementService.GetAchievementHistory,
	)

	// Upload attachments
	achievements.Post("/:id/attachments",
		achievementService.UploadAttachments,
	)

	// 5.5 Students & Lecturers
	students := api.Group("/students")
	students.Get("/", adminService.GetAllStudents)
	students.Get("/:id", adminService.GetStudentByID)
	students.Get("/:id/achievements", achievementService.GetStudentAchievements)
	students.Put("/:id/advisor", 
		roleMiddleware.RequireRole("admin", "super_admin"),
		adminService.UpdateStudentAdvisor,
	)

	lecturers := api.Group("/lecturers")
	lecturers.Get("/", adminService.GetAllLecturers)
	lecturers.Get("/:id/advisees", lecturerService.GetAdvisees)

	// 5.8 Reports & Analytics
	reports := api.Group("/reports")
	reports.Get("/statistics", statisticsService.GetAllStatistics)
	reports.Get("/student/:id", statisticsService.GetStudentReport)

	// Admin routes (legacy support)
	admin := api.Group("/admin", roleMiddleware.RequireRole("admin", "super_admin"))
	admin.Post("/students/profile", adminService.CreateStudentProfile)
	admin.Post("/lecturers/profile", adminService.CreateLecturerProfile)
	admin.Get("/roles", adminService.GetAllRoles)
}