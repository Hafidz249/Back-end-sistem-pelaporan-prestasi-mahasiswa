package middleware

import (
	"strings"

	"POJECT_UAS/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header format must be Bearer {token}"})
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// Ambil data dari claims
		var userID string
		var roleID string
		var username string

		if v, exists := claims["sub"].(string); exists {
			userID = v
		} else if v, exists := claims["id"].(string); exists {
			userID = v
		}

		if r, exists := claims["role_id"].(string); exists {
			roleID = r
		}

		if u, exists := claims["username"].(string); exists {
			username = u
		}

		c.Locals("id", userID)
		c.Locals("role_id", roleID)
		c.Locals("username", username)

		return c.Next()
	}
}