# Sistem Pelaporan Prestasi Mahasiswa
## Backend API dengan Go Fiber

---

## 📋 Overview Project

**Sistem Pelaporan Prestasi Mahasiswa** adalah aplikasi backend yang memungkinkan mahasiswa melaporkan prestasi mereka, dosen wali memverifikasi prestasi, dan admin mengelola seluruh sistem. Sistem ini dibangun menggunakan **Go Fiber** dengan arsitektur yang scalable dan secure.

---

## 🎯 Tujuan Sistem

1. **Digitalisasi** proses pelaporan prestasi mahasiswa
2. **Otomatisasi** workflow verifikasi oleh dosen wali
3. **Centralized management** untuk admin
4. **Real-time statistics** dan reporting
5. **Audit trail** untuk semua aktivitas

---

## 👥 User Roles & Permissions

### 1. **Mahasiswa (Student)**
- ✅ Submit prestasi baru (draft)
- ✅ Edit prestasi yang masih draft
- ✅ Submit prestasi untuk verifikasi
- ✅ Hapus prestasi draft (soft delete)
- ✅ Lihat prestasi sendiri
- ✅ Lihat statistik prestasi sendiri

### 2. **Dosen Wali (Lecturer)**
- ✅ Lihat prestasi mahasiswa bimbingan
- ✅ Approve/verify prestasi mahasiswa
- ✅ Reject prestasi dengan catatan
- ✅ Lihat statistik mahasiswa bimbingan

### 3. **Admin**
- ✅ Kelola semua user (CRUD)
- ✅ Lihat semua prestasi di sistem
- ✅ Assign dosen wali ke mahasiswa
- ✅ Lihat statistik global sistem
- ✅ Generate reports

---

## 🏗️ Arsitektur Sistem

### **Tech Stack:**
- **Backend Framework:** Go Fiber (Fast HTTP framework)
- **Database:** PostgreSQL (relational data) + MongoDB (flexible documents)
- **Authentication:** JWT (JSON Web Token)
- **Authorization:** RBAC (Role-Based Access Control)
- **Documentation:** Swagger OpenAPI 3.0
- **Testing:** Go testing dengan mocks
- **Deployment:** Docker + GitHub Actions CI/CD

### **Dual Database Architecture:**
```
PostgreSQL (Structured Data)     MongoDB (Flexible Data)
├── Users & Roles               ├── Achievement Details
├── Achievement References      ├── Dynamic Fields
├── Notifications              ├── File Attachments
└── Audit Logs                 └── Rich Content
```

**Mengapa Dual Database?**
- **PostgreSQL:** Untuk data terstruktur, relasi, dan konsistensi ACID
- **MongoDB:** Untuk data prestasi yang fleksibel (setiap prestasi bisa punya field berbeda)

---

## 🔐 Security Features

### **Authentication & Authorization:**
- ✅ **JWT Token** dengan expiration time
- ✅ **Refresh Token** untuk session management
- ✅ **RBAC Middleware** untuk permission checking
- ✅ **Password Hashing** dengan bcrypt
- ✅ **Token Blacklisting** untuk logout

### **Data Protection:**
- ✅ **Input Validation** di semua endpoint
- ✅ **SQL Injection Prevention** dengan prepared statements
- ✅ **CORS Configuration** untuk web security
- ✅ **Rate Limiting** untuk API protection

---

## 📊 Fitur Utama (Functional Requirements)

### **FR-001: Login dengan JWT Authentication**
- Login dengan username/email + password
- Generate JWT token dengan role & permissions
- Return user profile tanpa password

### **FR-002: RBAC Middleware System**
- Token validation middleware
- Permission checking middleware
- Role requirement middleware
- Helper functions untuk authorization

### **FR-003: Submit Prestasi (Mahasiswa)**
- Create prestasi baru dengan status "draft"
- Dual database storage (PostgreSQL + MongoDB)
- Flexible schema untuk berbagai jenis prestasi

### **FR-004: Submit untuk Verifikasi**
- Update status dari "draft" ke "submitted"
- Kirim notifikasi ke dosen wali
- Set timestamp submitted_at

### **FR-005: Hapus Prestasi (Soft Delete)**
- Soft delete untuk audit trail
- Hanya prestasi draft yang bisa dihapus
- Update status di kedua database

### **FR-006: View Prestasi Mahasiswa Bimbingan (Dosen)**
- Dosen lihat prestasi mahasiswa bimbingannya
- Pagination dan filtering
- Batch fetching untuk performance

### **FR-007: Verify Prestasi (Approve)**
- Dosen approve prestasi mahasiswa
- Update status ke "verified"
- Kirim notifikasi ke mahasiswa

### **FR-008: Reject Prestasi**
- Dosen reject dengan catatan
- Update status ke "rejected"
- Kirim notifikasi dengan alasan penolakan

### **FR-009: Manage Users (Admin)**
- CRUD operations untuk users
- Create student/lecturer profiles
- Assign advisor relationships

### **FR-010: View All Achievements (Admin)**
- Admin lihat semua prestasi
- Advanced filtering dan sorting
- Export capabilities

### **FR-011: Achievement Statistics**
- Multi-role statistics (student, lecturer, admin)
- Total by type, period, top students
- Competition level distribution

---

## 🛠️ API Endpoints Structure

### **Base URL:** `/api/v1/`

### **Authentication:**
```
POST /api/v1/auth/login      - User login
POST /api/v1/auth/refresh    - Refresh token
POST /api/v1/auth/logout     - User logout
GET  /api/v1/auth/profile    - Get user profile
```

### **Users (Admin):**
```
GET    /api/v1/users         - Get all users
GET    /api/v1/users/:id     - Get user by ID
POST   /api/v1/users         - Create user
PUT    /api/v1/users/:id     - Update user
DELETE /api/v1/users/:id     - Delete user
PUT    /api/v1/users/:id/role - Update user role
```

### **Achievements:**
```
GET    /api/v1/achievements           - List achievements
GET    /api/v1/achievements/:id       - Get achievement detail
POST   /api/v1/achievements           - Create achievement
PUT    /api/v1/achievements/:id       - Update achievement
DELETE /api/v1/achievements/:id       - Delete achievement
POST   /api/v1/achievements/:id/submit - Submit for verification
POST   /api/v1/achievements/:id/verify - Verify achievement
POST   /api/v1/achievements/:id/reject - Reject achievement
```

### **Statistics & Reports:**
```
GET /api/v1/reports/statistics    - Get statistics
GET /api/v1/reports/student/:id   - Get student report
```

---

## 🧪 Testing Strategy

### **Unit Testing:**
- ✅ **Mock Dependencies:** Database, HTTP, External services
- ✅ **Service Layer Tests:** Business logic testing
- ✅ **Repository Layer Tests:** Data access testing
- ✅ **Middleware Tests:** Authentication & authorization
- ✅ **Test Fixtures:** Reusable test data
- ✅ **Coverage Reports:** Comprehensive test coverage

### **Test Structure:**
```
tests/
├── mocks/              # Mock implementations
├── unit/               # Unit tests
├── integration/        # Integration tests
├── fixtures/           # Test data
└── helpers/            # Test utilities
```

---

## 📚 Documentation

### **Swagger API Documentation:**
- ✅ **Interactive UI** di `/swagger/index.html`
- ✅ **OpenAPI 3.0** specification
- ✅ **Request/Response** examples
- ✅ **Authentication** documentation
- ✅ **Error Handling** examples

### **Code Documentation:**
- ✅ **Go Comments** untuk semua functions
- ✅ **README.md** dengan setup instructions
- ✅ **API Documentation** dengan examples
- ✅ **Database Schema** documentation

---

## 🚀 Deployment & DevOps

### **Docker Configuration:**
- ✅ **Multi-stage build** untuk optimized image
- ✅ **Docker Compose** untuk development
- ✅ **Environment variables** configuration
- ✅ **Health checks** untuk monitoring

### **CI/CD Pipeline (GitHub Actions):**
- ✅ **Automated Testing** pada setiap push
- ✅ **Code Quality Checks** dengan linting
- ✅ **Security Scanning** untuk vulnerabilities
- ✅ **Automated Deployment** ke staging/production

---

## 📈 Performance & Scalability

### **Database Optimization:**
- ✅ **Indexing** pada frequently queried fields
- ✅ **Connection Pooling** untuk database connections
- ✅ **Batch Operations** untuk bulk data processing
- ✅ **Pagination** untuk large datasets

### **API Performance:**
- ✅ **Caching** untuk frequently accessed data
- ✅ **Rate Limiting** untuk API protection
- ✅ **Compression** untuk response optimization
- ✅ **Monitoring** dengan metrics collection

---

## 🔍 Monitoring & Logging

### **Application Monitoring:**
- ✅ **Structured Logging** dengan log levels
- ✅ **Error Tracking** dengan stack traces
- ✅ **Performance Metrics** collection
- ✅ **Health Check** endpoints

### **Security Monitoring:**
- ✅ **Authentication Logs** untuk security audit
- ✅ **Failed Login Attempts** tracking
- ✅ **Permission Violations** logging
- ✅ **Suspicious Activity** detection

---

## 🎯 Key Benefits

### **Untuk Mahasiswa:**
- ✅ **Mudah submit** prestasi kapan saja
- ✅ **Track status** verifikasi real-time
- ✅ **Lihat statistik** prestasi sendiri
- ✅ **History lengkap** semua prestasi

### **Untuk Dosen Wali:**
- ✅ **Centralized view** prestasi mahasiswa bimbingan
- ✅ **Efficient verification** process
- ✅ **Detailed statistics** untuk evaluasi
- ✅ **Notification system** untuk update

### **Untuk Admin:**
- ✅ **Complete system control** dan management
- ✅ **Comprehensive reporting** dan analytics
- ✅ **User management** yang mudah
- ✅ **System monitoring** dan maintenance

---

## 🏆 Technical Achievements

1. **Clean Architecture** dengan separation of concerns
2. **Scalable Design** yang bisa handle growth
3. **Security Best Practices** implementation
4. **Comprehensive Testing** dengan high coverage
5. **Production-Ready** dengan monitoring & logging
6. **Developer-Friendly** dengan good documentation
7. **Modern Tech Stack** dengan industry standards

---

## 📝 Kesimpulan

Sistem Pelaporan Prestasi Mahasiswa ini adalah **production-ready backend API** yang:

- ✅ **Memenuhi semua requirement** sesuai SRS
- ✅ **Menggunakan best practices** dalam development
- ✅ **Scalable dan maintainable** untuk jangka panjang
- ✅ **Secure dan reliable** untuk production use
- ✅ **Well-documented** untuk maintenance
- ✅ **Fully tested** dengan comprehensive test suite

Sistem ini siap untuk **deployment production** dan dapat **di-scale** sesuai kebutuhan institusi.

---

**Developed with ❤️ using Go Fiber**