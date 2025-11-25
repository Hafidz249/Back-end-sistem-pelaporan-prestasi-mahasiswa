package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewApp(db *sql.DB) *fiber.App {
    app := fiber.New()
    
    app.Use(logger.New())

    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Welcome to the Alumni API",
            "success": true,
        })
    })

    return app
}