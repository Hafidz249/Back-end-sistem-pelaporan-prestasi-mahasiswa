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

	// initialize DB
	db := config.InitDB()

	// Initialize repository
	authRepo := &repository.AuthRepository{
		DB:        db,
		JWTSecret: config.GetJWTSecret(),
	}

	// Initialize service (service sekarang juga sebagai handler)
	authService := service.NewAuthService(authRepo)

	// Initialize middleware
	permMiddleware := middleware.NewPermissionMiddleware(db)
	roleMiddleware := middleware.NewRoleMiddleware(db)

	app := fiber.New()

	// simple logger middleware
	app.Use(config.LoggerMiddleware)

	// register routes with middleware
	route.SetupRoutes(app, authService, permMiddleware, roleMiddleware)

	port := config.GetAppPort()
	fmt.Printf("Starting server on :%s\n", port)
	app.Listen(":" + port)
}
