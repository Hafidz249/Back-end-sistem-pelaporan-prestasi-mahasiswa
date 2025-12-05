# Arsitektur Project

## Struktur Folder

```
POJECT_UAS/
├── Config/              # Konfigurasi aplikasi
│   ├── config.go       # Environment variables
│   ├── env.go          # Load .env
│   ├── logger.go       # Logger middleware
│   └── token.go        # Database connection
│
├── model/              # Data models
│   ├── Users.go        # User model + Login request/response
│   ├── Roles.go        # Role model
│   ├── Permission.go   # Permission model
│   ├── Role_Permission.go  # Role-Permission junction
│   ├── Student.go      # Student model
│   ├── Lecturers.go    # Lecturer model
│   ├── achievement_references.go  # Achievement reference (PostgreSQL)
│   └── achievements.go # Achievement model (MongoDB)
│
├── repository/         # Database layer
│   ├── authRepo.go     # Auth repository (login, JWT)
│   └── user_repository.go
│
├── service/            # Business logic + HTTP handlers
│   └── authService.go  # Auth service (login handler)
│
├── middleware/         # Middleware
│   ├── TokenMiddleware.go      # JWT authentication
│   ├── PermissionMiddleware.go # Permission check
│   ├── RoleRequarment.go       # Role check
│   └── PermissionHelper.go     # Helper functions
│
├── Routes/             # Route definitions
│   └── Router.go       # Setup routes
│
├── main.go             # Entry point
├── .env                # Environment variables
└── go.mod              # Go modules
```

## Layer Architecture

### 1. Model Layer
- Definisi struct untuk database tables
- Request/Response models
- Tidak ada business logic

### 2. Repository Layer
- Akses database (PostgreSQL, MongoDB)
- Query dan CRUD operations
- Return data atau error

### 3. Service Layer (Handler + Business Logic)
- **Menggabungkan handler dan business logic**
- Menerima `*fiber.Ctx` sebagai parameter
- Validasi input
- Memanggil repository
- Return HTTP response

### 4. Middleware Layer
- Authentication (JWT)
- Authorization (Permission, Role)
- Helper functions

### 5. Route Layer
- Definisi endpoint
- Mapping endpoint ke service
- Apply middleware

## Flow Request

```
Client Request
    ↓
Router (Routes/Router.go)
    ↓
Middleware (JWT, Permission, Role)
    ↓
Service (service/*.go)
    ↓
Repository (repository/*.go)
    ↓
Database (PostgreSQL, MongoDB)
    ↓
Response to Client
```

## Contoh Flow Login

```
POST /api/auth/login
    ↓
Router.go → authService.Login
    ↓
authService.Login(c *fiber.Ctx)
    - Parse request body
    - Validasi input
    - Call authRepo.Login()
    ↓
authRepo.Login(req LoginRequest)
    - Query database
    - Validate password
    - Get user profile + permissions
    - Generate JWT token
    ↓
Return LoginResponse
    - token
    - profile (user + role + permissions)
```

## Contoh Flow Protected Endpoint

```
GET /api/achievements
    ↓
Router.go → middleware.JWTAuth()
    - Extract JWT from header
    - Validate token
    - Load user data to context
    ↓
Router.go → permMiddleware.RequirePermission("achievements", "read")
    - Check permission from context
    - Allow/Deny request
    ↓
achievementService.List(c *fiber.Ctx)
    - Get user info from context
    - Call achievementRepo.GetAll()
    - Return response
```

## Dependency Injection

```go
// main.go
func main() {
    // 1. Load config
    config.LoadEnv()
    db := config.InitDB()
    
    // 2. Initialize repository
    authRepo := &repository.AuthRepository{
        DB: db,
        JWTSecret: config.GetJWTSecret(),
    }
    
    // 3. Initialize service
    authService := service.NewAuthService(authRepo)
    
    // 4. Initialize middleware
    permMiddleware := middleware.NewPermissionMiddleware(db)
    roleMiddleware := middleware.NewRoleMiddleware(db)
    
    // 5. Setup routes
    route.SetupRoutes(app, authService, permMiddleware, roleMiddleware)
}
```

## Keuntungan Arsitektur Ini

### 1. Sederhana
- Tidak ada layer handler terpisah
- Service langsung handle HTTP request
- Lebih sedikit boilerplate code

### 2. Mudah Dipahami
- Flow jelas: Route → Middleware → Service → Repository
- Satu file service untuk satu resource

### 3. Mudah Testing
- Repository bisa di-mock
- Service bisa di-test dengan mock repository

### 4. Scalable
- Mudah menambah endpoint baru
- Mudah menambah middleware baru
- Mudah menambah service baru

## Cara Menambah Fitur Baru

### 1. Buat Model
```go
// model/achievement.go
type Achievement struct {
    ID          uuid.UUID `json:"id"`
    StudentID   uuid.UUID `json:"student_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
}
```

### 2. Buat Repository
```go
// repository/achievementRepo.go
type AchievementRepository struct {
    DB *sql.DB
}

func (r *AchievementRepository) GetAll() ([]model.Achievement, error) {
    // Query database
}
```

### 3. Buat Service
```go
// service/achievementService.go
type AchievementService struct {
    AchievementRepo *repository.AchievementRepository
}

func (s *AchievementService) List(c *fiber.Ctx) error {
    achievements, err := s.AchievementRepo.GetAll()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"data": achievements})
}
```

### 4. Register Route
```go
// Routes/Router.go
achievements := api.Group("/achievements")
achievements.Get("/", 
    permMiddleware.RequirePermission("achievements", "read"),
    achievementService.List,
)
```

### 5. Inject di main.go
```go
// main.go
achievementRepo := &repository.AchievementRepository{DB: db}
achievementService := service.NewAchievementService(achievementRepo)
route.SetupRoutes(app, authService, achievementService, ...)
```

## Best Practices

1. **Service method harus menerima `*fiber.Ctx`**
   ```go
   func (s *Service) Method(c *fiber.Ctx) error
   ```

2. **Repository method tidak boleh menerima `*fiber.Ctx`**
   ```go
   func (r *Repository) Method(params) (result, error)
   ```

3. **Validasi input di service layer**
   ```go
   if req.Field == "" {
       return c.Status(400).JSON(fiber.Map{"error": "field required"})
   }
   ```

4. **Error handling yang jelas**
   ```go
   if err != nil {
       return c.Status(500).JSON(fiber.Map{"error": err.Error()})
   }
   ```

5. **Gunakan middleware helper untuk ambil user info**
   ```go
   userID := middleware.GetUserID(c)
   username := middleware.GetUsername(c)
   ```
