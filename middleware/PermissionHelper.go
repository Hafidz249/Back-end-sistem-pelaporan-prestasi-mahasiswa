package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// HasPermission check apakah user memiliki permission tertentu (untuk digunakan di handler)
func HasPermission(c *fiber.Ctx, resource, action string) bool {
	permissions, ok := c.Locals("permissions").([]map[string]interface{})
	if !ok {
		return false
	}

	for _, perm := range permissions {
		permResource, _ := perm["resource"].(string)
		permAction, _ := perm["action"].(string)

		if permResource == resource && permAction == action {
			return true
		}
	}

	return false
}

// GetUserID helper untuk ambil user_id dari context
func GetUserID(c *fiber.Ctx) string {
	userID, _ := c.Locals("user_id").(string)
	return userID
}

// GetUsername helper untuk ambil username dari context
func GetUsername(c *fiber.Ctx) string {
	username, _ := c.Locals("username").(string)
	return username
}

// GetEmail helper untuk ambil email dari context
func GetEmail(c *fiber.Ctx) string {
	email, _ := c.Locals("email").(string)
	return email
}

// GetRoleID helper untuk ambil role_id dari context
func GetRoleID(c *fiber.Ctx) string {
	roleID, _ := c.Locals("role_id").(string)
	return roleID
}

// GetRoleName helper untuk ambil role_name dari context
func GetRoleName(c *fiber.Ctx) string {
	roleName, _ := c.Locals("role_name").(string)
	return roleName
}

// GetPermissions helper untuk ambil semua permissions dari context
func GetPermissions(c *fiber.Ctx) []map[string]interface{} {
	permissions, _ := c.Locals("permissions").([]map[string]interface{})
	return permissions
}
