# Contoh Penggunaan RBAC di Service (Handler)

## 1. Menggunakan Middleware di Route

```go
// Di Routes/Router.go
achievements := api.Group("/achievements")

// Semua user yang login bisa read
achievements.Get("/", 
    permMiddleware.RequirePermission("achievements", "read"),
    achievementService.List,
)

// Hanya user dengan permission create
achievements.Post("/",
    permMiddleware.RequirePermission("achievements", "create"),
    achievementService.Create,
)

// Hanya user dengan permission update
achievements.Put("/:id",
    permMiddleware.RequirePermission("achievements", "update"),
    achievementService.Update,
)

// Hanya user dengan permission delete
achievements.Delete("/:id",
    permMiddleware.RequirePermission("achievements", "delete"),
    achievementService.Delete,
)
```

## 2. Menggunakan Helper di Service

```go
package service

import (
    "POJECT_UAS/middleware"
    "POJECT_UAS/repository"
    "github.com/gofiber/fiber/v2"
)

type AchievementService struct {
    AchievementRepo *repository.AchievementRepository
}

func NewAchievementService(repo *repository.AchievementRepository) *AchievementService {
    return &AchievementService{
        AchievementRepo: repo,
    }
}

func (s *AchievementService) List(c *fiber.Ctx) error {
    // Ambil user info dari context
    userID := middleware.GetUserID(c)
    username := middleware.GetUsername(c)
    
    // Logic untuk list achievements
    achievements, err := s.AchievementRepo.GetAll()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to get achievements",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "list achievements",
        "data": achievements,
        "user_id": userID,
        "username": username,
    })
}

func (s *AchievementService) Update(c *fiber.Ctx) error {
    achievementID := c.Params("id")
    userID := middleware.GetUserID(c)
    
    // Check ownership: user hanya bisa update achievement miliknya sendiri
    // kecuali dia admin
    isAdmin := middleware.HasPermission(c, "achievements", "update_all")
    
    if !isAdmin {
        // Verify ownership
        // ... query database untuk check apakah achievement ini milik user
    }
    
    // Logic untuk update achievement
    // ...
    
    return c.JSON(fiber.Map{
        "message": "achievement updated",
        "id": achievementID,
    })
}

func (s *AchievementService) Delete(c *fiber.Ctx) error {
    achievementID := c.Params("id")
    
    // Check apakah user punya permission delete_all
    canDeleteAll := middleware.HasPermission(c, "achievements", "delete_all")
    
    if canDeleteAll {
        // Admin bisa delete semua achievement
        // ... delete logic
    } else {
        // User biasa hanya bisa delete miliknya sendiri
        userID := middleware.GetUserID(c)
        // ... verify ownership dan delete
    }
    
    return c.JSON(fiber.Map{
        "message": "achievement deleted",
        "id": achievementID,
    })
}
```

## 3. Conditional Logic Berdasarkan Permission

```go
func (s *StudentService) GetProfile(c *fiber.Ctx) error {
    studentID := c.Params("id")
    userID := middleware.GetUserID(c)
    
    // Check apakah user bisa view semua student
    canViewAll := middleware.HasPermission(c, "students", "read_all")
    
    var student Student
    var err error
    
    if canViewAll {
        // Admin/Lecturer bisa view semua student
        student, err = s.StudentRepo.GetByID(studentID)
    } else {
        // Student hanya bisa view profile sendiri
        student, err = s.StudentRepo.GetByUserID(userID)
        
        if student.ID != studentID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "you can only view your own profile",
            })
        }
    }
    
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "student not found",
        })
    }
    
    return c.JSON(fiber.Map{
        "data": student,
    })
}
```

## 4. Multiple Permission Check

```go
func (s *AchievementService) Verify(c *fiber.Ctx) error {
    achievementID := c.Params("id")
    
    // Verify achievement butuh 2 permissions:
    // 1. achievements:verify
    // 2. achievements:update
    
    canVerify := middleware.HasPermission(c, "achievements", "verify")
    canUpdate := middleware.HasPermission(c, "achievements", "update")
    
    if !canVerify || !canUpdate {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "you don't have permission to verify achievements",
        })
    }
    
    verifierID := middleware.GetUserID(c)
    
    // Logic untuk verify achievement
    // ... update status, set verified_by, verified_at
    
    return c.JSON(fiber.Map{
        "message": "achievement verified",
        "id": achievementID,
        "verified_by": verifierID,
    })
}
```

## 5. Role-Based Logic

```go
func (s *UserService) ListUsers(c *fiber.Ctx) error {
    roleName := middleware.GetRoleName(c)
    userID := middleware.GetUserID(c)
    
    var users []User
    var err error
    
    switch roleName {
    case "admin", "super_admin":
        // Admin bisa lihat semua user
        users, err = s.UserRepo.GetAll()
        
    case "lecturer":
        // Lecturer hanya bisa lihat student yang dia bimbing
        users, err = s.UserRepo.GetStudentsByAdvisor(userID)
        
    case "student":
        // Student tidak bisa akses endpoint ini
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "students cannot view user list",
        })
        
    default:
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "unknown role",
        })
    }
    
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to get users",
        })
    }
    
    return c.JSON(fiber.Map{
        "data": users,
    })
}
```

## 6. Dynamic Permission Check

```go
func (s *ResourceService) AccessResource(c *fiber.Ctx) error {
    resourceType := c.Params("type") // achievements, students, users, etc.
    action := c.Query("action", "read") // read, create, update, delete
    
    // Dynamic permission check
    if !middleware.HasPermission(c, resourceType, action) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "forbidden",
            "message": fmt.Sprintf("you don't have permission to %s %s", action, resourceType),
        })
    }
    
    // Logic untuk access resource
    // ...
    
    return c.JSON(fiber.Map{
        "message": fmt.Sprintf("%s %s success", action, resourceType),
    })
}
```

## 7. Audit Log dengan User Info

```go
func (s *AchievementService) Create(c *fiber.Ctx) error {
    var req CreateAchievementRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request",
        })
    }
    
    // Ambil user info untuk audit log
    userID := middleware.GetUserID(c)
    username := middleware.GetUsername(c)
    
    // Create achievement
    achievement, err := s.AchievementRepo.Create(req, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "failed to create achievement",
        })
    }
    
    // Log activity (optional)
    // s.auditLog.Log(AuditLog{
        UserID:    userID,
        Username:  username,
        Action:    "create_achievement",
        Resource:  "achievements",
        ResourceID: achievement.ID,
        Timestamp: time.Now(),
    })
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "achievement created",
        "data": achievement,
    })
}
```
