# Install Dependencies untuk FR-003

## MongoDB Driver

Untuk menggunakan fitur Submit Prestasi (FR-003), Anda perlu install MongoDB driver:

```bash
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/bson/primitive
```

## Atau install semua dependencies sekaligus:

```bash
go mod tidy
```

## Setup MongoDB

### 1. Install MongoDB
- Download dari: https://www.mongodb.com/try/download/community
- Atau gunakan Docker:
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 2. Update .env file
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=uas_prestasi
```

### 3. Verify MongoDB Connection
```bash
# Jika menggunakan MongoDB Compass
mongodb://localhost:27017

# Atau menggunakan mongosh
mongosh mongodb://localhost:27017
```

## Setup PostgreSQL Tables

Pastikan tabel-tabel berikut sudah dibuat:

### 1. Table: users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 2. Table: roles
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 3. Table: permissions
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);
```

### 4. Table: role_permissions
```sql
CREATE TABLE role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
```

### 5. Table: students
```sql
CREATE TABLE students (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    student_id VARCHAR(20) UNIQUE NOT NULL,
    program_study VARCHAR(100),
    academic_year VARCHAR(10),
    advisor_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (advisor_id) REFERENCES lecturers(id)
);
```

### 6. Table: lecturers
```sql
CREATE TABLE lecturers (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    lecturer_id VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 7. Table: achievement_references
```sql
CREATE TABLE achievement_references (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    mongo_achievement_id VARCHAR(24) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'submitted', 'verified', 'rejected')),
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID,
    rejection_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (verified_by) REFERENCES users(id)
);
```

## Seed Data untuk Testing

### 1. Insert Role Student
```sql
INSERT INTO roles (id, name, description) 
VALUES ('550e8400-e29b-41d4-a716-446655440001', 'student', 'Student role');
```

### 2. Insert Permission untuk Student
```sql
-- Permission: create achievement
INSERT INTO permissions (id, name, resource, action, description)
VALUES ('550e8400-e29b-41d4-a716-446655440010', 'create_achievement', 'achievements', 'create', 'Create achievement');

-- Permission: read own achievement
INSERT INTO permissions (id, name, resource, action, description)
VALUES ('550e8400-e29b-41d4-a716-446655440011', 'read_achievement', 'achievements', 'read', 'Read own achievement');
```

### 3. Assign Permission ke Role Student
```sql
INSERT INTO role_permissions (role_id, permission_id)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440010'),
    ('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440011');
```

### 4. Insert User Student
```sql
-- Password: password123 (hashed dengan bcrypt)
INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
VALUES (
    '550e8400-e29b-41d4-a716-446655440100',
    'john_student',
    'john@student.com',
    '$2a$10$YourBcryptHashHere',
    'John Doe',
    '550e8400-e29b-41d4-a716-446655440001',
    true
);
```

### 5. Insert Student Data
```sql
INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id)
VALUES (
    '550e8400-e29b-41d4-a716-446655440200',
    '550e8400-e29b-41d4-a716-446655440100',
    '2021001',
    'Teknik Informatika',
    '2021',
    NULL
);
```

## Run Application

```bash
go run main.go
```

Server akan berjalan di: http://localhost:8080

## Test Endpoints

### 1. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "john_student",
    "password": "password123"
  }'
```

### 2. Submit Achievement
```bash
curl -X POST http://localhost:8080/api/achievements \
  -H "Authorization: Bearer <token_from_login>" \
  -H "Content-Type: application/json" \
  -d '{
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon",
    "description": "Memenangkan hackathon nasional",
    "details": {
      "competition_name": "Hackathon Nasional 2024",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold"
    }
  }'
```
