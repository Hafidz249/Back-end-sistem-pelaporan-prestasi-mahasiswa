# 🎓 Back-end Sistem Pelaporan Prestasi Mahasiswa (UAS Semester 5)

![Status: Selesai (Final Project)](https://img.shields.io/badge/Status-Complete-brightgreen)
![Lisensi: MIT](https://img.shields.io/badge/License-MIT-blue)
![Teknologi Utama: Node.js | Express.js](https://img.shields.io/badge/Tech-Node.js%20%7C%20Express-red)

Repositori ini menyimpan kode sumber *back-end* (server-side) untuk **Sistem Pelaporan Prestasi Mahasiswa**. Sistem ini dirancang untuk memfasilitasi pencatatan, verifikasi, dan pelaporan prestasi akademik maupun non-akademik mahasiswa secara terstruktur.

## ✨ Fitur Utama (Fungsionalitas API)

Sistem *back-end* ini menyediakan *endpoint* RESTful API untuk mendukung fungsionalitas inti berikut:

1.  **Otentikasi & Autorasi (JWT):**
    * Registrasi dan *Login* pengguna dengan peran berbeda (Mahasiswa, Dosen/Verifikator, Admin).
    * Pengamanan *endpoint* menggunakan *JSON Web Tokens* (JWT) dan kontrol akses berbasis peran.
2.  **Manajemen Mahasiswa & Pengguna:**
    * CRUD Data Mahasiswa oleh Admin.
    * Pembaruan profil pengguna.
3.  **Pelaporan Prestasi:**
    * Mahasiswa dapat membuat, melihat, dan memperbarui laporan prestasi mereka.
    * Mendukung pengunggahan *file* bukti/sertifikat prestasi.
4.  **Verifikasi Prestasi:**
    * Dosen/Verifikator dapat meninjau, menyetujui, atau menolak laporan prestasi yang masuk.
5.  **Pencarian & Filtrasi:**
    * Pencarian data prestasi berdasarkan NIM, jenis prestasi, atau status verifikasi.

## 🛠️ Stack Teknologi

| Kategori | Teknologi | Deskripsi |
| :--- | :--- | :--- |
| **Bahasa** | **Node.js** | Lingkungan runtime untuk eksekusi kode server. |
| **Framework** | **Express.js** | Framework minimalis dan fleksibel untuk membangun API RESTful. |
| **Database** | **PostgreSQL** | Sistem manajemen basis data relasional yang handal. |
| **ORM/Query** | [Tambahkan ORM Anda, contoh: Sequelize atau Knex.js] | Alat untuk berinteraksi dengan database. |
| **Otentikasi** | **JWT** (JSON Web Tokens) | Standar untuk otentikasi yang aman. |

## ⚙️ Instalasi Lokal

Ikuti langkah-langkah di bawah ini untuk menjalankan proyek secara lokal di lingkungan pengembangan Anda.

### 1. Prasyarat

Pastikan Anda telah menginstal perangkat lunak berikut:

* **Node.js** (Versi 18+)
* **npm** (Node Package Manager)
* **PostgreSQL** (Server Database)

### 2. Kloning Repositori

```bash
git clone [https://github.com/Hafidz249/Back-end-sistem-pelaporan-prestasi-mahasiswa.git](https://github.com/Hafidz249/Back-end-sistem-pelaporan-prestasi-mahasiswa.git)
cd Back-end-sistem-pelaporan-prestasi-mahasiswa
