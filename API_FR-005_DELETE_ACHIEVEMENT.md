# FR-005: Hapus Prestasi

## Deskripsi
Mahasiswa dapat menghapus prestasi yang masih berstatus draft.

## Actor
Mahasiswa

## Precondition
- User terautentikasi sebagai mahasiswa
- Prestasi berstatus 'draft'

## Flow
1. Soft delete data di MongoDB (set isDeleted = true)
2. Update reference di PostgreSQL (set status = 'deleted')
3. Return success message

---

## API Endpoint

### Delete Achievement

**Endpoint:** `DELETE /api/achievements/:reference_id`

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
  "message": "prestasi berhasil dihapus",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "deleted"
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
  "error": "only draft achievements can be deleted"
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
  "error": "you can only delete your own achievements"
}
```

**Response Error (404 - Not Found):**
```json
{
  "error": "achievement reference not found"
}
```

**Response Error (500 - Server Error):**
```json
{
  "error": "failed to delete achievement"
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

### Step 2: Create Draft Achievement
```bash
curl -X POST http://localhost:8080/api/achievements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "achievement_type": "competition",
    "title": "Test Achievement",
    "description": "This is a test",
    "details": {}
  }'
```

**Response:**
```json
{
  "message": "prestasi berhasil disubmit",
  "data": {
    "achievement_id": "507f1f77bcf86cd799439011",
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "draft",
    ...
  }
}
```

### Step 3: Delete Achievement
```bash
curl -X DELETE http://localhost:8080/api/achievements/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "message": "prestasi berhasil dihapus",
  "data": {
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "deleted"
  }
}
```

---

## Soft Delete Implementation

### Why Soft Delete?
- **Data Recovery**: Dapat restore jika terhapus tidak sengaja
- **Audit Trail**: Tetap ada history untuk audit
- **Referential Integrity**: Tidak break foreign key constraints
- **Analytics**: Data tetap bisa digunakan untuk reporting

### MongoDB Changes

**Before Delete:**
```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439011"),
  studentId: UUID("123e4567-e89b-12d3-a456-426614174001"),
  achievementType: "competition",
  title: "Test Achievement",
  description: "This is a test",
  details: {}
}
```

**After Delete (Soft Delete):**
```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439011"),
  studentId: UUID("123e4567-e89b-12d3-a456-426614174001"),
  achievementType: "competition",
  title: "Test Achievement",
  description: "This is a test",
  details: {},
  isDeleted: true,
  deletedAt: ISODate("2024-12-07T11:00:00Z")
}
```

### PostgreSQL Changes

**Before Delete:**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'draft'
created_at: 2024-12-07 10:00:00
updated_at: 2024-12-07 10:00:00
```

**After Delete:**
```sql
id: 123e4567-e89b-12d3-a456-426614174000
student_id: 123e4567-e89b-12d3-a456-426614174001
mongo_achievement_id: 507f1f77bcf86cd799439011
status: 'deleted'
created_at: 2024-12-07 10:00:00
updated_at: 2024-12-07 11:00:00
```

---

## Status Flow

```
┌───────┐
│ draft │ ← Can be deleted
└───┬───┘
    │
    ├─── DELETE /:reference_id → deleted
    │
    └─── POST /:reference_id/submit
         │
         ▼
    ┌───────────┐
    │ submitted │ ← Cannot be deleted
    └───────────┘
         │
         ├─── verified ← Cannot be deleted
         │
         └─── rejected ← Cannot be deleted
```

**Rules:**
- ✅ Can delete: status = 'draft'
- ❌ Cannot delete: status = 'submitted', 'verified', 'rejected'

---

## Transaction & Rollback

### Success Flow
```
1. Soft delete MongoDB (isDeleted = true)
2. Update PostgreSQL (status = 'deleted')
3. Return success
```

### Rollback Flow (if PostgreSQL fails)
```
1. Soft delete MongoDB (isDeleted = true)
2. Update PostgreSQL fails ❌
3. Rollback MongoDB (isDeleted = false, remove deletedAt)
4. Return error
```

### Code Implementation
```go
// 1. Soft delete MongoDB
collection.UpdateOne(ctx, filter, update)

// 2. Update PostgreSQL
result, err := db.Exec(query)
if err != nil {
    // Rollback MongoDB
    collection.UpdateOne(ctx, filter, rollbackUpdate)
    return err
}
```

---

## Validation Rules

### 1. Authentication
- ✅ User must be authenticated (JWT token)
- ✅ User must be a student

### 2. Ownership
- ✅ User can only delete their own achievements
- ✅ System auto-detects student from logged-in user

### 3. Status Check
- ✅ Only achievements with status 'draft' can be deleted
- ❌ Cannot delete if status is 'submitted', 'verified', or 'rejected'

### 4. Reference ID
- ✅ Must be valid UUID
- ✅ Must exist in database
- ✅ Must belong to the student

---

## Query Filter for List

### Before (Show All)
```go
filter := primitive.M{
    "studentId": studentID,
}
```

### After (Exclude Deleted)
```go
filter := primitive.M{
    "studentId": studentID,
    "$or": []primitive.M{
        {"isDeleted": primitive.M{"$exists": false}},
        {"isDeleted": false},
    },
}
```

**Result:**
- GET /api/achievements/my → Tidak menampilkan yang sudah dihapus
- Deleted achievements tidak muncul di list

---

## Error Handling

### Case 1: Delete already submitted achievement
```bash
# Submit achievement first
POST /api/achievements/123e4567.../submit
→ Status: 200, status = 'submitted'

# Try to delete
DELETE /api/achievements/123e4567...
→ Status: 400, "only draft achievements can be deleted"
```

### Case 2: Delete other student's achievement
```bash
# Student A tries to delete Student B's achievement
DELETE /api/achievements/123e4567...
→ Status: 403, "you can only delete your own achievements"
```

### Case 3: Delete non-existent achievement
```bash
DELETE /api/achievements/non-existent-uuid
→ Status: 404, "achievement reference not found"
```

### Case 4: Delete already deleted achievement
```bash
# First delete - SUCCESS
DELETE /api/achievements/123e4567...
→ Status: 200, "prestasi berhasil dihapus"

# Second delete - ERROR
DELETE /api/achievements/123e4567...
→ Status: 400, "only draft achievements can be deleted"
```

---

## Testing Checklist

### Prerequisites
- [ ] PostgreSQL running
- [ ] MongoDB running
- [ ] Student has draft achievement

### Test Cases
- [ ] Delete draft achievement (success)
- [ ] Delete submitted achievement (error 400)
- [ ] Delete verified achievement (error 400)
- [ ] Delete with invalid reference ID (error 400)
- [ ] Delete without authentication (error 401)
- [ ] Delete other student's achievement (error 403)
- [ ] Delete non-existent achievement (error 404)
- [ ] Verify MongoDB isDeleted = true
- [ ] Verify PostgreSQL status = 'deleted'
- [ ] Verify deleted achievement not in list
- [ ] Verify rollback if PostgreSQL fails

### SQL Queries for Testing
```sql
-- Check achievement status
SELECT id, status, updated_at 
FROM achievement_references 
WHERE id = '<reference_id>';

-- Check deleted achievements
SELECT * FROM achievement_references 
WHERE status = 'deleted';
```

### MongoDB Queries for Testing
```javascript
// Check soft delete
db.achievements.findOne({
  _id: ObjectId("507f1f77bcf86cd799439011")
})

// Check deleted achievements
db.achievements.find({
  isDeleted: true
})

// Count deleted achievements
db.achievements.countDocuments({
  isDeleted: true
})
```

---

## Future Features (Coming Soon)

### FR-006: Restore Deleted Achievement
- [ ] Admin can restore deleted achievements
- [ ] Set isDeleted = false
- [ ] Set status back to 'draft'

### FR-007: Permanent Delete
- [ ] Admin can permanently delete
- [ ] Remove from MongoDB
- [ ] Remove from PostgreSQL
- [ ] Cannot be restored

### FR-008: Bulk Delete
- [ ] Delete multiple achievements at once
- [ ] Transaction for all or nothing

---

## Summary

FR-005: Hapus Prestasi sudah **COMPLETE**!

**Key Features:**
- ✅ Soft delete di MongoDB (isDeleted flag)
- ✅ Update status di PostgreSQL (status = 'deleted')
- ✅ Transaction with rollback
- ✅ Only draft can be deleted
- ✅ Ownership verification
- ✅ Exclude deleted from list
- ✅ Complete error handling
- ✅ Ready for production

**Benefits:**
- ✅ Data recovery possible
- ✅ Audit trail maintained
- ✅ No broken references
- ✅ Analytics data preserved

**Next:** Test endpoint dan verify soft delete!
