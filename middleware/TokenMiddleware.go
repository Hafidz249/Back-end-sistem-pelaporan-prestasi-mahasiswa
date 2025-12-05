package middleware

import (
	"strings"

	config "POJECT_UAS/Config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth memvalidasi header Authorization Bearer dan menyimpan klaim ke Locals
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ekstrak JWT dari header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header format must be Bearer {token}",
			})
		}

		tokenString := parts[1]

		// 2. Validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
			}
			return []byte(config.GetJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		// 3. Load user data dan permissions dari token
		userID, _ := claims["user_id"].(string)
		username, _ := claims["username"].(string)
		email, _ := claims["email"].(string)
		roleID, _ := claims["role_id"].(string)

		// Extract permissions dari token
		var permissions []map[string]interface{}
		if perms, ok := claims["permissions"].([]interface{}); ok {
			for _, p := range perms {
				if perm, ok := p.(map[string]interface{}); ok {
					permissions = append(permissions, perm)
				}
			}
		}

		// Simpan ke context untuk digunakan di handler/middleware berikutnya
		c.Locals("user_id", userID)
		c.Locals("username", username)
		c.Locals("email", email)
		c.Locals("role_id", roleID)
		c.Locals("permissions", permissions)

		return c.Next()
	}
}