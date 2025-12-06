# FR-003: Submit Prestasi

## Deskripsi
Mahasiswa dapat menambahkan laporan prestasi dengan menyimpan data ke MongoDB (achievement) dan PostgreSQL (reference).

## Actor
Mahasiswa

## Precondition
- User terautentikasi sebagai mahasiswa
- User memiliki permission `achievements:create`

## Flow
1. Mahasiswa mengisi data prestasi
2. Mahasiswa upload dokumen pendukung (optional)
3. Sistem simpan ke MongoDB (achievement) dan PostgreSQL (reference)
4. Status awal: 'draft'
5. Return achievement data

---

## API Endpoints

### 1. Submit Prestasi (Create)

**Endpoint:** `POST /api/achievements`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "achievement_type": "competition",
  "title": "Juara 1 Hackathon Nasional 2024",
  "description": "Memenangkan kompetisi hackathon tingkat nasional dengan tema AI for Education",
  "details": {
    "competition_name": "Hackathon Nasional 2024",
    "competition_level": "national",
    "rank": 1,
    "medal_type": "gold"
  }
}
```

**Achievement Types:**
- `academic` - Prestasi akademik
- `competition` - Kompetisi/lomba
- `organization` - Organisasi
- `publication` - Publikasi
- `certification` - Sertifikasi
- `other` - Lainnya

**Details untuk Competition:**
```json
{
  "competition_name": "string (optional)",
  "competition_level": "international|national|regional|local (optional)",
  "rank": "number (optional)",
  "medal_type": "string (optional)"
}
```

**Details untuk Publication:**
```json
{
  "publication_type": "journal|conference|book (optional)",
  "publication_title": "string (optional)",
  "authors": ["string"] (optional),
  "publisher": "string (optional)",
  "issn": "string (optional)"
}
```

**Response Success (201):**
```json
{
  "message": "prestasi berhasil disubmit",
  "data": {
    "achievement_id": "507f1f77bcf86cd799439011",
    "achievement_reference_id": "123e4567-e89b-12d3-a456-426614174000",
    "student_id": "123e4567-e89b-12d3-a456-426614174001",
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024",
    "description": "Memenangkan kompetisi hackathon tingkat nasional dengan tema AI for Education",
    "status": "draft",
    "created_at": "2024-12-06T10:30:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "error": "achievement_type harus diisi"
}
```

**Response Error (401):**
```json
{
  "error": "user not authenticated"
}
```

**Response Error (403):**
```json
{
  "error": "forbidden",
  "message": "you don't have permission to create achievements"
}
```

**Response Error (404):**
```json
{
  "error": "student not found"
}
```

---

### 2. Get My Achievements

**Endpoint:** `GET /api/achievements/my`

**Headers:**
```
Authorization: Bearer <token>
```

**Response Success (200):**
```json
{
  "message": "success",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "student_id": "123e4567-e89b-12d3-a456-426614174001",
      "achievement_type": "competition",
      "title": "Juara 1 Hackathon Nasional 2024",
      "description": "Memenangkan kompetisi hackathon tingkat nasional",
      "details": {
        "competition_name": "Hackathon Nasional 2024",
        "competition_level": "national",
        "rank": 1,
        "medal_type": "gold"
      }
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "student_id": "123e4567-e89b-12d3-a456-426614174001",
      "achievement_type": "publication",
      "title": "Paper tentang Machine Learning",
      "description": "Publikasi paper di jurnal internasional",
      "details": {
        "publication_type": "journal",
        "publication_title": "Advanced ML Techniques",
        "authors": ["John Doe", "Jane Smith"],
        "publisher": "IEEE",
        "issn": "1234-5678"
      }
    }
  ]
}
```

---

### 3. Get Achievement Detail

**Endpoint:** `GET /api/achievements/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Response Success (200):**
```json
{
  "message": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "student_id": "123e4567-e89b-12d3-a456-426614174001",
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024",
    "description": "Memenangkan kompetisi hackathon tingkat nasional",
    "details": {
      "competition_name": "Hackathon Nasional 2024",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold"
    }
  }
}
```

**Response Error (403):**
```json
{
  "error": "you can only view your own achievements"
}
```

**Response Error (404):**
```json
{
  "error": "achievement not found"
}
```

---

## Testing dengan curl

### 1. Login dulu untuk dapat token
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "student_username",
    "password": "password123"
  }'
```

### 2. Submit prestasi competition
```bash
curl -X POST http://localhost:8080/api/achievements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "achievement_type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024",
    "description": "Memenangkan kompetisi hackathon tingkat nasional dengan tema AI for Education",
    "details": {
      "competition_name": "Hackathon Nasional 2024",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold"
    }
  }'
```

### 3. Submit prestasi publication
```bash
curl -X POST http://localhost:8080/api/achievements \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "achievement_type": "publication",
    "title": "Paper tentang Machine Learning",
    "description": "Publikasi paper di jurnal internasional",
    "details": {
      "publication_type": "journal",
      "publication_title": "Advanced ML Techniques",
      "authors": ["John Doe", "Jane Smith"],
      "publisher": "IEEE",
      "issn": "1234-5678"
    }
  }'
```

### 4. Lihat prestasi sendiri
```bash
curl -X GET http://localhost:8080/api/achievements/my \
  -H "Authorization: Bearer <token>"
```

### 5. Lihat detail prestasi
```bash
curl -X GET http://localhost:8080/api/achievements/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token>"
```

---

## Database Schema

### MongoDB Collection: achievements
```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439011"),
  studentId: UUID("123e4567-e89b-12d3-a456-426614174001"),
  achievementType: "competition",
  title: "Juara 1 Hackathon Nasional 2024",
  description: "Memenangkan kompetisi hackathon tingkat nasional",
  details: {
    competitionName: "Hackathon Nasional 2024",
    competitionLevel: "national",
    rank: 1,
    medalType: "gold"
  }
}
```

### PostgreSQL Table: achievement_references
```sql
id: UUID PRIMARY KEY
student_id: UUID FOREIGN KEY -> students.id
mongo_achievement_id: VARCHAR(24) NOT NULL
status: ENUM('draft', 'submitted', 'verified', 'rejected')
submitted_at: TIMESTAMP
verified_at: TIMESTAMP
verified_by: UUID FOREIGN KEY -> users.id
rejection_note: TEXT
created_at: TIMESTAMP DEFAULT NOW()
updated_at: TIMESTAMP DEFAULT NOW()
```

---

## Validasi

### Required Fields:
- `achievement_type` - Harus diisi dan valid
- `title` - Harus diisi
- `description` - Harus diisi

### Valid Achievement Types:
- academic
- competition
- organization
- publication
- certification
- other

### Permission Required:
- `achievements:create` - Untuk submit prestasi

---

## Status Flow

```
draft → submitted → verified
                 ↓
              rejected
```

- **draft**: Status awal saat prestasi dibuat
- **submitted**: Prestasi sudah disubmit untuk verifikasi (coming soon)
- **verified**: Prestasi sudah diverifikasi oleh dosen/admin (coming soon)
- **rejected**: Prestasi ditolak (coming soon)

---

## Notes

1. **Dual Database**: 
   - MongoDB untuk data prestasi (flexible schema)
   - PostgreSQL untuk reference dan status tracking

2. **Ownership**: 
   - Mahasiswa hanya bisa submit prestasi untuk diri sendiri
   - System otomatis ambil student_id dari user yang login

3. **Status**: 
   - Status awal selalu 'draft'
   - Mahasiswa bisa edit selama status masih 'draft'

4. **Details Field**: 
   - Field `details` bersifat flexible
   - Bisa berisi CompetitionDetails atau PublicationDetails
   - Atau custom object sesuai achievement_type

5. **Future Features**:
   - Upload dokumen pendukung
   - Submit untuk verifikasi
   - Edit prestasi (status draft)
   - Delete prestasi (status draft)
