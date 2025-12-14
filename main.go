// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0.0
// @description API untuk sistem pelaporan prestasi mahasiswa dengan fitur Authentication, RBAC, Achievement Management, Statistics & Reporting
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"

	config "POJECT_UAS/Config"
	_ "POJECT_UAS/docs"
	route "POJECT_UAS/Routes"
	"POJECT_UAS/middleware"
	"POJECT_UAS/repository"
	"POJECT_UAS/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// load .env
	config.LoadEnv()

	// initialize PostgreSQL
	postgresDB := config.InitDB()

	// initialize MongoDB
	mongoDB := config.InitMongoDB()

	// Initialize repository
	authRepo := &repository.AuthRepository{
		DB:        postgresDB,
		JWTSecret: config.GetJWTSecret(),
	}

	achievementRepo := repository.NewAchievementRepository(postgresDB, mongoDB)
	userRepo := repository.NewUserRepository(postgresDB)

	// Initialize service (service sekarang juga sebagai handler)
	authService := service.NewAuthService(authRepo)
	achievementService := service.NewAchievementService(achievementRepo)
	lecturerService := service.NewLecturerService(achievementRepo)
	adminService := service.NewAdminService(userRepo, achievementRepo)
	statisticsService := service.NewStatisticsService(achievementRepo)

	// Initialize middleware
	permMiddleware := middleware.NewPermissionMiddleware(postgresDB)
	roleMiddleware := middleware.NewRoleMiddleware(postgresDB)

	app := fiber.New(fiber.Config{
		AppName: "Sistem Pelaporan Prestasi Mahasiswa API v1.0.0",
	})

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// simple logger middleware
	app.Use(config.LoggerMiddleware)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// register routes with middleware
	route.SetupRoutes(app, authService, achievementService, lecturerService, adminService, statisticsService, permMiddleware, roleMiddleware)

	port := config.GetAppPort()
	fmt.Printf("Starting server on :%s\n", port)
	app.Listen(":" + port)
}
