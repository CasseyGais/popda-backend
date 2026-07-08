# POPDA 2026 — Project Overview

> **Sistem Informasi Pekan Olahraga Pelajar Daerah (POPDA) Provinsi Banten Tahun 2026**  
> Mengelola proses pendaftaran kontingen, atlet, pelatih, official, serta pencatatan pertandingan.

---

## Daftar Isi

1. [Teknologi Stack](#1-teknologi-stack)
2. [Struktur Proyek Backend](#2-struktur-proyek-backend)
3. [Database](#3-database)
4. [Autentikasi & Otorisasi](#4-autentikasi--otorisasi)
5. [Alur Pendaftaran (Tahap 1-2-3)](#5-alur-pendaftaran-tahap-1-2-3)
6. [Daftar Endpoint API](#6-daftar-endpoint-api)
7. [Dokumentasi Per Fitur](#7-dokumentasi-per-fitur)

---

## 1. Teknologi Stack

### Backend

| Komponen | Teknologi | Versi |
|---|---|---|
| Bahasa | Go | 1.24.0 |
| Web Framework | Gin | v1.11.0 |
| ORM | GORM | v1.31.1 |
| Database Driver | gorm.io/driver/mysql | v1.6.0 |
| JWT | golang-jwt/jwt/v5 | v5.3.0 |
| CORS | gin-contrib/cors | v1.7.6 |
| PDF Generator | jung-kurt/gofpdf | v1.16.2 |
| Excel Generator | xuri/excelize/v2 | v2.8.1 |
| Env Config | joho/godotenv | v1.5.1 |
| Password Hash | SHA-256 (crypto/sha256) | stdlib |

### Database

| Komponen | Teknologi |
|---|---|
| DBMS | MariaDB 10.4.32 |
| Nama Database | `popda_2026` |
| Charset | `utf8mb4_general_ci` |

### Backend Server

| Setting | Nilai |
|---|---|
| Port default | `8000` |
| Mode | `gin.ReleaseMode` |
| Base URL | `http://localhost:8000` |
| Static files | `/avatar/`, `/uploads/` |

---

## 2. Struktur Proyek Backend

```
backend/
├── cmd/api/main.go               ← Entry point, DI & routing
├── go.mod / go.sum               ← Go module dependencies
├── .env                          ← DB_DSN, PORT, JWT_SECRET
├── pkg/
│   └── jwt/jwt.go                ← Token generate & parse
├── internal/
│   ├── auth/                     ← Login, logout
│   ├── users/                    ← Manajemen user
│   ├── roles/                    ← Manajemen role
│   ├── permissions/              ← Manajemen permission
│   ├── rolepermissions/          ← Relasi role ↔ permission
│   ├── modules/                  ← Manajemen modul sistem
│   ├── territories/              ← Wilayah (provinsi/kab/kota)
│   ├── kontingen/                ← Data kontingen
│   ├── kontingenidentitas/       ← Identitas kepala kontingen & PIC
│   ├── tahap1/                   ← Entry By Sport (cabor + kuota)
│   ├── tahap2/                   ← Entry By Number (nomor pertandingan)
│   ├── tahap3/                   ← Entry By Name (atlet, pelatih, official)
│   ├── masterpelatih/            ← Master data pelatih
│   ├── masterofficial/           ← Master data official
│   ├── atlet/                    ← Master data atlet (legacy)
│   ├── cabor/                    ← Master cabang olahraga
│   ├── nomor/                    ← Master nomor pertandingan
│   ├── sekolah/                  ← Master data sekolah
│   ├── transaksi/                ← Transaksi legacy
│   ├── pengaturantahap/          ← Buka/tutup tahap pendaftaran
│   ├── validasipendaftaran/      ← Validasi data kontingen oleh superadmin
│   ├── sertifikat/               ← Penerbitan sertifikat
│   ├── laporanpertandingan/      ← Laporan hasil pertandingan
│   └── shared/
│       ├── database/database.go  ← Koneksi DB
│       └── middleware/
│           ├── auth.go           ← AuthRequired, PermissionRequired, SuperadminOnly, RolesAllowed
│           └── tahap.go          ← TahapOpen (cek is_open tahap)
└── uploads/                      ← File upload (foto, dokumen, laporan)
    ├── kepala/
    ├── pic/
    ├── atlet/
    ├── pelatih/
    ├── official/
    ├── sertifikat/
    └── laporan/
```

### Pola Arsitektur Per Modul

Setiap modul mengikuti pola 4 file konsisten:

```
modul/
├── entity.go       ← Struct DB mapping (GORM) + request/response structs
├── repository.go   ← Query database (GORM)
├── service.go      ← Business logic & validasi
└── handler.go      ← HTTP handler (Gin)
```

---

## 3. Database

**Database:** `popda_2026` — MariaDB 10.4.32

### Tabel Utama

| Tabel | Keterangan |
|---|---|
| `users` | User login (id: int) |
| `roles` | Role: SUPERADMIN, ADMIN, STAFF_LAPANGAN |
| `user_roles` | Relasi user ↔ role |
| `permissions` | Daftar permission per modul |
| `role_permissions` | Relasi role ↔ permission |
| `modules` | Daftar modul sistem |
| `territories` | Wilayah: PROVINSI / KABUPATEN / KOTA |
| `user_territories` | Relasi user ↔ territory |
| `kontingen` | Data kontingen per territory + status tahap + validasi |
| `kontingen_identitas` | Kepala kontingen & PIC |
| `master_cabor` | Cabang olahraga + kuota atlet |
| `master_nomor` | Nomor/kelas pertandingan per cabor |
| `master_atlet` | Data atlet |
| `master_pelatih` | Data pelatih |
| `master_official` | Data official |
| `master_sekolah` | Data sekolah (opsional) |
| `trx_kontingen_cabor` | Tahap 1: cabor yang dipilih kontingen |
| `trx_kontingen_nomor` | Tahap 2: nomor yang dipilih kontingen |
| `trx_pendaftaran_atlet` | Tahap 3: pendaftaran atlet ke nomor |
| `trx_pendaftaran_pelatih` | Tahap 3: pendaftaran pelatih ke cabor |
| `trx_pendaftaran_official` | Tahap 3: pendaftaran official |
| `pengaturan_tahap` | Buka/tutup tahap 1/2/3 secara global |
| `sertifikat` | Sertifikat untuk atlet/pelatih/official |
| `laporan_pertandingan` | Laporan hasil pertandingan |
| `laporan_pertandingan_atlet` | Atlet per sisi per pertandingan |

### Kolom Status di `kontingen`

```
tahap1_status              ENUM('DRAFT','SUBMITTED')
tahap1_submitted_at        DATETIME
tahap1_validasi_status     ENUM('PENDING','VALID','REVISI')
tahap1_validasi_catatan    TEXT
tahap1_validasi_at         DATETIME
(dan kolom yang sama untuk tahap2 dan tahap3)
```

---

## 4. Autentikasi & Otorisasi

### Login & JWT

```
POST /login
Body: { "email": "...", "password": "..." }
Response: { "token": "Bearer eyJ...", "user": {...}, "role": "ADMIN" }
```

Password di-hash dengan **SHA-256** sebelum disimpan. JWT berlaku **24 jam**.

### JWT Claims

```go
type Claims struct {
    UserID      uint    // ID user
    KontingenID uint    // 0 jika superadmin/staff
    Email       string
    Role        string  // "SUPERADMIN" / "ADMIN" / "STAFF_LAPANGAN"
    Name        string
    KabKota     string
    FotoProfil  string
}
```

### Role & Akses

| Role | Role ID | Akses |
|---|---|---|
| `SUPERADMIN` | 1 | Full access + bypass semua permission + akses data semua kontingen |
| `ADMIN` | 2 | Akses terbatas + hanya data kontingen sendiri |
| `STAFF_LAPANGAN` | 3 | Akses sertifikat + laporan pertandingan |

### Deteksi Superadmin (3 Lapis)

Backend mendeteksi superadmin secara **case-insensitive** via 3 kondisi:
1. `claims.Role == "superadmin"` (cara utama)
2. `claims.KontingenID == 0` (fallback token lama)
3. `?territory_id` ada di query → selalu override JWT

### Middleware yang Tersedia

| Middleware | Fungsi |
|---|---|
| `AuthRequired()` | Validasi JWT, wajib di semua `/admin/*` |
| `PermissionRequired(db, "cabor.read")` | Cek permission spesifik, superadmin bypass |
| `SuperadminOnly()` | Hanya superadmin |
| `RolesAllowed("SUPERADMIN", "STAFF_LAPANGAN")` | Daftar role yang diizinkan |
| `TahapOpen(db, tahap)` | Blokir operasi tulis jika tahap belum dibuka |

---

## 5. Alur Pendaftaran (Tahap 1-2-3)

### Overview

```
Superadmin buka Tahap 1
        ↓
Kontingen daftar Cabor + kuota atlet  (Tahap 1 — Entry By Sport)
        ↓
Kontingen submit → validasi_status = PENDING
        ↓
Superadmin buka Tahap 2
        ↓
Kontingen pilih Nomor/Kelas dari cabor yang dipilih  (Tahap 2 — Entry By Number)
        ↓
Kontingen submit → validasi_status = PENDING
        ↓
Superadmin buka Tahap 3
        ↓
Kontingen input data Atlet, Pelatih, Official  (Tahap 3 — Entry By Name)
        ↓
Kontingen submit → validasi_status = PENDING
        ↓
Superadmin review → set VALID atau REVISI (dengan catatan)
        ↓
Jika REVISI → kontingen perbaiki + submit ulang → PENDING lagi
```

### Aturan Buka Tahap

- Tahap 2 tidak bisa dibuka sebelum Tahap 1 pernah dibuka
- Tahap 3 tidak bisa dibuka sebelum Tahap 2 pernah dibuka
- Superadmin **bypass** semua gate tahap — bisa akses meski tahap tutup

### Dependency Antar Tahap

| Tahap | Syarat Akses |
|---|---|
| Tahap 1 | Tahap 1 is_open = true |
| Tahap 2 | Tahap 1 SUBMITTED + Tahap 2 is_open = true |
| Tahap 3 | Tahap 1 SUBMITTED + Tahap 2 SUBMITTED (opsional) + Tahap 3 is_open = true |

### Superadmin — Akses via Territory

Superadmin tidak punya `kontingen_id` di JWT (nilainya 0). Untuk akses data kontingen tertentu, wajib kirim `?territory_id=X`:

```
GET /admin/tahap1?territory_id=2   ← Superadmin
GET /admin/tahap1                  ← Admin biasa (kontingen_id dari JWT)
```

---

## 6. Daftar Endpoint API

**Base URL:** `http://localhost:8000`  
**Auth:** Header `Authorization: Bearer <token>` wajib di semua `/admin/*`

### Public

| Method | URL | Keterangan |
|---|---|---|
| `POST` | `/login` | Login, return JWT token |
| `POST` | `/logout` | Logout (stateless, hapus token di client) |
| `GET` | `/territories` | Semua wilayah (tanpa auth) |

---

### Kontingen & Identitas

| Method | URL | Keterangan |
|---|---|---|
| `GET` | `/admin/identitas` | Identitas kontingen sendiri |
| `PUT` | `/admin/identitas` | Update identitas |
| `GET` | `/admin/kontingen/:id` | Detail kontingen |
| `GET` | `/admin/kontingen/territory/:territory_id` | Kontingen by territory |
| `POST` | `/admin/kontingen` | Buat kontingen baru |
| `PUT` | `/admin/kontingen/:id` | Update kontingen |
| `GET/POST/PUT/DELETE` | `/admin/kontingen-identitas/*` | CRUD identitas + upload foto kepala/PIC |

---

### Tahap 1 — Entry By Sport

| Method | URL | Keterangan |
|---|---|---|
| `GET` | `/admin/tahap1` | Status + daftar cabor. Superadmin: `?territory_id=X` |
| `PUT` | `/admin/tahap1` | Upsert satu cabor (form-data). Dicek TahapOpen |
| `DELETE` | `/admin/tahap1/:cabor_id` | Hapus cabor. Dicek TahapOpen |
| `POST` | `/admin/tahap1/submit` | Kunci Tahap 1. Dicek TahapOpen |
| `GET` | `/admin/tahap1/export/pdf` | Download PDF rekap cabor |
| `GET` | `/admin/tahap1/export/excel` | Download Excel rekap cabor |

---

### Tahap 2 — Entry By Number

| Method | URL | Keterangan |
|---|---|---|
| `GET` | `/admin/tahap2` | Status + daftar nomor tersedia. Superadmin: `?territory_id=X` |
| `POST` | `/admin/tahap2/nomor/:nomor_id` | Daftar satu nomor. Dicek TahapOpen |
| `DELETE` | `/admin/tahap2/nomor/:nomor_id` | Batal daftar nomor. Dicek TahapOpen |
| `POST` | `/admin/tahap2/submit` | Kunci Tahap 2. Dicek TahapOpen |
| `GET` | `/admin/tahap2/export/pdf` | Download PDF rekap nomor |
| `GET` | `/admin/tahap2/export/excel` | Download Excel rekap nomor |

---

### Tahap 3 — Entry By Name

| Method | URL | Keterangan |
|---|---|---|
| `GET` | `/admin/tahap3` | Semua data tahap 3 sekaligus |
| `GET` | `/admin/tahap3/cabor` | Cabor terpilih dari tahap 1 |
| `GET` | `/admin/tahap3/nomor` | Nomor terdaftar dari tahap 2 |
| `POST` | `/admin/tahap3/submit` | Kunci Tahap 3. Dicek TahapOpen |
| `GET/POST/PUT/DELETE` | `/admin/tahap3/atlet/*` | CRUD atlet + upload foto/file |
| `GET` | `/admin/tahap3/atlet/export/pdf` | Export PDF atlet |
| `GET` | `/admin/tahap3/atlet/export/excel` | Export Excel atlet |
| `GET/POST/PUT/DELETE` | `/admin/tahap3/pelatih/*` | CRUD pelatih + upload file |
| `GET/POST/PUT/DELETE` | `/admin/tahap3/official/*` | CRUD official + upload file |
| `POST/DELETE` | `/admin/tahap3/trx/atlet` | Transaksi pendaftaran atlet |
| `POST/DELETE` | `/admin/tahap3/trx/pelatih` | Transaksi pendaftaran pelatih |
| `POST/DELETE` | `/admin/tahap3/trx/official` | Transaksi pendaftaran official |
| `GET` | `/admin/tahap3/export/pdf` | PDF semua (atlet+pelatih+official) |
| `GET` | `/admin/tahap3/export/excel` | Excel 3 sheet |

---

### Master Data

| Method | URL | Keterangan |
|---|---|---|
| `GET/POST/PUT/DELETE` | `/admin/master/cabor` | CRUD cabang olahraga |
| `GET/POST/PUT/DELETE` | `/admin/master/nomor` | CRUD nomor pertandingan |
| `GET/POST/PUT/DELETE` | `/admin/master/sekolah` | CRUD sekolah |
| `GET/POST/PUT/DELETE` | `/admin/master/pelatih` | CRUD pelatih + upload foto/file + trx |
| `GET/POST/PUT/DELETE` | `/admin/master/official` | CRUD official + upload foto/file + trx |

---

### Administrasi (Superadmin)

| Method | URL | Keterangan |
|---|---|---|
| `GET/POST/PUT/DELETE` | `/admin/users/*` | CRUD user + assign role/territory |
| `GET/POST/PUT/DELETE` | `/admin/roles/*` | CRUD role + assign permission |
| `GET/POST/PUT/DELETE` | `/admin/permissions/*` | CRUD permission |
| `GET/POST/PUT/DELETE` | `/admin/modules/*` | CRUD modul sistem |
| `GET/POST/PUT/DELETE` | `/admin/territories/*` | CRUD wilayah + assign ke user |
| `GET/POST/PUT/DELETE` | `/admin/role-permissions/*` | Relasi role ↔ permission |

---

### Pengaturan Tahap

| Method | URL | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/pengaturan-tahap` | Semua | Status buka/tutup tahap 1/2/3 |
| `PUT` | `/admin/pengaturan-tahap/:tahap` | Superadmin | Toggle buka/tutup + set tanggal |

---

### Validasi Pendaftaran

| Method | URL | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/validasi-pendaftaran/status` | Semua | Status validasi kontingen sendiri (widget) |
| `GET` | `/admin/validasi-pendaftaran` | Superadmin | List semua kontingen + status validasi |
| `PUT` | `/admin/validasi-pendaftaran/:id/tahap/:tahap` | Superadmin | Set VALID atau REVISI |
| `GET` | `/admin/rekap-pendaftaran` | Semua | Semua data kontingen dalam satu response |

---

### Sertifikat

| Method | URL | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/sertifikat` | SUPERADMIN, STAFF_LAPANGAN | List sertifikat |
| `POST` | `/admin/sertifikat` | SUPERADMIN, STAFF_LAPANGAN | Buat sertifikat |
| `GET/PUT/DELETE` | `/admin/sertifikat/:id` | SUPERADMIN, STAFF_LAPANGAN | CRUD sertifikat |
| `PUT` | `/admin/sertifikat/:id/file` | SUPERADMIN, STAFF_LAPANGAN | Upload file PDF manual |
| `POST` | `/admin/sertifikat/:id/export/pdf` | SUPERADMIN, STAFF_LAPANGAN | Export PDF landscape + tanda tangan (body JSON) |
| `POST` | `/admin/sertifikat/export/batch/pdf` | SUPERADMIN, STAFF_LAPANGAN | Export PDF batch + tanda tangan (body JSON) |
| `GET` | `/admin/sertifikat/penerima/atlet` | SUPERADMIN, STAFF_LAPANGAN | Dropdown atlet |
| `GET` | `/admin/sertifikat/penerima/pelatih` | SUPERADMIN, STAFF_LAPANGAN | Dropdown pelatih |
| `GET` | `/admin/sertifikat/penerima/official` | SUPERADMIN, STAFF_LAPANGAN | Dropdown official |

---

### Laporan Pertandingan

| Method | URL | Keterangan |
|---|---|---|
| `GET` | `/admin/laporan-pertandingan` | List + filter (tanggal, cabor, babak, pemenang) |
| `POST` | `/admin/laporan-pertandingan` | Buat laporan baru |
| `POST` | `/admin/laporan-pertandingan/export/pdf` | Export PDF batch + tanda tangan |
| `GET` | `/admin/laporan-pertandingan/dropdown/kontingen` | Dropdown Tim A/B |
| `GET` | `/admin/laporan-pertandingan/dropdown/cabor` | Dropdown cabor |
| `GET` | `/admin/laporan-pertandingan/dropdown/nomor` | Dropdown nomor (`?cabor_id`) |
| `GET` | `/admin/laporan-pertandingan/dropdown/atlet` | Dropdown atlet terdaftar |
| `GET/PUT/DELETE` | `/admin/laporan-pertandingan/:id` | CRUD laporan |
| `PUT` | `/admin/laporan-pertandingan/:id/foto` | Upload foto bukti |
| `PUT` | `/admin/laporan-pertandingan/:id/video` | Upload video bukti |
| `POST` | `/admin/laporan-pertandingan/:id/export/pdf` | Export PDF satu laporan + tanda tangan |

---

## 7. Dokumentasi Per Fitur

File dokumentasi detail per fitur tersedia di folder `backend/`:

| File | Isi |
|---|---|
| `TAHAP1_DOCUMENTATION.md` | Entry By Sport — endpoint, aturan bisnis, alur frontend |
| `TAHAP2_DOCUMENTATION.md` | Entry By Number — endpoint, aturan bisnis, alur frontend |
| `TAHAP3_DOCUMENTATION.md` | Entry By Name — endpoint lengkap atlet/pelatih/official |
| `FEATURE_PENGATURAN_TAHAP.md` | Spec backend pengaturan buka/tutup tahap |
| `FEATURE_PENGATURAN_TAHAP_FRONTEND.md` | Dokumentasi frontend pengaturan tahap |
| `FEATURE_VALIDASI_PENDAFTARAN.md` | Spec backend validasi pendaftaran |
| `FEATURE_VALIDASI_PENDAFTARAN_FRONTEND.md` | Dokumentasi frontend validasi & rekap pendaftaran |
| `FEATURE_SERTIFIKAT_FRONTEND.md` | Dokumentasi frontend sertifikat |
| `FEATURE_LAPORAN_PERTANDINGAN_FRONTEND.md` | Dokumentasi frontend laporan pertandingan |
| `SUPERADMIN_PANEL.md` | Panduan superadmin panel (users, roles, territories, dll) |
| `BACKEND_BUG_TAHAP1_TERRITORY.md` | Dokumentasi bug territory superadmin yang sudah di-fix |

---

## 8. Konvensi Response API

Semua endpoint menggunakan format JSON yang konsisten:

### Success

```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": { ... }
}
```

### Error

```json
{
  "success": false,
  "message": "Pesan error yang deskriptif"
}
```

### Konvensi Khusus

| Hal | Konvensi |
|---|---|
| Password | Dikirim **plaintext** dari frontend, di-hash SHA-256 di backend |
| File upload | `multipart/form-data`, disimpan di `uploads/` |
| Tanggal `DATE` kolom | Response selalu `"YYYY-MM-DD"` (bukan ISO timestamp) |
| `created_by` | Diisi otomatis dari JWT, tidak perlu dikirim frontend |
| `kontingen_id` | Diisi otomatis dari JWT atau territory resolve |
| Partial update `PUT` | Hanya field yang dikirim yang diupdate |
| Superadmin `?territory_id` | Wajib di endpoint yang resolve kontingen |
| Enum `babak`/`pemenang` | Case-insensitive di backend (auto `ToUpper`) |

---

## 9. Konfigurasi (.env)

```env
DB_DSN=root:password@tcp(127.0.0.1:3306)/popda_2026?charset=utf8mb4&parseTime=True&loc=Local
PORT=8000
JWT_SECRET=rahasiasuperkuat1234567890
```

---

## 10. Menjalankan Server

```bash
# Development
go run cmd/api/main.go

# Build executable
go build -o api.exe cmd/api/main.go
./api.exe
```

Server berjalan di `http://localhost:8000`

---

## 11. User Akun Default (Seed Data)

| User | Email | Password | Role |
|---|---|---|---|
| Superadmin | `superadmin@popda.id` | `1` | SUPERADMIN |
| Admin Kab Tangerang | `admin.kabtangerang@popda.id` | `popda2026` | ADMIN |
| Admin Kab Serang | `admin.kabserang@popda.id` | `popda2026` | ADMIN |
| Admin Kab Lebak | `admin.kablebak@popda.id` | `popda2026` | ADMIN |
| Admin Kab Pandeglang | `admin.kabpandeglang@popda.id` | `admin123` | ADMIN |
| Admin Kota Tangerang | `admin.kotatangerang@popda.id` | `admin123` | ADMIN |
| Admin Tangsel | `admin.tangsel@popda.id` | `admin123` | ADMIN |
| Admin Kota Serang | `admin.kotaserang@popda.id` | `admin123` | ADMIN |
| Admin Kota Cilegon | `admin.kotacilegon@popda.id` | `admin123` | ADMIN |

> Password disimpan sebagai SHA-256 hash. Kirim plaintext ke `/login`.

---

## 12. Wilayah (Territories)

| ID | Nama | Tipe |
|---|---|---|
| 2 | Kabupaten Tangerang | KABUPATEN |
| 3 | Kabupaten Serang | KABUPATEN |
| 4 | Kabupaten Lebak | KABUPATEN |
| 5 | Kabupaten Pandeglang | KABUPATEN |
| 6 | Kota Tangerang | KOTA |
| 7 | Kota Tangerang Selatan | KOTA |
| 8 | Kota Serang | KOTA |
| 9 | Kota Cilegon | KOTA |

> Territory ID 1 tidak ada — dimulai dari 2 (sesuai seed data).
