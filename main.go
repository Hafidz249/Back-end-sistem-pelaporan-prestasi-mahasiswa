package main

import (
	"fmt"

	config "POJECT_UAS/Config"
	route "POJECT_UAS/Routes"
	"POJECT_UAS/middleware"
	"POJECT_UAS/repository"
	"POJECT_UAS/service"

	"github.com/gofiber/fiber/v2"
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

	// Initialize service (service sekarang juga sebagai handler)
	authService := service.NewAuthService(authRepo)
	achievementService := service.NewAchievementService(achievementRepo)
	lecturerService := service.NewLecturerService(achievementRepo)

	// Initialize middleware
	permMiddleware := middleware.NewPermissionMiddleware(postgresDB)
	roleMiddleware := middleware.NewRoleMiddleware(postgresDB)

	app := fiber.New()

	// simple logger middleware
	app.Use(config.LoggerMiddleware)

	// register routes with middleware
	route.SetupRoutes(app, authService, achievementService, lecturerService, permMiddleware, roleMiddleware)

	port := config.GetAppPort()
	fmt.Printf("Starting server on :%s\n", port)
	app.Listen(":" + port)
}
