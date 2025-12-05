package middleware

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RoleMiddleware untuk check role
type RoleMiddleware struct {
	DB *sql.DB
}

func NewRoleMiddleware(db *sql.DB) *RoleMiddleware {
	return &RoleMiddleware{DB: db}
}

// RequireRole check apakah user memiliki salah satu role yang diizinkan
func (rm *RoleMiddleware) RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "role not found",
			})
		}

		roleUUID, err := uuid.Parse(roleID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "invalid role id",
			})
		}

		// Get role name from database
		var roleName string
		query := `SELECT name FROM roles WHERE id = $1`
		err = rm.DB.QueryRow(query, roleUUID).Scan(&roleName)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "role not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to check role",
			})
		}

		// Check apakah role ada di allowed roles
		for _, allowedRole := range allowedRoles {
			if roleName == allowedRole {
				c.Locals("role_name", roleName)
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "forbidden",
			"message": "insufficient role",
		})
	}
}