# FR-010: View All Achievements (Admin)

## Deskripsi
Admin dapat melihat semua prestasi mahasiswa dengan filter dan pagination.

## Endpoint
```
GET /api/admin/achievements
```

## Headers
```
Authorization: Bearer <jwt_token>
```

## Query Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | integer | No | 1 | Halaman yang ingin ditampilkan |
| per_page | integer | No | 10 | Jumlah item per halaman (max: 100) |
| status | string | No | - | Filter berdasarkan status (draft, submitted, verified, rejected, deleted) |
| achievement_type | string | No | - | Filter berdasarkan tipe prestasi |

## Response Success (200)
```json
{
  "message": "success",
  "data": {
    "achievements": [
      {
        "reference_id": "uuid",
        "achievement_id": "mongodb_object_id",
        "student_id": "uuid",
        "achievement_type": "akademik",
        "title": "Juara 1 Lomba Programming",
        "description": "Deskripsi prestasi",
        "status": "verified",
        "submitted_at": "2024-01-15T10:30:00Z",
        "verified_at": "2024-01-16T14:20:00Z",
        "verified_by": "uuid_lecturer",
        "rejection_note": null,
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-16T14:20:00Z",
        "student_info": {
          "student_id_number": "2021001",
          "full_name": "John Doe",
          "program_study": "Teknik Informatika",
          "academic_year": "2021"
        }
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

## Response Error (401)
```json
{
  "error": "unauthorized"
}
```

## Response Error (403)
```json
{
  "error": "insufficient permissions"
}
```

## Response Error (500)
```json
{
  "error": "failed to get achievement references"
}
```

## Contoh Penggunaan

### 1. Get semua achievements (page 1)
```bash
curl -X GET "http://localhost:8080/api/admin/achievements" \
  -H "Authorization: Bearer <jwt_token>"
```

### 2. Get achievements dengan filter status
```bash
curl -X GET "http://localhost:8080/api/admin/achievements?status=verified&page=1&per_page=20" \
  -H "Authorization: Bearer <jwt_token>"
```

### 3. Get achievements dengan filter tipe
```bash
curl -X GET "http://localhost:8080/api/admin/achievements?achievement_type=akademik" \
  -H "Authorization: Bearer <jwt_token>"
```

### 4. Get achievements dengan multiple filter
```bash
curl -X GET "http://localhost:8080/api/admin/achievements?status=submitted&achievement_type=non-akademik&page=2" \
  -H "Authorization: Bearer <jwt_token>"
```

## Fitur
- ✅ Pagination dengan limit maksimal 100 item per halaman
- ✅ Filter berdasarkan status prestasi
- ✅ Filter berdasarkan tipe prestasi
- ✅ Sorting berdasarkan tanggal dibuat (terbaru dulu)
- ✅ Informasi lengkap mahasiswa (nama, NIM, prodi, angkatan)
- ✅ Batch fetching untuk performa optimal
- ✅ Protected dengan role admin/super_admin

## Catatan
- Endpoint ini hanya bisa diakses oleh user dengan role `admin` atau `super_admin`
- Data diambil dari PostgreSQL (references) dan MongoDB (detail prestasi)
- Menggunakan batch fetching untuk mengoptimalkan query ke database
- Filter achievement_type dilakukan di level aplikasi setelah data diambil dari MongoDB