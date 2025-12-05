package middleware

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PermissionMiddleware untuk check permission berdasarkan resource dan action
type PermissionMiddleware struct {
	DB *sql.DB
}

func NewPermissionMiddleware(db *sql.DB) *PermissionMiddleware {
	return &PermissionMiddleware{DB: db}
}

// RequirePermission check apakah user memiliki permission yang diperlukan
// Usage: RequirePermission("achievements", "read")
func (pm *PermissionMiddleware) RequirePermission(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 4. Check apakah user memiliki permission yang diperlukan
		permissions, ok := c.Locals("permissions").([]map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "no permissions found",
			})
		}

		// Check permission dari token
		hasPermission := false
		for _, perm := range permissions {
			permResource, _ := perm["resource"].(string)
			permAction, _ := perm["action"].(string)

			if permResource == resource && permAction == action {
				hasPermission = true
				break
			}
		}

		// 5. Allow/deny request
		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "forbidden",
				"message": fmt.Sprintf("you don't have permission to %s %s", action, resource),
			})
		}

		return c.Next()
	}
}

// RequirePermissionWithCache check permission dengan fallback ke database jika tidak ada di token
func (pm *PermissionMiddleware) RequirePermissionWithCache(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check permission dari token (cache)
		permissions, ok := c.Locals("permissions").([]map[string]interface{})
		if ok {
			for _, perm := range permissions {
				permResource, _ := perm["resource"].(string)
				permAction, _ := perm["action"].(string)

				if permResource == resource && permAction == action {
					return c.Next()
				}
			}
		}

		// Fallback: Load dari database jika tidak ada di token
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

		hasPermission, err := pm.checkPermissionFromDB(roleUUID, resource, action)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to check permission",
			})
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "forbidden",
				"message": fmt.Sprintf("you don't have permission to %s %s", action, resource),
			})
		}

		return c.Next()
	}
}

// checkPermissionFromDB query database untuk check permission
func (pm *PermissionMiddleware) checkPermissionFromDB(roleID uuid.UUID, resource, action string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1 AND p.resource = $2 AND p.action = $3
	`

	var count int
	err := pm.DB.QueryRow(query, roleID, resource, action).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RequireAnyPermission check apakah user memiliki salah satu dari permissions yang diperlukan
func (pm *PermissionMiddleware) RequireAnyPermission(permissions []Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPermissions, ok := c.Locals("permissions").([]map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "no permissions found",
			})
		}

		// Check apakah user memiliki salah satu permission
		for _, requiredPerm := range permissions {
			for _, userPerm := range userPermissions {
				permResource, _ := userPerm["resource"].(string)
				permAction, _ := userPerm["action"].(string)

				if permResource == requiredPerm.Resource && permAction == requiredPerm.Action {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "forbidden",
			"message": "you don't have any of the required permissions",
		})
	}
}

// RequireAllPermissions check apakah user memiliki semua permissions yang diperlukan
func (pm *PermissionMiddleware) RequireAllPermissions(permissions []Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userPermissions, ok := c.Locals("permissions").([]map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "no permissions found",
			})
		}

		// Check apakah user memiliki semua permission
		for _, requiredPerm := range permissions {
			hasPermission := false
			for _, userPerm := range userPermissions {
				permResource, _ := userPerm["resource"].(string)
				permAction, _ := userPerm["action"].(string)

				if permResource == requiredPerm.Resource && permAction == requiredPerm.Action {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":   "forbidden",
					"message": fmt.Sprintf("you don't have permission to %s %s", requiredPerm.Action, requiredPerm.Resource),
				})
			}
		}

		return c.Next()
	}
}

// Permission struct untuk multiple permission check
type Permission struct {
	Resource string
	Action   string
}
