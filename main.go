package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"

	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "POJECT_UAS/docs"
)

// JWT Secret key untuk demo
const jwtSecret = "your-super-secret-jwt-key-for-demo-presentation-2024"

// JWT Claims structure
type Claims struct {
	UserID      string   `json:"sub"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// Generate JWT Token dengan expiration
func generateJWTToken(userID, username, email, role string, permissions []string) (string, error) {
	// Set expiration time (24 hours from now)
	expirationTime := time.Now().Add(24 * time.Hour)
	
	// Create claims
	claims := &Claims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "prestasi-api",
			Subject:   userID,
		},
	}
	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token with secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Sistem Pelaporan Prestasi Mahasiswa API v1.0.0",
	})

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// =========================
	// Swagger UI
	// =========================
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
			"status":  "API is running",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := app.Group("/api/v1")

	// API info
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Sistem Pelaporan Prestasi Mahasiswa API",
			"version": "1.0.0",
			"features": []string{
				"JWT Authentication & Authorization",
				"Role-Based Access Control (RBAC)",
				"Achievement Management (CRUD)",
				"Verification Workflow",
				"Statistics & Reporting",
				"Dual Database (PostgreSQL + MongoDB)",
			},
		})
	})

	// =========================
	// Authentication
	// =========================
	auth := v1.Group("/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		// Parse request body
		var loginReq struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		if err := c.BodyParser(&loginReq); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
				"message": "Please provide valid JSON with username and password",
			})
		}
		
		// Demo user data (in real app, this would be from database)
		var userID, username, email, role string
		var permissions []string
		
		// Simple demo authentication (in real app, verify against database)
		switch loginReq.Username {
		case "student123", "student@example.com":
			userID = "123e4567-e89b-12d3-a456-426614174000"
			username = "student123"
			email = "student@example.com"
			role = "student"
			permissions = []string{"achievements:create", "achievements:read", "achievements:update", "achievements:delete"}
		case "lecturer123", "lecturer@example.com":
			userID = "223e4567-e89b-12d3-a456-426614174001"
			username = "lecturer123"
			email = "lecturer@example.com"
			role = "lecturer"
			permissions = []string{"achievements:verify", "achievements:reject", "students:read", "reports:advisee"}
		case "admin123", "admin@example.com":
			userID = "323e4567-e89b-12d3-a456-426614174002"
			username = "admin123"
			email = "admin@example.com"
			role = "admin"
			permissions = []string{"users:*", "achievements:*", "reports:*", "system:*"}
		default:
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid credentials",
				"message": "Username or password is incorrect",
			})
		}
		
		// Generate JWT token with expiration
		token, err := generateJWTToken(userID, username, email, role, permissions)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "token generation failed",
				"message": "Failed to generate authentication token",
			})
		}
		
		// Calculate expiration time for response
		expiresAt := time.Now().Add(24 * time.Hour)
		
		return c.JSON(fiber.Map{
			"message": "Login successful",
			"data": fiber.Map{
				"token": token,
				"expires_at": expiresAt.Format(time.RFC3339),
				"expires_in": 86400, // 24 hours in seconds
				"token_type": "Bearer",
				"user": fiber.Map{
					"id":          userID,
					"username":    username,
					"email":       email,
					"full_name":   username + " User",
					"role":        role,
					"permissions": permissions,
				},
			},
		})
	})

	// Token validation endpoint
	auth.Post("/validate", func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing authorization header",
				"message": "Please provide Authorization header with Bearer token",
			})
		}
		
		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid authorization format",
				"message": "Authorization header must be in format: Bearer <token>",
			})
		}
		
		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
				"message": "Token is invalid or expired",
				"details": err.Error(),
			})
		}
		
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return c.JSON(fiber.Map{
				"message": "Token is valid",
				"data": fiber.Map{
					"user_id":     claims.UserID,
					"username":    claims.Username,
					"email":       claims.Email,
					"role":        claims.Role,
					"permissions": claims.Permissions,
					"expires_at":  claims.ExpiresAt.Format(time.RFC3339),
					"issued_at":   claims.IssuedAt.Format(time.RFC3339),
				},
			})
		}
		
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token claims",
			"message": "Token claims are invalid",
		})
	})

	// Refresh token endpoint
	auth.Post("/refresh", func(c *fiber.Ctx) error {
		// Get current token
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}
		
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}
		
		// Parse token (even if expired, we can still get claims)
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		
		if err != nil {
			// For refresh, we try to parse even expired tokens
			// In JWT v5, we handle this differently
			parser := jwt.NewParser(jwt.WithoutClaimsValidation())
			token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			
			if err != nil {
				return c.Status(401).JSON(fiber.Map{
					"error": "invalid token",
					"message": "Cannot refresh invalid token",
				})
			}
			
			// If we can get claims (even from expired token), generate new one
			if claims, ok := token.Claims.(*Claims); ok {
				// Generate new token
				newToken, err := generateJWTToken(
					claims.UserID,
					claims.Username,
					claims.Email,
					claims.Role,
					claims.Permissions,
				)
				if err != nil {
					return c.Status(500).JSON(fiber.Map{
						"error": "failed to generate new token",
					})
				}
				
				expiresAt := time.Now().Add(24 * time.Hour)
				return c.JSON(fiber.Map{
					"message": "Token refreshed successfully",
					"data": fiber.Map{
						"token":      newToken,
						"expires_at": expiresAt.Format(time.RFC3339),
						"expires_in": 86400,
						"token_type": "Bearer",
					},
				})
			}
			
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
				"message": "Cannot refresh invalid token",
			})
		}
		
		// If token is still valid, return new one anyway
		if claims, ok := token.Claims.(*Claims); ok {
			newToken, err := generateJWTToken(
				claims.UserID,
				claims.Username,
				claims.Email,
				claims.Role,
				claims.Permissions,
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "failed to generate new token",
				})
			}
			
			expiresAt := time.Now().Add(24 * time.Hour)
			return c.JSON(fiber.Map{
				"message": "Token refreshed successfully",
				"data": fiber.Map{
					"token":      newToken,
					"expires_at": expiresAt.Format(time.RFC3339),
					"expires_in": 86400,
					"token_type": "Bearer",
				},
			})
		}
		
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token claims",
		})
	})

	auth.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"id":        "123",
			"username":  "student123",
			"email":     "student@example.com",
			"full_name": "John Doe",
			"role":      "student",
		})
	})

	// =========================
	// Achievements
	// =========================
	achievements := v1.Group("/achievements")

	achievements.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Achievement list (demo)",
		})
	})

	achievements.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Achievement created (demo)",
		})
	})

	// =========================
	// Users (Admin)
	// =========================
	users := v1.Group("/users")

	users.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Users list (admin demo)",
		})
	})

	// =========================
	// Reports
	// =========================
	reports := v1.Group("/reports")

	reports.Get("/statistics", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Statistics demo",
		})
	})

	// =========================
	// Start server
	// =========================
	port := "8080"

	fmt.Println("🚀 Sistem Pelaporan Prestasi Mahasiswa API")
	fmt.Println("📍 Server :", "http://localhost:"+port)
	fmt.Println("📚 Swagger:", "http://localhost:"+port+"/swagger/index.html")
	fmt.Println("💚 Health :", "http://localhost:"+port+"/ping")

	log.Fatal(app.Listen(":" + port))
}
