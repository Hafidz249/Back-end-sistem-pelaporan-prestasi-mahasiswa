# FR-002: RBAC Middleware

## Deskripsi
Setiap endpoint dilindungi dengan permission check menggunakan JWT token dan database.

## Flow
1. **Ekstrak JWT dari header** - Ambil token dari Authorization header
2. **Validasi token** - Verify signature dan expiration
3. **Load user permissions** - Ambil dari token (cache) atau database
4. **Check permission** - Validasi apakah user punya permission yang diperlukan
5. **Allow/deny request** - Lanjutkan atau tolak request

## Middleware yang Tersedia

### 1. JWTAuth()
Middleware untuk autentikasi JWT token.

**Usage:**
```go
api := app.Group("/api", middleware.JWTAuth())
```

**Response jika gagal:**
```json
{
  "error": "Authorization header required"
}
```
atau
```json
{
  "error": "invalid or expired token"
}
```

### 2. RequirePermission(resource, action)
Check permission berdasarkan resource dan action.

**Usage:**
```go
achievements.Get("/", 
    permMiddleware.RequirePermission("achievements", "read"),
    handler,
)
```

**Response jika gagal:**
```json
{
  "error": "forbidden",
  "message": "you don't have permission to read achievements"
}
```

### 3. RequirePermissionWithCache(resource, action)
Check permission dari token, fallback ke database jika tidak ada.

**Usage:**
```go
achievements.Post("/", 
    permMiddleware.RequirePermissionWithCache("achievements", "create"),
    handler,
)
```

### 4. RequireAnyPermission(permissions)
User harus punya minimal 1 dari permissions yang diminta.

**Usage:**
```go
permMiddleware.RequireAnyPermission([]middleware.Permission{
    {Resource: "achievements", Action: "read"},
    {Resource: "achievements", Action: "update"},
})
```

### 5. RequireAllPermissions(permissions)
User harus punya semua permissions yang diminta.

**Usage:**
```go
permMiddleware.RequireAllPermissions([]middleware.Permission{
    {Resource: "achievements", Action: "read"},
    {Resource: "students", Action: "read"},
})
```

### 6. RequireRole(roles...)
Check apakah user punya salah satu role yang diizinkan.

**Usage:**
```go
admin := api.Group("/admin", roleMiddleware.RequireRole("admin", "super_admin"))
```

**Response jika gagal:**
```json
{
  "error": "forbidden",
  "message": "insufficient role"
}
```

## Contoh Penggunaan

### Protected Route dengan Permission
```go
// Hanya user dengan permission "achievements:read"
achievements.Get("/", 
    permMiddleware.RequirePermission("achievements", "read"),
    achievementHandler.List,
)

// Hanya user dengan permission "achievements:create"
achievements.Post("/",
    permMiddleware.RequirePermission("achievements", "create"),
    achievementHandler.Create,
)
```

### Protected Route dengan Role
```go
// Hanya admin dan super_admin
admin := api.Group("/admin", 
    roleMiddleware.RequireRole("admin", "super_admin"),
)
admin.Get("/users", userHandler.ListAll)
```

### Kombinasi JWT + Permission
```go
// Semua route di group ini butuh JWT valid
api := app.Group("/api", middleware.JWTAuth())

// Endpoint spesifik butuh permission tertentu
api.Get("/achievements", 
    permMiddleware.RequirePermission("achievements", "read"),
    achievementHandler.List,
)
```

## Testing dengan curl

### 1. Login dulu untuk dapat token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "john_doe",
    "password": "password123"
  }'
```

Response:
```json
{
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "profile": {...}
  }
}
```

### 2. Akses protected endpoint dengan token
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Akses endpoint dengan permission check
```bash
# Berhasil jika user punya permission "achievements:read"
curl -X GET http://localhost:8080/api/achievements \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Gagal jika user tidak punya permission
# Response: {"error": "forbidden", "message": "you don't have permission to read achievements"}
```

### 4. Akses endpoint admin
```bash
# Berhasil jika user role = admin atau super_admin
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Gagal jika user role = student atau lecturer
# Response: {"error": "forbidden", "message": "insufficient role"}
```

## Context Locals yang Tersedia

Setelah melewati `JWTAuth()` middleware, data berikut tersedia di context:

```go
c.Locals("user_id")      // string - UUID user
c.Locals("username")     // string - username
c.Locals("email")        // string - email
c.Locals("role_id")      // string - UUID role
c.Locals("permissions")  // []map[string]interface{} - list permissions
c.Locals("role_name")    // string - nama role (setelah RequireRole)
```

## Permission Format di JWT Token

```json
{
  "user_id": "uuid",
  "username": "john_doe",
  "email": "john@example.com",
  "role_id": "uuid",
  "permissions": [
    {
      "name": "view_achievements",
      "resource": "achievements",
      "action": "read"
    },
    {
      "name": "create_achievements",
      "resource": "achievements",
      "action": "create"
    }
  ],
  "exp": 1234567890
}
```

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```
```json
{
  "error": "invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "forbidden",
  "message": "you don't have permission to read achievements"
}
```
```json
{
  "error": "forbidden",
  "message": "insufficient role"
}
```

## Best Practices

1. **Gunakan RequirePermission untuk granular access control**
   ```go
   // Lebih spesifik dan flexible
   permMiddleware.RequirePermission("achievements", "read")
   ```

2. **Gunakan RequireRole untuk role-based access**
   ```go
   // Untuk endpoint yang hanya boleh diakses role tertentu
   roleMiddleware.RequireRole("admin")
   ```

3. **Kombinasikan JWT + Permission untuk keamanan maksimal**
   ```go
   api := app.Group("/api", middleware.JWTAuth())
   api.Get("/data", permMiddleware.RequirePermission("data", "read"), handler)
   ```

4. **Gunakan RequirePermissionWithCache untuk performa lebih baik**
   - Check dari token dulu (fast)
   - Fallback ke database jika perlu (accurate)
