# POPDA Backend API Documentation

> **Base URL:** `http://localhost:8000`  
> **Tech Stack:** Go (Gin), MySQL (GORM), JWT Authentication  
> **Module:** `popda_bulutangkis`

---

## Daftar Isi

1. [Autentikasi](#1-autentikasi)
2. [Alur Pendaftaran (Tahap 1 → 3)](#2-alur-pendaftaran)
3. [Kontingen](#3-kontingen)
4. [Transaksi](#4-transaksi)
5. [Master Data — Cabor](#5-master-data--cabor)
6. [Master Data — Nomor](#6-master-data--nomor)
7. [Master Data — Sekolah](#7-master-data--sekolah)
8. [Atlet](#8-atlet)
9. [Kontingen Identitas](#9-kontingen-identitas)
10. [Master Pelatih](#10-master-pelatih)
11. [Master Official](#11-master-official)
12. [Users](#12-users)
13. [Roles](#13-roles)
14. [Permissions](#14-permissions)
15. [Role Permissions](#15-role-permissions)
16. [Territories](#16-territories)
17. [Static Files](#17-static-files)
18. [Model Referensi](#18-model-referensi)

---

## Konvensi Umum

### Request Headers (Protected Routes)
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Format Response Sukses
```json
{
  "success": true,
  "message": "...",
  "data": { ... }
}
```

### Format Response Error
```json
{
  "success": false,
  "message": "...",
  "error": "detail error"
}
```

> **Catatan:** Beberapa endpoint (terutama transaksi) menggunakan format `{ "data": ... }` atau `{ "message": "..." }` tanpa wrapper `success`.

---

## 1. Autentikasi

### POST `/login`
Login dan dapatkan JWT token.

**Request Body (JSON):**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Nama User",
      "email": "user@example.com",
      "avatar": "/avatar/user.png",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "territories": [
        { "id": 1, "name": "Kabupaten Tangerang", "type": "KABUPATEN" }
      ]
    },
    "role": "admin"
  }
}
```

**Response 401:**
```json
{
  "success": false,
  "message": "Login gagal",
  "error": "password salah"
}
```

---

### POST `/logout` 🔒
Logout user (stateless — hapus token di sisi client).

**Response 200:**
```json
{
  "success": true,
  "message": "Logout berhasil"
}
```

---

### Cara Pakai Token JWT

Setelah login, simpan `token` dan sertakan di setiap request protected:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Token berisi claims:
| Field | Keterangan |
|---|---|
| `user_id` | ID user yang login |
| `kontingen_id` | ID kontingen milik user |
| `role` | Role user (string) |
| `email` | Email user |
| `name` | Nama user |
| `territory_name` | Nama wilayah user |
| `avatar` | Path foto avatar |

---

## 2. Alur Pendaftaran

Pendaftaran kontingen dibagi menjadi 3 tahap berurutan. Semua endpoint tahap menggunakan `kontingen_id` dari JWT token secara otomatis.

```
Tahap 1 → Pilih Cabor & Jumlah Personel → Submit
Tahap 2 → Pilih Nomor Pertandingan → Submit
Tahap 3 → Lihat & Submit Data Atlet/Pelatih/Official
```

---

### Tahap 1 — Cabor & Personel

#### GET `/admin/tahap1` 🔒
Ambil data tahap 1 milik kontingen yang sedang login.

**Response 200:**
```json
{
  "data": {
    "cabor_list": [...],
    "jumlah_atlet": 10,
    "jumlah_pelatih": 2,
    "jumlah_official": 1,
    "tahap1_submitted": false,
    "submitted_at": null
  }
}
```

---

#### PUT `/admin/tahap1` 🔒
Simpan/update data tahap 1.

**Request Body (multipart/form-data):**
| Field | Tipe | Keterangan |
|---|---|---|
| `caborList[]` | array of uint | ID cabor yang dipilih |
| `jumlahAtlet` | int | Jumlah atlet |
| `jumlahPelatih` | int | Jumlah pelatih |
| `jumlahOfficial` | int | Jumlah official |

**Response 200:**
```json
{ "message": "Berhasil disimpan" }
```

---

#### POST `/admin/tahap1/submit` 🔒
Submit tahap 1 (tidak bisa diubah setelah submit).

**Response 200:**
```json
{ "message": "Tahap 1 berhasil disubmit" }
```

---

### Tahap 2 — Nomor Pertandingan

#### GET `/admin/tahap2` 🔒
Ambil daftar nomor pertandingan beserta status keikutsertaan.

**Response 200:**
```json
{
  "data": {
    "tahap2_submitted": false,
    "submitted_at": null,
    "available_cabor": ["Bulutangkis", "Renang"],
    "events": [
      {
        "event_id": 1,
        "cabor": "Bulutangkis",
        "nama_event": "Tunggal Putra",
        "jenis_kelamin": "PUTRA",
        "ikut": true
      }
    ]
  }
}
```

---

#### PUT `/admin/tahap2` 🔒
Update pilihan nomor pertandingan.

**Request Body (JSON):**
```json
{
  "nomor_ids": [1, 3, 5]
}
```

**Response 200:**
```json
{ "message": "Berhasil disimpan" }
```

---

#### POST `/admin/tahap2/submit` 🔒
Submit tahap 2.

**Response 200:**
```json
{ "message": "Tahap 2 berhasil disubmit" }
```

---

### Tahap 3 — Data Atlet, Pelatih & Official

#### GET `/admin/tahap3` 🔒
Ambil semua data atlet, pelatih, dan official milik kontingen.

**Response 200:**
```json
{
  "data": {
    "tahap3_submitted": false,
    "submitted_at": null,
    "atlets": [...],
    "pelatihs": [...],
    "officials": [...]
  }
}
```

---

#### POST `/admin/tahap3/submit` 🔒
Submit tahap 3.

**Response 200:**
```json
{ "message": "Tahap 3 berhasil disubmit" }
```

---

## 3. Kontingen

### GET `/admin/identitas` 🔒
Ambil identitas kontingen milik user yang login.

**Response 200:**
```json
{
  "success": true,
  "message": "Data identitas berhasil diambil",
  "data": {
    "id": 1,
    "kontingen_id": 1,
    "kepala_nama": "Drs. Budi Santoso",
    "kepala_jabatan": "Kepala Dinas",
    "kepala_nip": "197001011990031001",
    "kepala_telepon": "08123456789",
    "kepala_foto": "/uploads/kepala/foto.jpg",
    "pic_nama": "Siti Rahayu",
    "pic_jabatan": "Koordinator",
    "pic_telepon": "08987654321",
    "pic_foto": "/uploads/pic/foto.jpg",
    "alamat": "Jl. Contoh No. 1",
    "email_instansi": "dinas@example.com",
    "phone_instansi": "02112345678",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### PUT `/admin/identitas` 🔒
Update identitas kontingen. Mendukung upload foto.

**Request Body (multipart/form-data):**
| Field | Tipe | Keterangan |
|---|---|---|
| `alamat` | string | Alamat instansi |
| `email_instansi` | string | Email instansi |
| `phone_instansi` | string | Telepon instansi |
| `kepala_nama` | string | Nama kepala |
| `kepala_jabatan` | string | Jabatan kepala |
| `kepala_nip` | string | NIP kepala |
| `kepala_telepon` | string | Telepon kepala |
| `kepala_foto` | file | Foto kepala (opsional) |
| `pic_nama` | string | Nama PIC |
| `pic_jabatan` | string | Jabatan PIC |
| `pic_telepon` | string | Telepon PIC |
| `pic_foto` | file | Foto PIC (opsional) |

---

### GET `/admin/kontingen/:id` 🔒
Ambil data kontingen berdasarkan ID.

### GET `/admin/kontingen/territory/:territory_id` 🔒
Ambil kontingen berdasarkan territory ID.

### POST `/admin/kontingen` 🔒
Buat kontingen baru.

**Request Body (JSON):**
```json
{
  "territory_id": 1,
  "nama_kontingen": "Kontingen Tangerang"
}
```

### PUT `/admin/kontingen/:id` 🔒
Update data kontingen.

---

## 4. Transaksi

> Semua endpoint transaksi menggunakan `kontingen_id` dari JWT token secara otomatis.  
> Request menggunakan **multipart/form-data** (bukan JSON).

### Transaksi Cabor

#### POST `/admin/trx/cabor` 🔒
Daftarkan cabor untuk kontingen.

**Request Body (form-data):**
| Field | Tipe | Keterangan |
|---|---|---|
| `cabor_id` | uint | ID cabor |
| `putra` | int | Jumlah atlet putra |
| `putri` | int | Jumlah atlet putri |
| `pelatih` | int | Jumlah pelatih |

#### GET `/admin/trx/cabor` 🔒
Ambil daftar cabor yang didaftarkan kontingen.

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "kontingen_id": 1,
      "cabor_id": 2,
      "putra": 3,
      "putri": 2,
      "pelatih": 1,
      "total_atlet": 5,
      "total_personel": 6
    }
  ]
}
```

#### PUT `/admin/trx/cabor` 🔒
Update data cabor kontingen (form-data sama seperti POST).

---

### Transaksi Nomor

#### POST `/admin/trx/nomor` 🔒
Tambah nomor pertandingan.

**Request Body (form-data):**
| Field | Tipe | Keterangan |
|---|---|---|
| `nomor_id` | uint | ID nomor pertandingan |

#### GET `/admin/trx/nomor` 🔒
Ambil daftar nomor yang diikuti kontingen.

#### DELETE `/admin/trx/nomor/:nomor_id` 🔒
Hapus nomor pertandingan dari kontingen.

---

### Transaksi Pendaftaran Atlet

#### POST `/admin/trx/atlet` 🔒
Daftarkan atlet ke nomor pertandingan.

**Request Body (form-data):**
| Field | Tipe | Keterangan |
|---|---|---|
| `atlet_id` | uint | ID atlet |
| `nomor_id` | uint | ID nomor pertandingan |
| `kelas_id` | uint | ID kelas (opsional) |

#### GET `/admin/trx/atlet` 🔒
Ambil daftar pendaftaran atlet milik kontingen.

#### PUT `/admin/trx/atlet/:atlet_id/:nomor_id` 🔒
Update status pendaftaran atlet.

**Request Body (form-data):**
| Field | Tipe | Keterangan |
|---|---|---|
| `status` | string | Status pendaftaran |

---

## 5. Master Data — Cabor

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/master/cabor` | Semua cabor |
| GET | `/admin/master/cabor/:id` | Cabor by ID |
| POST | `/admin/master/cabor` | Buat cabor baru |
| PUT | `/admin/master/cabor/:id` | Update cabor |
| DELETE | `/admin/master/cabor/:id` | Hapus cabor |

**Model Cabor:**
```json
{
  "id": 1,
  "nama": "Bulutangkis",
  "max_putra": 10,
  "max_putri": 10,
  "max_pelatih": 3,
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Request Body (POST/PUT JSON):**
```json
{
  "nama": "Bulutangkis",
  "max_putra": 10,
  "max_putri": 10,
  "max_pelatih": 3,
  "is_active": true
}
```

---

## 6. Master Data — Nomor

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/master/nomor` | Semua nomor |
| GET | `/admin/master/nomor/:id` | Nomor by ID |
| GET | `/admin/master/nomor/cabor/:cabor_id` | Nomor by Cabor ID |
| POST | `/admin/master/nomor` | Buat nomor baru |
| PUT | `/admin/master/nomor/:id` | Update nomor |
| DELETE | `/admin/master/nomor/:id` | Hapus nomor |

**Model Nomor:**
```json
{
  "ID": 1,
  "nama": "Tunggal Putra",
  "cabor_id": 1,
  "is_active": true,
  "Cabor": {
    "id": 1,
    "nama": "Bulutangkis"
  }
}
```

---

## 7. Master Data — Sekolah

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/master/sekolah` | Semua sekolah |
| GET | `/admin/master/sekolah/:id` | Sekolah by ID |
| GET | `/admin/master/sekolah/search` | Cari sekolah |
| POST | `/admin/master/sekolah` | Buat sekolah baru |
| PUT | `/admin/master/sekolah/:id` | Update sekolah |
| DELETE | `/admin/master/sekolah/:id` | Hapus sekolah |

**Model Sekolah:**
```json
{
  "ID": 1,
  "nama": "SMP Negeri 1 Tangerang",
  "npsn": "20603001",
  "alamat": "Jl. Contoh No. 1",
  "kabupaten": "Kabupaten Tangerang",
  "is_active": true
}
```

---

## 8. Atlet

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/atlet` | Semua atlet |
| GET | `/admin/atlet/:id` | Atlet by ID |
| GET | `/admin/atlet/kontingen/:kontingen_id` | Atlet by Kontingen |
| GET | `/admin/atlet/sekolah/:sekolah_id` | Atlet by Sekolah |
| GET | `/admin/atlet/status/:status` | Atlet by Status Verifikasi |
| POST | `/admin/atlet` | Buat atlet baru |
| PUT | `/admin/atlet/:id` | Update atlet |
| DELETE | `/admin/atlet/:id` | Hapus atlet |
| PUT | `/admin/atlet/:id/status` | Update status verifikasi |
| PUT | `/admin/atlet/:id/foto` | Update foto atlet |

**Model Atlet:**
```json
{
  "id": 1,
  "kontingen_id": 1,
  "sekolah_id": 2,
  "nisn": "1234567890",
  "nama": "Ahmad Fauzi",
  "jenis_kelamin": "PUTRA",
  "tanggal_lahir": "2008-05-15T00:00:00Z",
  "kelas": "VIII",
  "tinggi": 165,
  "berat": 55.5,
  "foto": "/uploads/atlet/foto.jpg",
  "status_verifikasi": "PENDING",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**POST `/admin/atlet` — Request Body (JSON):**
```json
{
  "kontingen_id": 1,
  "sekolah_id": 2,
  "nisn": "1234567890",
  "nama": "Ahmad Fauzi",
  "jenis_kelamin": "PUTRA",
  "tanggal_lahir": "2008-05-15T00:00:00Z",
  "kelas": "VIII",
  "tinggi": 165,
  "berat": 55.5,
  "foto": "/uploads/atlet/foto.jpg"
}
```

> `jenis_kelamin` harus salah satu dari: `PUTRA` | `PUTRI`

**PUT `/admin/atlet/:id/status` — Request Body (JSON):**
```json
{
  "status": "VALID"
}
```
> `status` harus salah satu dari: `PENDING` | `VALID` | `DITOLAK`

**PUT `/admin/atlet/:id/foto` — Request (multipart/form-data atau JSON):**
- Multipart: field `foto` berisi file gambar
- JSON fallback: `{ "foto": "/path/ke/foto.jpg" }`

---

## 9. Kontingen Identitas

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/kontingen-identitas` | Semua identitas |
| GET | `/admin/kontingen-identitas/:id` | By ID |
| GET | `/admin/kontingen-identitas/kontingen/:kontingen_id` | By Kontingen ID |
| POST | `/admin/kontingen-identitas` | Buat identitas baru |
| PUT | `/admin/kontingen-identitas/:id` | Update identitas |
| DELETE | `/admin/kontingen-identitas/:id` | Hapus identitas |
| PUT | `/admin/kontingen-identitas/:id/kepala-foto` | Update foto kepala |
| PUT | `/admin/kontingen-identitas/:id/pic-foto` | Update foto PIC |

> Setiap kontingen hanya boleh memiliki **satu** identitas (unique per `kontingen_id`).

---

## 10. Master Pelatih

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/master/pelatih` | Semua pelatih |
| GET | `/admin/master/pelatih/:id` | Pelatih by ID |
| GET | `/admin/master/pelatih/kontingen/:kontingen_id` | Pelatih by Kontingen |
| POST | `/admin/master/pelatih` | Buat pelatih baru |
| PUT | `/admin/master/pelatih/:id` | Update pelatih |
| DELETE | `/admin/master/pelatih/:id` | Hapus pelatih |
| PUT | `/admin/master/pelatih/:id/foto` | Update foto pelatih |
| PUT | `/admin/master/pelatih/:id/sertifikat` | Update sertifikat pelatih |

**Model Pelatih:**
```json
{
  "id": 1,
  "kontingen_id": 1,
  "nama": "Budi Pelatih",
  "no_hp": "08123456789",
  "sertifikat": "/uploads/sertifikat/file.pdf",
  "foto": "/uploads/pelatih/foto.jpg",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

## 11. Master Official

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/master/official` | Semua official |
| GET | `/admin/master/official/:id` | Official by ID |
| GET | `/admin/master/official/kontingen/:kontingen_id` | Official by Kontingen |
| POST | `/admin/master/official` | Buat official baru |
| PUT | `/admin/master/official/:id` | Update official |
| DELETE | `/admin/master/official/:id` | Hapus official |

**Model Official:**
```json
{
  "id": 1,
  "kontingen_id": 1,
  "nama": "Siti Official",
  "jabatan": "Manajer Tim",
  "no_hp": "08987654321",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

## 12. Users

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/users` | Semua user |
| GET | `/admin/users/:id` | User by ID |
| GET | `/admin/users/email/:email` | User by email |
| POST | `/admin/users` | Buat user baru |
| PUT | `/admin/users/:id` | Update user |
| DELETE | `/admin/users/:id` | Hapus user |
| PUT | `/admin/users/:id/avatar` | Update avatar |
| PUT | `/admin/users/:id/password` | Update password |
| PUT | `/admin/users/:id/status` | Update status aktif |
| GET | `/admin/users/:id/roles` | Ambil roles user |
| POST | `/admin/users/:id/roles/:role_id` | Assign role ke user |
| DELETE | `/admin/users/:id/roles/:role_id` | Hapus role dari user |
| GET | `/admin/users/:id/territories` | Ambil territories user |
| POST | `/admin/users/:id/territories/:territory_id` | Assign territory ke user |
| DELETE | `/admin/users/:id/territories/:territory_id` | Hapus territory dari user |

**Model User:**
```json
{
  "id": 1,
  "name": "Admin POPDA",
  "email": "admin@popda.id",
  "is_active": true,
  "avatar": "/avatar/user.png",
  "created_at": "2024-01-01T00:00:00Z"
}
```

> Password tidak pernah dikembalikan dalam response.

---

## 13. Roles

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/roles` | Semua role |
| GET | `/admin/roles/:id` | Role by ID |
| GET | `/admin/roles/user/:user_id` | Role by User ID |
| POST | `/admin/roles` | Buat role baru |
| PUT | `/admin/roles/:id` | Update role |
| DELETE | `/admin/roles/:id` | Hapus role |
| POST | `/admin/roles/:id/permissions/:permission_id` | Assign permission ke role |
| DELETE | `/admin/roles/:id/permissions/:permission_id` | Hapus permission dari role |
| GET | `/admin/roles/:id/permissions` | Ambil permissions role |

**Model Role:**
```json
{
  "id": 1,
  "name": "admin",
  "description": "Administrator sistem"
}
```

---

## 14. Permissions

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/permissions` | Semua permission |
| GET | `/admin/permissions/:id` | Permission by ID |
| GET | `/admin/permissions/role/:role_id` | Permission by Role ID |
| POST | `/admin/permissions` | Buat permission baru |
| PUT | `/admin/permissions/:id` | Update permission |
| DELETE | `/admin/permissions/:id` | Hapus permission |

---

## 15. Role Permissions

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/role-permissions` | Semua relasi role-permission |
| GET | `/admin/role-permissions/role/:id` | Permission by Role ID |
| GET | `/admin/role-permissions/permission/:id` | Role by Permission ID |
| POST | `/admin/role-permissions` | Assign permission ke role |
| DELETE | `/admin/role-permissions/role/:id/permission/:permissionId` | Hapus relasi spesifik |
| DELETE | `/admin/role-permissions/role/:id` | Hapus semua permission dari role |
| DELETE | `/admin/role-permissions/permission/:id` | Hapus permission dari semua role |

---

## 16. Territories

> Endpoint `GET /territories` tersedia **tanpa autentikasi** (public).  
> Endpoint di bawah `/admin/territories` memerlukan autentikasi.

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| GET | `/territories` | ❌ | Semua territory (public) |
| GET | `/admin/territories` | ✅ | Semua territory |
| GET | `/admin/territories/:id` | ✅ | Territory by ID |
| GET | `/admin/territories/type/:type` | ✅ | Territory by tipe |
| GET | `/admin/territories/provinces` | ✅ | Semua provinsi |
| GET | `/admin/territories/kabupatens` | ✅ | Semua kabupaten |
| GET | `/admin/territories/kotas` | ✅ | Semua kota |
| GET | `/admin/territories/user/:user_id` | ✅ | Territory by User ID |
| POST | `/admin/territories` | ✅ | Buat territory baru |
| PUT | `/admin/territories/:id` | ✅ | Update territory |
| DELETE | `/admin/territories/:id` | ✅ | Hapus territory |
| POST | `/admin/territories/user/:user_id/:territory_id` | ✅ | Assign territory ke user |
| DELETE | `/admin/territories/user/:user_id/:territory_id` | ✅ | Hapus territory dari user |

**Model Territory:**
```json
{
  "id": 1,
  "name": "Kabupaten Tangerang",
  "type": "KABUPATEN"
}
```

> `type` harus salah satu dari: `PROVINSI` | `KABUPATEN` | `KOTA`

---

## 17. Static Files

File statis dapat diakses langsung via URL:

| Path | Keterangan |
|---|---|
| `/avatar/{filename}` | Avatar user (folder `./avatar`) |
| `/uploads/{path}` | File upload (folder `./uploads`) |

Contoh path upload yang digunakan sistem:
- `/uploads/atlet/{filename}` — Foto atlet
- `/uploads/kepala/{filename}` — Foto kepala kontingen
- `/uploads/pic/{filename}` — Foto PIC kontingen
- `/uploads/pelatih/{filename}` — Foto pelatih
- `/uploads/sertifikat/{filename}` — Sertifikat pelatih

---

## 18. Model Referensi

### Status Verifikasi Atlet
| Nilai | Keterangan |
|---|---|
| `PENDING` | Menunggu verifikasi (default) |
| `VALID` | Sudah diverifikasi |
| `DITOLAK` | Ditolak |

### Status Tahap Kontingen
| Nilai | Keterangan |
|---|---|
| `DRAFT` | Belum disubmit (default) |
| `SUBMITTED` | Sudah disubmit |

### Jenis Kelamin
| Nilai | Keterangan |
|---|---|
| `PUTRA` | Laki-laki |
| `PUTRI` | Perempuan |

### Tipe Territory
| Nilai | Keterangan |
|---|---|
| `PROVINSI` | Tingkat provinsi |
| `KABUPATEN` | Tingkat kabupaten |
| `KOTA` | Tingkat kota |

---

## CORS

Backend mengizinkan request dari origin berikut:
- `http://localhost:5173` (Vite dev server)
- `http://localhost:8000`
- `http://127.0.0.1:8000`

Method yang diizinkan: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`  
Headers yang diizinkan: `Origin`, `Content-Type`, `Authorization`  
Credentials: **diizinkan** (`withCredentials: true`)

---

---

## 19. Changelog & Bug Fixes

### Fix #1 — Login Superadmin Gagal (NULL kontingen_id)

**File:** `internal/auth/repository.go` — fungsi `GetKontingenIDByUser`

**Root cause:**
Query berikut menggunakan `LEFT JOIN` ke tabel `kontingen`. Superadmin memiliki banyak territory tapi tidak semua territory punya kontingen, sehingga `k.id` bisa bernilai `NULL`:

```sql
SELECT k.id
FROM user_territories ut
JOIN territories t ON t.id = ut.territory_id
LEFT JOIN kontingen k ON k.territory_id = t.id
WHERE ut.user_id = 1
LIMIT 1
```

Go tidak bisa scan `NULL` ke tipe `uint`, sehingga muncul error:
```
sql: Scan error on column index 0, name "id": converting NULL to uint is unsupported
```
Login langsung return `401` meski password benar.

**Fix:**
Ganti variabel scan dari `uint` ke `*uint` (pointer). Pointer bisa menampung `nil` (NULL dari DB). Kalau `nil`, return `0` — login tetap jalan, `kontingen_id: 0` di JWT berarti user tidak terikat kontingen tertentu (sesuai untuk superadmin).

```go
// Sebelum
var kontingenID uint

// Sesudah
var kontingenID *uint  // pointer — bisa nil kalau tidak ada kontingen

// ...setelah query...
if kontingenID == nil {
    return 0, nil  // superadmin / user tanpa kontingen → lanjut login
}
return *kontingenID, nil
```

**Dampak:** Superadmin dengan `kontingen_id: 0` di JWT bisa login normal. Endpoint yang butuh kontingen spesifik (seperti `/admin/identitas`) tetap bekerja untuk admin biasa karena mereka punya `kontingen_id > 0`.

---

### Fix #2 — Upload Foto Superadmin di Kontingen Identitas

**File:** `internal/kontingenidentitas/handler.go` — fungsi `UpdateKepalaFoto` dan `UpdatePICFoto`

**Root cause:**
Endpoint foto kontingen identitas sebelumnya hanya menerima **JSON path string**:
```json
{ "foto": "/uploads/kepala/namafile.jpg" }
```
Tidak ada mekanisme upload file langsung. Superadmin yang ingin upload foto baru tidak bisa melakukannya karena tidak ada endpoint upload file untuk kontingen identitas (berbeda dengan `/admin/identitas` milik admin biasa yang sudah support multipart).

**Fix:**
Kedua handler diubah menjadi **dual-mode** — coba terima file upload dulu, kalau tidak ada fallback ke JSON:

```
Mode 1 (multipart/form-data):
  field "foto" berisi file gambar
  → disimpan ke uploads/kepala/ atau uploads/pic/
  → path disimpan ke DB
  → response mengembalikan path baru

Mode 2 (JSON fallback):
  { "foto": "/path/ke/foto.jpg" }
  → path langsung disimpan ke DB (behavior lama, tidak berubah)
```

**Cara pakai dari frontend (superadmin upload foto):**

```js
// Upload foto kepala
const formData = new FormData()
formData.append('foto', fileInput.files[0])

fetch(`/admin/kontingen-identitas/${id}/kepala-foto`, {
  method: 'PUT',
  headers: { Authorization: `Bearer ${token}` },
  body: formData  // jangan set Content-Type, biarkan browser set boundary
})

// Upload foto PIC
const formData = new FormData()
formData.append('foto', fileInput.files[0])

fetch(`/admin/kontingen-identitas/${id}/pic-foto`, {
  method: 'PUT',
  headers: { Authorization: `Bearer ${token}` },
  body: formData
})
```

**Response sukses (mode file upload):**
```json
{
  "success": true,
  "message": "Foto kepala berhasil diupdate",
  "foto": "/uploads/kepala/1717123456789_namafile.jpg"
}
```

**Path penyimpanan file:**
| Endpoint | Folder | Contoh path |
|---|---|---|
| `/kepala-foto` | `uploads/kepala/` | `/uploads/kepala/1717123456789_foto.jpg` |
| `/pic-foto` | `uploads/pic/` | `/uploads/pic/1717123456789_foto.jpg` |

> File bisa diakses publik via `http://localhost:8000/uploads/kepala/namafile.jpg`

---

*Dokumentasi ini di-generate dari source code backend POPDA Bulutangkis.*
