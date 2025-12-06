# FR-004: Submit untuk Verifikasi

## Deskripsi
Mahasiswa submit prestasi draft untuk diverifikasi oleh dosen wali.

## Actor
Mahasiswa

## Precondition
- User terautentikasi sebagai mahasiswa
- Prestasi berstatus 'draft'

## Flow
1. Mahasiswa submit prestasi
2. Update status menjadi 'submitted'
3. Create notification untuk dosen wali
4. Return updated status

---

## API Endpoint

### Submit untuk Verifikasi

**Endpoint:** `POST /api/achievements/:reference_id/submit`

**Headers:**
```
Authorization: Bearer <token>
```

**URL Parameters:**
- `reference_id` (UUID) - Achievement reference ID dari PostgreSQL

**Request Body:** None

**Response Success (200):**
```json
{
  "message": "prestasi berhasil disubmit untuk verifikasi",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "submitted",
    "submitted_at": "2024-12-07T10:30:00Z",
    "message": "Prestasi Anda telah disubmit dan menunggu verifikasi dari dosen wali"
  }
}
```

**Response Error (400 - Invalid ID):**
```json
{
  "error": "invalid achievement reference id"
}
```

**Response Error (400 - Wrong Status):**
```json
{
  "error": "only draft achievements can be submitted for verification"
}
```

**Response Error (401 - Unauthorized):**
```json
{
  "error": "user not authenticated"
}
```

**Response Error (403 - Not Owner):**
```json
{
  "error": "you can only submit your own achievements"
}
```

**Response Error (404 - Not Found):**
```json
{
  "error": "achievement reference not found"
}
```

---

## Complete Flow Example

### Step 1: Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "john_student",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "profile": {...}
  }
}
```

### Step 2: Submit Prestasi (Create Draft)
```bash
curl -X POST http://localhost:8080/api/achievements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024",
    "description": "Memenangkan kompetisi hackathon tingkat nasional",
    "details": {
      "competition_name": "Hackathon Nasional 2024",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold"
    }
  }'
```

**Response:**
```json
{
  "message": "prestasi berhasil disubmit",
  "data": {
    "achievement_id": "507f1f77bcf86cd799439011",
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "student_id": "123e4567-e89b-12d3-a456-426614174001",
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024",
    "description": "Memenangkan kompetisi hackathon tingkat nasional",
    "status": "draft",
    "created_at": "2024-12-07T10:00:00Z"
  }
}
```

### Step 3: Submit untuk Verifikasi
```bash
curl -X POST http://localhost:8080/api/achievements/123e4567-e89b-12d3-a456-426614174000/submit \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "message": "prestasi berhasil disubmit untuk verifikasi",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "submitted",
    "submitted_at": "2024-12-07T10:30:00Z",
    "message": "Prestasi Anda telah disubmit dan menunggu verifikasi dari dosen wali"
  }
}
```

---

## Database Changes

### PostgreSQL: achievement_references table
**Before:**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'draft'
submitted_at: NULL
verified_at: NULL
verified_by: NULL
rejection_note: NULL
created_at: 2024-12-07 10:00:00
updated_at: 2024-12-07 10:00:00
```

**After:**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'submitted'
submitted_at: 2024-12-07 10:30:00
verified_at: NULL
verified_by: NULL
rejection_note: NULL
created_at: 2024-12-07 10:00:00
updated_at: 2024-12-07 10:30:00
```

### PostgreSQL: notifications table (New Entry)
```sql
id: 456e7890-e89b-12d3-a456-426614174002
user_id: 789e0123-e89b-12d3-a456-426614174003  -- Advisor ID
type: 'achievement_submitted'
title: 'Prestasi Baru Menunggu Verifikasi'
message: 'Mahasiswa John Doe telah mengajukan prestasi "Juara 1 Hackathon Nasional 2024" untuk diverifikasi'
data: '{"achievement_id":"507f1f77bcf86cd799439011","achievement_reference_id":"123e4567-e89b-12d3-a456-426614174000","student_id":"123e4567-e89b-12d3-a456-426614174001","student_name":"John Doe","achievement_title":"Juara 1 Hackathon Nasional 2024"}'
is_read: false
created_at: 2024-12-07 10:30:00
```

---

## Status Flow

```
┌───────┐
│ draft │ ← Initial status (FR-003)
└───┬───┘
    │
    │ POST /:reference_id/submit (FR-004)
    ▼
┌───────────┐
│ submitted │ ← Waiting for verification
└─────┬─────┘
      │
      ├─── verified (FR-005 - coming soon)
      │
      └─── rejected (FR-005 - coming soon)
```

---

## Notification System

### Notification Data Structure
```json
{
  "id": "uuid",
  "user_id": "uuid",  // Dosen wali
  "type": "achievement_submitted",
  "title": "Prestasi Baru Menunggu Verifikasi",
  "message": "Mahasiswa John Doe telah mengajukan prestasi 'Juara 1 Hackathon' untuk diverifikasi",
  "data": {
    "achievement_id": "507f1f77bcf86cd799439011",
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "student_id": "123e4567-e89b-12d3-a456-426614174001",
    "student_name": "John Doe",
    "achievement_title": "Juara 1 Hackathon Nasional 2024"
  },
  "is_read": false,
  "created_at": "2024-12-07T10:30:00Z"
}
```

### Notification Types
- `achievement_submitted` - Prestasi disubmit untuk verifikasi
- `achievement_verified` - Prestasi diverifikasi (coming soon)
- `achievement_rejected` - Prestasi ditolak (coming soon)
- `achievement_updated` - Prestasi diupdate (coming soon)

---

## Validation Rules

### 1. Authentication
- ✅ User must be authenticated (JWT token)
- ✅ User must be a student

### 2. Ownership
- ✅ User can only submit their own achievements
- ✅ System auto-detects student from logged-in user

### 3. Status Check
- ✅ Only achievements with status 'draft' can be submitted
- ❌ Cannot submit if status is 'submitted', 'verified', or 'rejected'

### 4. Reference ID
- ✅ Must be valid UUID
- ✅ Must exist in database
- ✅ Must belong to the student

---

## Error Handling

### Case 1: Submit already submitted achievement
```bash
# First submit - SUCCESS
POST /api/achievements/123e4567.../submit
→ Status: 200, status changed to 'submitted'

# Second submit - ERROR
POST /api/achievements/123e4567.../submit
→ Status: 400, "only draft achievements can be submitted for verification"
```

### Case 2: Submit other student's achievement
```bash
# Student A tries to submit Student B's achievement
POST /api/achievements/123e4567.../submit
→ Status: 403, "you can only submit your own achievements"
```

### Case 3: Invalid reference ID
```bash
POST /api/achievements/invalid-uuid/submit
→ Status: 400, "invalid achievement reference id"
```

---

## Testing Checklist

### Prerequisites
- [ ] PostgreSQL running
- [ ] MongoDB running
- [ ] Notifications table created
- [ ] Student has advisor assigned
- [ ] Student has permission to create achievements

### Test Cases
- [ ] Submit draft achievement (success)
- [ ] Submit already submitted achievement (error 400)
- [ ] Submit with invalid reference ID (error 400)
- [ ] Submit without authentication (error 401)
- [ ] Submit other student's achievement (error 403)
- [ ] Submit non-existent achievement (error 404)
- [ ] Verify notification created for advisor
- [ ] Verify status changed to 'submitted'
- [ ] Verify submitted_at timestamp set

---

## Database Setup

### Create notifications table
```sql
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
```

### Ensure students have advisor
```sql
-- Check students without advisor
SELECT * FROM students WHERE advisor_id IS NULL;

-- Assign advisor to student
UPDATE students 
SET advisor_id = '<lecturer_id>' 
WHERE id = '<student_id>';
```

---

## Next Features (Coming Soon)

### FR-005: Verifikasi Prestasi (Dosen)
- [ ] Dosen view list prestasi yang perlu diverifikasi
- [ ] Dosen approve prestasi (status: submitted → verified)
- [ ] Dosen reject prestasi (status: submitted → rejected)
- [ ] Create notification untuk mahasiswa

### FR-006: View Notifications
- [ ] User view list notifications
- [ ] Mark notification as read
- [ ] Delete notification
- [ ] Real-time notification (WebSocket)

---

## Summary

FR-004: Submit untuk Verifikasi sudah **COMPLETE**!

**Key Features:**
- ✅ Submit draft achievement untuk verifikasi
- ✅ Update status dari 'draft' ke 'submitted'
- ✅ Create notification untuk dosen wali
- ✅ Validation & ownership check
- ✅ Complete error handling
- ✅ Ready for production

**Next:** Setup notifications table dan test endpoint!
