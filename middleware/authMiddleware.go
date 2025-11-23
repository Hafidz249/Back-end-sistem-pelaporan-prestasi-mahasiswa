package middleware

import (
	"strings"

	"POJECT_UAS/Config"
	"POJECT_UAS/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(userRepo *model.UserRepository) fiber.Handler {
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

        
        var userID float64
        var role string
        var username string

        if v, ok := claims["sub"].(float64); ok {
            userID = v
        } else if v, ok := claims["id"].(float64); ok {
            userID = v
        }

        if r, exists := claims["role"].(string); exists {
            role = r
        }
        
        if u, exists := claims["username"].(string); exists {
            username = u
        }


        if userID != 0 {
            c.Locals("id", int(userID))
        }
        c.Locals("role", role)
        c.Locals("username", username)

        return c.Next()
    }
}