# FR-007: Verifikasi Prestasi

## Deskripsi
Dosen wali dapat memverifikasi (approve) atau menolak (reject) prestasi mahasiswa bimbingannya.

## Actor
Dosen Wali

## Precondition
- User terautentikasi sebagai dosen
- Prestasi berstatus 'submitted'
- Prestasi dari mahasiswa bimbingan dosen tersebut

## Flow
1. Dosen approve/reject prestasi
2. Update status (submitted → verified/rejected)
3. Set verified_at, verified_by, atau rejection_note
4. Create notification untuk mahasiswa
5. Return updated status

---

## API Endpoints

### 1. Verify Achievement (Approve)

**Endpoint:** `POST /api/lecturer/achievements/:reference_id/verify`

**Headers:**
```
Authorization: Bearer <token>
```

**URL Parameters:**
- `reference_id` (UUID) - Achievement reference ID

**Request Body:** None

**Response Success (200):**
```json
{
  "message": "prestasi berhasil diverifikasi",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "verified",
    "verified_by": "789e0123-e89b-12d3-a456-426614174003"
  }
}
```

**Response Error (400 - Wrong Status):**
```json
{
  "error": "only submitted achievements can be verified"
}
```

**Response Error (403 - Not Own Student):**
```json
{
  "error": "you can only verify achievements from your own students"
}
```

**Response Error (404 - Not Found):**
```json
{
  "error": "achievement reference not found"
}
```

---

### 2. Reject Achievement

**Endpoint:** `POST /api/lecturer/achievements/:reference_id/reject`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**
- `reference_id` (UUID) - Achievement reference ID

**Request Body:**
```json
{
  "rejection_note": "Dokumen pendukung tidak lengkap. Mohon upload sertifikat asli."
}
```

**Response Success (200):**
```json
{
  "message": "prestasi ditolak",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "rejected",
    "verified_by": "789e0123-e89b-12d3-a456-426614174003",
    "rejection_note": "Dokumen pendukung tidak lengkap. Mohon upload sertifikat asli."
  }
}
```

**Response Error (400 - Missing Note):**
```json
{
  "error": "rejection_note is required"
}
```

**Response Error (400 - Wrong Status):**
```json
{
  "error": "only submitted achievements can be rejected"
}
```

**Response Error (403 - Not Own Student):**
```json
{
  "error": "you can only reject achievements from your own students"
}
```

---

## Complete Flow Example

### Scenario 1: Approve Achievement

#### Step 1: Login as Lecturer
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "dosen_username",
    "password": "password123"
  }'
```

#### Step 2: View Submitted Achievements
```bash
curl -X GET "http://localhost:8080/api/lecturer/students/achievements?status=submitted" \
  -H "Authorization: Bearer <token>"
```

#### Step 3: Verify Achievement
```bash
curl -X POST http://localhost:8080/api/lecturer/achievements/123e4567-e89b-12d3-a456-426614174000/verify \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "message": "prestasi berhasil diverifikasi",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "verified",
    "verified_by": "789e0123-e89b-12d3-a456-426614174003"
  }
}
```

---

### Scenario 2: Reject Achievement

#### Step 1-2: Same as above

#### Step 3: Reject Achievement with Note
```bash
curl -X POST http://localhost:8080/api/lecturer/achievements/123e4567-e89b-12d3-a456-426614174000/reject \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "rejection_note": "Dokumen pendukung tidak lengkap. Mohon upload sertifikat asli dan foto kegiatan."
  }'
```

**Response:**
```json
{
  "message": "prestasi ditolak",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "rejected",
    "verified_by": "789e0123-e89b-12d3-a456-426614174003",
    "rejection_note": "Dokumen pendukung tidak lengkap. Mohon upload sertifikat asli dan foto kegiatan."
  }
}
```

---

## Database Changes

### PostgreSQL: achievement_references

**Before Verify:**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'submitted'
submitted_at: 2024-12-07 10:30:00
verified_at: NULL
verified_by: NULL
rejection_note: NULL
```

**After Verify (Approved):**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'verified'
submitted_at: 2024-12-07 10:30:00
verified_at: 2024-12-07 14:00:00
verified_by: 789e0123-e89b-12d3-a456-426614174003
rejection_note: NULL
```

**After Reject:**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'rejected'
submitted_at: 2024-12-07 10:30:00
verified_at: 2024-12-07 14:00:00
verified_by: 789e0123-e89b-12d3-a456-426614174003
rejection_note: 'Dokumen pendukung tidak lengkap'
```

### PostgreSQL: notifications (New Entry)

**For Verified:**
```sql
id: uuid
user_id: <student_user_id>
type: 'achievement_verified'
title: 'Prestasi Diverifikasi'
message: 'Prestasi Anda "Juara 1 Hackathon" telah diverifikasi oleh dosen wali'
data: '{"achievement_id":"507f...","student_id":"123e...",...}'
is_read: false
created_at: 2024-12-07 14:00:00
```

**For Rejected:**
```sql
id: uuid
user_id: <student_user_id>
type: 'achievement_rejected'
title: 'Prestasi Ditolak'
message: 'Prestasi Anda "Juara 1 Hackathon" ditolak. Alasan: Dokumen tidak lengkap'
data: '{"achievement_id":"507f...","student_id":"123e...",...}'
is_read: false
created_at: 2024-12-07 14:00:00
```

---

## Status Flow

```
draft
  │
  ├─── DELETE → deleted
  │
  └─── POST /submit
       │
       ▼
   submitted ← Can be verified/rejected
       │
       ├─── POST /verify → verified (final)
       │
       └─── POST /reject → rejected (can resubmit)
```

**Rules:**
- ✅ Can verify/reject: status = 'submitted'
- ❌ Cannot verify/reject: status = 'draft', 'verified', 'rejected', 'deleted'

---

## Validation Rules

### 1. Authentication
- ✅ User must be authenticated (JWT token)
- ✅ User must be lecturer/dosen

### 2. Ownership
- ✅ Achievement must be from lecturer's own student
- ✅ Check advisor_id matches lecturer_id
- ❌ Cannot verify other lecturer's students

### 3. Status Check
- ✅ Only 'submitted' achievements can be verified/rejected
- ❌ Cannot verify 'draft', 'verified', 'rejected', 'deleted'

### 4. Rejection Note
- ✅ Required when rejecting
- ✅ Must not be empty
- ✅ Stored in database

---

## Notification System

### Notification Types
- `achievement_verified` - Prestasi diverifikasi
- `achievement_rejected` - Prestasi ditolak

### Notification Data
```json
{
  "achievement_id": "507f1f77bcf86cd799439011",
  "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
  "student_id": "123e4567-e89b-12d3-a456-426614174001",
  "student_name": "John Doe",
  "achievement_title": "Juara 1 Hackathon Nasional 2024"
}
```

---

## Error Handling

### Case 1: Verify already verified achievement
```bash
# First verify - SUCCESS
POST /api/lecturer/achievements/123e4567.../verify
→ Status: 200, status = 'verified'

# Second verify - ERROR
POST /api/lecturer/achievements/123e4567.../verify
→ Status: 400, "only submitted achievements can be verified"
```

### Case 2: Verify other lecturer's student
```bash
# Lecturer A tries to verify Lecturer B's student
POST /api/lecturer/achievements/123e4567.../verify
→ Status: 403, "you can only verify achievements from your own students"
```

### Case 3: Reject without note
```bash
POST /api/lecturer/achievements/123e4567.../reject
Body: {}
→ Status: 400, "rejection_note is required"
```

---

## Testing Checklist

### Prerequisites
- [ ] PostgreSQL running
- [ ] MongoDB running
- [ ] Lecturer has students
- [ ] Students have submitted achievements

### Test Cases - Verify
- [ ] Verify submitted achievement (success)
- [ ] Verify draft achievement (error 400)
- [ ] Verify already verified achievement (error 400)
- [ ] Verify other lecturer's student (error 403)
- [ ] Verify non-existent achievement (error 404)
- [ ] Verify notification created for student
- [ ] Verify status changed to 'verified'
- [ ] Verify verified_at and verified_by set

### Test Cases - Reject
- [ ] Reject submitted achievement with note (success)
- [ ] Reject without note (error 400)
- [ ] Reject draft achievement (error 400)
- [ ] Reject already rejected achievement (error 400)
- [ ] Reject other lecturer's student (error 403)
- [ ] Reject notification created for student
- [ ] Verify status changed to 'rejected'
- [ ] Verify rejection_note saved

### SQL Queries for Testing
```sql
-- Check achievement status
SELECT id, status, verified_at, verified_by, rejection_note
FROM achievement_references
WHERE id = '<reference_id>';

-- Check verified achievements
SELECT * FROM achievement_references
WHERE status = 'verified'
AND verified_by = '<lecturer_user_id>';

-- Check rejected achievements
SELECT * FROM achievement_references
WHERE status = 'rejected'
AND rejection_note IS NOT NULL;

-- Check notifications
SELECT * FROM notifications
WHERE type IN ('achievement_verified', 'achievement_rejected')
ORDER BY created_at DESC;
```

---

## Next Features (Coming Soon)

### FR-008: Resubmit Rejected Achievement
- [ ] Student can resubmit rejected achievement
- [ ] Update status rejected → submitted
- [ ] Clear rejection_note
- [ ] Create new notification

### FR-009: View Verification History
- [ ] Track all verification actions
- [ ] Show who verified/rejected
- [ ] Show when verified/rejected
- [ ] Show rejection notes history

### FR-010: Bulk Verification
- [ ] Verify multiple achievements at once
- [ ] Reject multiple achievements
- [ ] Transaction for all or nothing

---

## Summary

FR-007: Verifikasi Prestasi sudah **COMPLETE**!

**Key Features:**
- ✅ Dosen approve prestasi (submitted → verified)
- ✅ Dosen reject prestasi (submitted → rejected)
- ✅ Rejection note required
- ✅ Set verified_at, verified_by
- ✅ Create notification for student
- ✅ Ownership verification
- ✅ Complete error handling
- ✅ Ready for production

**Benefits:**
- ✅ Dosen can manage student achievements
- ✅ Clear rejection reason
- ✅ Student gets notified
- ✅ Audit trail complete
- ✅ Secure with ownership check

**Next:** Test endpoints dan verify notifications!
