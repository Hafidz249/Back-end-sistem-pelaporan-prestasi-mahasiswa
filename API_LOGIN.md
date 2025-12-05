# FR-001: Login API

## Endpoint
```
POST /api/auth/login
```

## Request Body
```json
{
  "credential": "username_atau_email",
  "password": "password123"
}
```

## Response Success (200)
```json
{
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "profile": {
      "id": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": {
        "id": "uuid",
        "name": "student",
        "description": "Student role"
      },
      "permissions": [
        {
          "id": "uuid",
          "name": "view_achievements",
          "resource": "achievements",
          "action": "read"
        }
      ]
    }
  }
}
```

## Response Error (400 - Validation Error)
```json
{
  "error": "username atau email harus diisi"
}
```

## Response Error (401 - Unauthorized)
```json
{
  "error": "kredensial salah"
}
```
atau
```json
{
  "error": "akun anda dinonaktifkan"
}
```

## Flow
1. User mengirim kredensial (username/email + password)
2. Sistem memvalidasi input
3. Sistem mencari user berdasarkan username ATAU email
4. Sistem mengecek status aktif user
5. Sistem memvalidasi password dengan bcrypt
6. Sistem mengambil role dan permissions dari database
7. Sistem generate JWT token dengan role dan permissions
8. Return token dan user profile (tanpa password)

## JWT Token Claims
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
    }
  ],
  "exp": 1234567890
}
```

## Testing dengan curl
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "john_doe",
    "password": "password123"
  }'
```
