# FR-009: Manage Users (Admin)

## Deskripsi
Admin dapat melakukan CRUD users, assign roles, set student/lecturer profile, dan set advisor untuk mahasiswa.

## Actor
Admin, Super Admin

## Precondition
- User terautentikasi sebagai admin atau super_admin
- JWT token valid dengan role admin/super_admin

---

## 1. Create User

**Endpoint:** `POST /api/admin/users`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123",
  "full_name": "John Doe",
  "role_id": "uuid-role-id"
}
```

**Response Success (201):**
```json
{
  "message": "user created successfully",
  "data": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "role_id": "uuid-role-id",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "error": "username, email, password, and full_name are required"
}
```

---

## 2. Get All Users

**Endpoint:** `GET /api/admin/users`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `page` (optional, default: 1)
- `per_page` (optional, default: 10, max: 100)

**Response Success (200):**
```json
{
  "message": "success",
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "john_doe",
        "email": "john@example.com",
        "full_name": "John Doe",
        "role_id": "uuid-role-id",
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "per_page": 10,
      "total_pages": 5,
      "total_items": 50
    }
  }
}
```

---

## 3. Update User

**Endpoint:** `PUT /api/admin/users/:user_id`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "role_id": "uuid-new-role-id",
  "is_active": true
}
```

**Response Success (200):**
```json
{
  "message": "user updated successfully",
  "data": {
    "id": "uuid",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe Updated",
    "role_id": "uuid-new-role-id",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "error": "invalid user id"
}
```

---

## 4. Delete User (Soft Delete)

**Endpoint:** `DELETE /api/admin/users/:user_id`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response Success (200):**
```json
{
  "message": "user deleted successfully"
}
```

**Response Error (400):**
```json
{
  "error": "invalid user id"
}
```

**Note:** Soft delete - set `is_active = false`

---

## 5. Create Student Profile

**Endpoint:** `POST /api/admin/students/profile`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "user_id": "uuid-user-id",
  "student_id_number": "2024001",
  "program_study": "Teknik Informatika",
  "academic_year": "2024/2025",
  "advisor_id": "uuid-lecturer-id"
}
```

**Response Success (201):**
```json
{
  "message": "student profile created successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid-user-id",
    "student_id": "2024001",
    "program_study": "Teknik Informatika",
    "academic_year": "2024/2025",
    "advisor_id": "uuid-lecturer-id",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "error": "student_id_number, program_study, and academic_year are required"
}
```

---

## 6. Create Lecturer Profile

**Endpoint:** `POST /api/admin/lecturers/profile`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "user_id": "uuid-user-id",
  "lecturer_id_number": "L001",
  "department": "Fakultas Teknik"
}
```

**Response Success (201):**
```json
{
  "message": "lecturer profile created successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid-user-id",
    "lecturer_id": "L001",
    "department": "Fakultas Teknik",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "error": "lecturer_id_number and department are required"
}
```

---

## 7. Update Student Advisor

**Endpoint:** `PUT /api/admin/students/:student_id/advisor`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "advisor_id": "uuid-new-lecturer-id"
}
```

**Response Success (200):**
```json
{
  "message": "advisor updated successfully"
}
```

**Response Error (400):**
```json
{
  "error": "invalid student id"
}
```

---

## 8. Get All Roles

**Endpoint:** `GET /api/admin/roles`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response Success (200):**
```json
{
  "message": "success",
  "data": [
    {
      "id": "uuid",
      "name": "admin",
      "description": "Administrator",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "name": "lecturer",
      "description": "Dosen",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "name": "student",
      "description": "Mahasiswa",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## Flow Diagram

### Create User Flow
1. Admin mengirim request create user dengan data lengkap
2. Sistem hash password dengan bcrypt
3. Sistem simpan user ke database
4. Return user data (tanpa password)

### Update User Flow
1. Admin mengirim request update user
2. Sistem validasi user_id
3. Sistem update data user (full_name, role_id, is_active)
4. Return updated user data

### Delete User Flow
1. Admin mengirim request delete user
2. Sistem validasi user_id
3. Sistem soft delete (set is_active = false)
4. Return success message

### Create Profile Flow
1. Admin mengirim request create student/lecturer profile
2. Sistem validasi user_id exists
3. Sistem simpan profile ke tabel students/lecturers
4. Return profile data

### Update Advisor Flow
1. Admin mengirim request update advisor
2. Sistem validasi student_id dan advisor_id
3. Sistem update advisor_id di tabel students
4. Return success message

---

## Security
- Endpoint dilindungi dengan JWT authentication
- Hanya role `admin` dan `super_admin` yang dapat akses
- Password di-hash dengan bcrypt sebelum disimpan
- Soft delete untuk preserve data integrity

## Database Tables
- `users` - User data
- `students` - Student profile
- `lecturers` - Lecturer profile
- `roles` - Role definitions

## Notes
- Password tidak pernah di-return dalam response
- Soft delete digunakan untuk user deletion
- Pagination default: page=1, per_page=10
- Maximum per_page: 100
