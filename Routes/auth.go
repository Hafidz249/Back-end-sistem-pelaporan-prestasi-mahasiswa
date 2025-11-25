package routes

import (
	"POJECT_UAS/model"
	"POJECT_UAS/service"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoutes(app *fiber.App, userRepo model.UserRepository) {
	authService := service.NewAuthService(userRepo)

	app.Post("/register", func(c *fiber.Ctx) error {
		var req model.Users
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
		}

		user := &model.Users{
			Email:    req.Email,
			Username: req.Username,
			Password: string(hashedPassword),
			Role:     req.Role,
		}

		if err := userRepo.Create(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		user.Password = ""
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "user registered successfully",
			"user":    user,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		body := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		token, user, err := authService.Login(body.Email, body.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user":  user,
		})
	})

	// contoh route untuk cek token (opsional)
	app.Get("/me", func(c *fiber.Ctx) error {
		uid := c.Locals("user_id")
		role := c.Locals("role")
		return c.JSON(fiber.Map{"user_id": uid, "role": role})
	})
}

func UserRoutes(app *fiber.App) {
	app.Get("/users", service.GetUsersService)
}