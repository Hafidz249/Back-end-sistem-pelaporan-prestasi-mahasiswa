# 🎥 Script Presentasi Video - Sistem Pelaporan Prestasi Mahasiswa

## ⏱️ Durasi Estimasi: 15-20 menit

---

## 🎬 PEMBUKAAN (2 menit)

### **Slide 1: Title Screen**
**Yang ditampilkan:** Desktop dengan folder project terbuka
**Yang dijelaskan:**
> "Selamat pagi/siang Bapak/Ibu Dosen. Saya akan mempresentasikan project Sistem Pelaporan Prestasi Mahasiswa yang telah saya kembangkan menggunakan Go Fiber sebagai backend API."

> "Project ini adalah sistem lengkap untuk mengelola prestasi mahasiswa dengan fitur authentication, authorization, dan management yang comprehensive."

### **Slide 2: Project Overview**
**Yang ditampilkan:** Buka file `PENJELASAN_PROJECT.md`
**Yang dijelaskan:**
> "Sistem ini memiliki 3 role utama: Mahasiswa untuk submit prestasi, Dosen Wali untuk verifikasi, dan Admin untuk management. Menggunakan dual database PostgreSQL dan MongoDB untuk fleksibilitas data."

---

## 🏗️ BAGIAN 1: STRUKTUR PROJECT (3 menit)

### **Demo 1.1: Project Structure**
**Yang ditampilkan:** Explorer VS Code dengan folder structure
**Yang dijelaskan:**
> "Mari kita lihat struktur project. Saya menggunakan clean architecture dengan separation of concerns:"

**Tunjukkan folder satu per satu:**
- `Config/` - "Konfigurasi database, JWT, dan environment"
- `model/` - "Data models untuk Users, Achievements, Notifications"
- `repository/` - "Data access layer untuk database operations"
- `service/` - "Business logic layer yang handle HTTP requests"
- `middleware/` - "JWT authentication dan RBAC authorization"
- `Routes/` - "API routing dan endpoint definitions"
- `tests/` - "Comprehensive unit testing dengan mocks"
- `docs/` - "Swagger API documentation"

### **Demo 1.2: Key Files**
**Yang ditampilkan:** Buka `main.go`
**Yang dijelaskan:**
> "Di main.go, kita setup Fiber app, database connections, dependency injection, dan routing. Ini adalah entry point aplikasi."

---

## 🔐 BAGIAN 2: AUTHENTICATION & SECURITY (3 menit)

### **Demo 2.1: JWT Authentication**
**Yang ditampilkan:** Buka `service/authService.go`
**Yang dijelaskan:**
> "Sistem menggunakan JWT untuk authentication. Di sini kita lihat login function yang validate credentials, generate JWT token dengan role dan permissions."

**Scroll ke function Login:**
> "Login bisa menggunakan username atau email, password di-hash dengan bcrypt, dan return JWT token beserta user profile."

### **Demo 2.2: RBAC Middleware**
**Yang ditampilkan:** Buka `middleware/PermissionMiddleware.go`
**Yang dijelaskan:**
> "Untuk authorization, saya implementasi RBAC middleware. Setiap endpoint bisa di-protect dengan permission atau role tertentu."

**Tunjukkan function RequirePermission:**
> "Middleware ini extract JWT token, validate, dan check apakah user punya permission yang dibutuhkan."

---

## 📊 BAGIAN 3: CORE FEATURES (5 menit)

### **Demo 3.1: Achievement Management**
**Yang ditampilkan:** Buka `service/achievementService.go`
**Yang dijelaskan:**
> "Ini adalah core feature untuk management prestasi. Mahasiswa bisa submit prestasi dengan status draft."

**Tunjukkan function SubmitAchievement:**
> "Data disimpan di dual database - reference di PostgreSQL, detail di MongoDB untuk fleksibilitas."

### **Demo 3.2: Verification Workflow**
**Yang ditampilkan:** Buka `service/lecturerService.go`
**Yang dijelaskan:**
> "Dosen wali bisa verify atau reject prestasi mahasiswa bimbingannya."

**Tunjukkan VerifyAchievement dan RejectAchievement:**
> "Setiap action update status dan kirim notification ke mahasiswa."

### **Demo 3.3: Admin Management**
**Yang ditampilkan:** Buka `service/adminService.go`
**Yang dijelaskan:**
> "Admin punya full control - manage users, lihat semua prestasi, assign dosen wali."

### **Demo 3.4: Statistics & Reporting**
**Yang ditampilkan:** Buka `service/statisticsService.go`
**Yang dijelaskan:**
> "Sistem generate statistics real-time - total by type, period, top students, competition levels. Multi-role access dengan data sesuai permission."

---

## 🛠️ BAGIAN 4: API ENDPOINTS (2 menit)

### **Demo 4.1: API Routes**
**Yang ditampilkan:** Buka `Routes/Router.go`
**Yang dijelaskan:**
> "Semua API endpoints menggunakan prefix /api/v1/ sesuai best practices. Ada grouping berdasarkan functionality."

**Scroll melalui routes:**
- "Authentication routes untuk login, refresh, logout"
- "Protected routes dengan JWT middleware"
- "Role-based routes dengan permission checking"

### **Demo 4.2: Swagger Documentation**
**Yang ditampilkan:** Buka browser ke `http://localhost:8080/swagger/index.html`
**Yang dijelaskan:**
> "Saya buat interactive Swagger documentation lengkap dengan request/response examples, authentication, dan error handling."

**Demo beberapa endpoints di Swagger UI:**
> "Developers bisa test API langsung dari documentation ini."

---

## 🧪 BAGIAN 5: TESTING & QUALITY (2 menit)

### **Demo 5.1: Unit Testing Structure**
**Yang ditampilkan:** Buka folder `tests/`
**Yang dijelaskan:**
> "Saya implementasi comprehensive unit testing dengan mock strategy untuk isolate dependencies."

**Tunjukkan struktur:**
- `mocks/` - "Mock implementations untuk database dan external services"
- `unit/` - "Unit tests untuk service dan repository layers"
- `fixtures/` - "Test data yang reusable"
- `helpers/` - "Test utilities dan helpers"

### **Demo 5.2: Mock Implementation**
**Yang ditampilkan:** Buka `tests/mocks/mock_auth_repo.go`
**Yang dijelaskan:**
> "Semua external dependencies di-mock untuk testing yang reliable dan fast. No database connection needed untuk unit tests."

### **Demo 5.3: Test Runner**
**Yang ditampilkan:** Buka `Makefile` atau `tests/test_runner.go`
**Yang dijelaskan:**
> "Ada test runner untuk menjalankan semua test suites dengan coverage reporting."

---

## 🚀 BAGIAN 6: DEPLOYMENT & DEVOPS (2 menit)

### **Demo 6.1: Docker Configuration**
**Yang ditampilkan:** Buka `docker-compose.yml` dan `Dockerfile`
**Yang dijelaskan:**
> "Project sudah containerized dengan Docker. Multi-stage build untuk optimized production image."

### **Demo 6.2: GitHub CI/CD**
**Yang ditampilkan:** Buka `.github/workflows/ci.yml`
**Yang dijelaskan:**
> "Ada automated CI/CD pipeline dengan GitHub Actions - testing, security scanning, building, dan deployment."

### **Demo 6.3: Environment Configuration**
**Yang ditampilkan:** Buka `.env.example`
**Yang dijelaskan:**
> "Environment variables untuk different deployment stages - development, staging, production."

---

## 🎯 BAGIAN 7: DEMO LIVE API (3 menit)

### **Demo 7.1: Start Application**
**Yang ditampilkan:** Terminal dengan command `go run main.go`
**Yang dijelaskan:**
> "Mari kita jalankan aplikasi dan demo beberapa endpoints."

### **Demo 7.2: Test Authentication**
**Yang ditampilkan:** Postman atau curl untuk login
**Yang dijelaskan:**
> "Test login endpoint - input username/password, dapat JWT token."

### **Demo 7.3: Test Protected Endpoints**
**Yang ditampilkan:** Test beberapa protected endpoints dengan JWT
**Yang dijelaskan:**
> "Test endpoints yang require authentication dan specific permissions."

### **Demo 7.4: Database Integration**
**Yang ditampilkan:** Database client showing data
**Yang dijelaskan:**
> "Data tersimpan di PostgreSQL dan MongoDB sesuai dengan dual database architecture."

---

## 🏆 PENUTUP (1 menit)

### **Summary & Achievements**
**Yang ditampilkan:** Kembali ke `PENJELASAN_PROJECT.md`
**Yang dijelaskan:**
> "Untuk summary, project ini berhasil mengimplementasi:"

**Highlight key points:**
- "✅ Semua functional requirements FR-001 sampai FR-011"
- "✅ Security dengan JWT dan RBAC"
- "✅ Scalable architecture dengan clean code"
- "✅ Comprehensive testing dan documentation"
- "✅ Production-ready dengan CI/CD pipeline"
- "✅ Modern tech stack dengan best practices"

> "Sistem ini siap untuk production deployment dan bisa di-scale sesuai kebutuhan institusi."

### **Q&A Preparation**
**Yang dijelaskan:**
> "Terima kasih atas perhatiannya. Saya siap menjawab pertanyaan yang mungkin ada."

---

## 📝 TIPS PRESENTASI VIDEO

### **Persiapan Sebelum Recording:**
1. ✅ **Close aplikasi tidak perlu** (browser tabs, notifications)
2. ✅ **Prepare terminal windows** dengan commands ready
3. ✅ **Setup database** dan pastikan data test tersedia
4. ✅ **Test run aplikasi** sekali untuk memastikan working
5. ✅ **Prepare Postman collection** untuk API testing
6. ✅ **Set screen resolution** yang optimal untuk recording

### **Selama Recording:**
1. 🎤 **Speak clearly** dan tidak terlalu cepat
2. 👆 **Point cursor** ke bagian yang sedang dijelaskan
3. ⏸️ **Pause sejenak** setelah setiap section
4. 🔍 **Zoom in** kalau perlu untuk readability
5. 📱 **Minimize distractions** (notifications, sounds)

### **Technical Checklist:**
- [ ] Audio quality bagus (test microphone)
- [ ] Screen resolution optimal (1920x1080 recommended)
- [ ] Font size cukup besar untuk dibaca
- [ ] Internet connection stable untuk live demo
- [ ] Backup plan kalau ada technical issues

### **Content Flow:**
1. **Start with big picture** → zoom into details
2. **Show code** → explain functionality → demo result
3. **Connect features** to business requirements
4. **Emphasize technical achievements** dan best practices
5. **End with impact** dan production readiness

---

## 🎬 SCRIPT BACKUP (Jika Ada Masalah Technical)

**Jika aplikasi tidak bisa dijalankan:**
> "Untuk demo live, saya akan tunjukkan melalui Swagger documentation dan explain expected behavior berdasarkan code implementation."

**Jika database connection error:**
> "Saya akan explain database schema dan tunjukkan mock data di test fixtures untuk demonstrate data flow."

**Jika ada error lain:**
> "Mari kita fokus ke code implementation dan architecture explanation, yang menunjukkan understanding terhadap concepts dan best practices."

---

**Good luck dengan presentasi! 🚀**