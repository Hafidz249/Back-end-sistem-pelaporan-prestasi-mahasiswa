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

	// Example: Achievements routes with permission check
	achievements := api.Group("/achievements")
	achievements.Get("/",
		permMiddleware.RequirePermission("achievements", "read"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list achievements"})
		},
	)
	achievements.Post("/",
		permMiddleware.RequirePermission("achievements", "create"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "create achievement"})
		},
	)
	achievements.Put("/:id",
		permMiddleware.RequirePermission("achievements", "update"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "update achievement"})
		},
	)
	achievements.Delete("/:id",
		permMiddleware.RequirePermission("achievements", "delete"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "delete achievement"})
		},
	)

	// Example: Admin routes with role check
	admin := api.Group("/admin", roleMiddleware.RequireRole("admin", "super_admin"))
	admin.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "list all users (admin only)"})
	})

	// Example: Student routes
	students := api.Group("/students")
	students.Get("/",
		permMiddleware.RequirePermission("students", "read"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list students"})
		},
	)

	// TODO: register other routes (users, roles, lecturers, etc.)
}