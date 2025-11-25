package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        r := c.Locals("role")
        roleStr, _ := r.(string)

        for _, ar := range allowedRoles {
            if roleStr == ar {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden: insufficient role"})
    }
}