# Export PDF & Excel — Tahap 1, 2, 3

> **Base URL:** `http://localhost:8000`  
> **Auth:** Semua endpoint wajib kirim header:
> ```
> Authorization: Bearer <token>
> ```
> **Format response:** File binary (bukan JSON)  
> **Superadmin:** Wajib kirim `?territory_id=X` di semua endpoint export

---

## Daftar Isi

1. [Overview](#1-overview)
2. [Endpoint Export Tahap 1 — Entry By Sport](#2-endpoint-export-tahap-1--entry-by-sport)
3. [Endpoint Export Tahap 2 — Entry By Number](#3-endpoint-export-tahap-2--entry-by-number)
4. [Endpoint Export Tahap 3 — Entry By Name](#4-endpoint-export-tahap-3--entry-by-name)
5. [Response Headers & Download](#5-response-headers--download)
6. [Contoh Implementasi Frontend](#6-contoh-implementasi-frontend)
7. [Superadmin — Query Parameter](#7-superadmin--query-parameter)
8. [Isi Konten File per Tahap](#8-isi-konten-file-per-tahap)
9. [Error Handling](#9-error-handling)
10. [Ringkasan Endpoint](#10-ringkasan-endpoint)

---

## 1. Overview

Setiap tahap memiliki dua endpoint export:
- `GET /admin/tahapX/export/pdf` → download PDF
- `GET /admin/tahapX/export/excel` → download Excel (`.xlsx`)

File yang dihasilkan berisi data lengkap milik kontingen untuk tahap tersebut, siap cetak.

**Query params:**

| Param | Keterangan |
|---|---|
| `territory_id` | Wajib untuk superadmin — menentukan kontingen mana yang di-export |

---

## 2. Endpoint Export Tahap 1 — Entry By Sport

### `GET /admin/tahap1/export/pdf` 🔒

Download rekap data cabang olahraga yang didaftarkan kontingen di Tahap 1 dalam format PDF.

**Contoh request:**
```
GET /admin/tahap1/export/pdf
Authorization: Bearer <token>

# Superadmin:
GET /admin/tahap1/export/pdf?territory_id=2
```

**Response:**
- Status: `200 OK`
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename="tahap1_<nama_kontingen>_<tanggal>.pdf"`
- Body: file PDF binary

---

### `GET /admin/tahap1/export/excel` 🔒

Download rekap data Tahap 1 dalam format Excel (`.xlsx`).

**Contoh request:**
```
GET /admin/tahap1/export/excel
Authorization: Bearer <token>

# Superadmin:
GET /admin/tahap1/export/excel?territory_id=2
```

**Response:**
- Status: `200 OK`
- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Content-Disposition: `attachment; filename="tahap1_<nama_kontingen>_<tanggal>.xlsx"`
- Body: file Excel binary

---

## 3. Endpoint Export Tahap 2 — Entry By Number

### `GET /admin/tahap2/export/pdf` 🔒

Download daftar nomor pertandingan yang sudah didaftarkan kontingen di Tahap 2 dalam format PDF.

**Syarat:** `tahap1_status` harus `SUBMITTED` (sama dengan syarat akses Tahap 2).

**Contoh request:**
```
GET /admin/tahap2/export/pdf
Authorization: Bearer <token>

# Superadmin:
GET /admin/tahap2/export/pdf?territory_id=2
```

**Response:**
- Status: `200 OK`
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename="tahap2_<nama_kontingen>_<tanggal>.pdf"`

---

### `GET /admin/tahap2/export/excel` 🔒

Download daftar nomor pertandingan Tahap 2 dalam format Excel.

**Contoh request:**
```
GET /admin/tahap2/export/excel
Authorization: Bearer <token>
```

**Response:**
- Status: `200 OK`
- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Content-Disposition: `attachment; filename="tahap2_<nama_kontingen>_<tanggal>.xlsx"`

---

## 4. Endpoint Export Tahap 3 — Entry By Name

Tahap 3 memiliki tiga sub-data (atlet, pelatih, official). Export bisa dilakukan per sub-data atau sekaligus (semua dalam satu file).

### `GET /admin/tahap3/export/pdf` 🔒

Download semua data Tahap 3 (atlet + pelatih + official) dalam satu file PDF dengan section terpisah.

**Contoh request:**
```
GET /admin/tahap3/export/pdf
Authorization: Bearer <token>

# Superadmin:
GET /admin/tahap3/export/pdf?territory_id=2
```

**Response:**
- Content-Type: `application/pdf`
- Content-Disposition: `attachment; filename="tahap3_<nama_kontingen>_<tanggal>.pdf"`

---

### `GET /admin/tahap3/export/excel` 🔒

Download semua data Tahap 3 dalam satu file Excel dengan sheet terpisah per sub-data.

**Sheet yang ada dalam file:**
1. `Atlet` — daftar semua atlet
2. `Pelatih` — daftar semua pelatih
3. `Official` — daftar semua official

**Contoh request:**
```
GET /admin/tahap3/export/excel
Authorization: Bearer <token>
```

**Response:**
- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Content-Disposition: `attachment; filename="tahap3_<nama_kontingen>_<tanggal>.xlsx"`

---

### Export per Sub-Data (opsional)

Jika backend ingin mendukung export per sub-data saja:

| Method | URL | Keterangan |
|---|---|---|
| GET | `/admin/tahap3/export/atlet/pdf` | Hanya data atlet — PDF |
| GET | `/admin/tahap3/export/atlet/excel` | Hanya data atlet — Excel |
| GET | `/admin/tahap3/export/pelatih/pdf` | Hanya data pelatih — PDF |
| GET | `/admin/tahap3/export/pelatih/excel` | Hanya data pelatih — Excel |
| GET | `/admin/tahap3/export/official/pdf` | Hanya data official — PDF |
| GET | `/admin/tahap3/export/official/excel` | Hanya data official — Excel |

> Endpoint per sub-data ini **opsional** — frontend bisa pakai endpoint gabungan saja jika cukup.

---

## 5. Response Headers & Download

Backend harus mengirim headers berikut agar browser otomatis trigger download:

### PDF

```http
HTTP/1.1 200 OK
Content-Type: application/pdf
Content-Disposition: attachment; filename="tahap1_Kab_Tangerang_2026-06-03.pdf"
Content-Length: 45678
Cache-Control: no-store
```

### Excel

```http
HTTP/1.1 200 OK
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="tahap1_Kab_Tangerang_2026-06-03.xlsx"
Content-Length: 12345
Cache-Control: no-store
```

**Konvensi nama file:**
```
tahap1_<nama_kontingen_underscore>_<YYYY-MM-DD>.<ekstensi>

Contoh:
  tahap1_Kab_Tangerang_2026-06-03.pdf
  tahap2_Kota_Serang_2026-06-03.xlsx
  tahap3_Kab_Lebak_2026-06-03.pdf
```

---

## 6. Contoh Implementasi Frontend

Frontend memanggil endpoint ini dengan `fetch` dan membuat link download dari Blob:

```typescript
// src/utils/exportHelper.ts

/**
 * Download file dari endpoint export backend.
 * Pakai Blob API — tidak perlu library tambahan.
 */
export async function downloadExport(
  url: string,
  token: string,
  filename: string
): Promise<void> {
  const res = await fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || `Export gagal (${res.status})`);
  }

  const blob = await res.blob();
  const objectUrl = URL.createObjectURL(blob);

  // Buat link sementara dan klik otomatis
  const link = document.createElement("a");
  link.href = objectUrl;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(objectUrl);
}
```

**Penggunaan di Bysport/MainPage.tsx:**

```typescript
import { downloadExport } from "../../../utils/exportHelper";
import { useAuth } from "../../../context/AuthContext";

// Di dalam komponen:
const { token } = useAuth();

const handleExportPDF = async () => {
  const url = territoryId
    ? `http://localhost:8000/admin/tahap1/export/pdf?territory_id=${territoryId}`
    : `http://localhost:8000/admin/tahap1/export/pdf`;

  await downloadExport(url, token!, `tahap1_${kontigenName}_${today}.pdf`);
};

const handleExportExcel = async () => {
  const url = territoryId
    ? `http://localhost:8000/admin/tahap1/export/excel?territory_id=${territoryId}`
    : `http://localhost:8000/admin/tahap1/export/excel`;

  await downloadExport(url, token!, `tahap1_${kontigenName}_${today}.xlsx`);
};
```

**Tombol di UI (ditambahkan ke header halaman):**

```tsx
{/* Tombol Export — tampil jika ada data */}
{!loading && entries.length > 0 && (
  <div className="flex items-center gap-2">
    <button
      onClick={handleExportPDF}
      className="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-red-200 text-red-600 hover:bg-red-50 dark:border-red-800/40 dark:hover:bg-red-900/20 text-sm font-medium transition-colors"
    >
      {/* PDF icon */}
      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
          d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
      </svg>
      PDF
    </button>
    <button
      onClick={handleExportExcel}
      className="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-green-200 text-green-600 hover:bg-green-50 dark:border-green-800/40 dark:hover:bg-green-900/20 text-sm font-medium transition-colors"
    >
      {/* Excel icon */}
      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
          d="M3 10h18M3 14h18M10 3v18M14 3v18M5 3h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2z" />
      </svg>
      Excel
    </button>
  </div>
)}
```

---

## 7. Superadmin — Query Parameter

Sama seperti endpoint data biasa, semua endpoint export wajib menerima `?territory_id=X` untuk superadmin.

| Role | Contoh URL |
|---|---|
| Admin biasa | `GET /admin/tahap1/export/pdf` |
| Superadmin | `GET /admin/tahap1/export/pdf?territory_id=3` |

**Error jika superadmin tidak kirim territory_id (400):**
```json
{ "success": false, "message": "Superadmin wajib kirim query parameter territory_id" }
```

**Error jika tidak ada data untuk di-export (404):**
```json
{ "success": false, "message": "Tidak ada data untuk di-export" }
```

---

## 8. Isi Konten File per Tahap

### Tahap 1 — PDF & Excel

**Header dokumen:**
- Judul: `REKAP ENTRY BY SPORT — POPDA 2026`
- Nama kontingen / kabupaten-kota
- Tanggal cetak
- Status: `DRAFT` atau `SUBMITTED`

**Tabel data:**

| No | Cabang Olahraga | Atlet Putra | Atlet Putri | Pelatih | Total Atlet | Total Personel |
|---|---|---|---|---|---|---|
| 1 | Bulutangkis | 3 | 2 | 1 | 5 | 6 |
| 2 | Atletik | 5 | 4 | 2 | 9 | 11 |

**Baris total (footer tabel):**

| | **Total** | 8 | 6 | 3 | 14 | 17 |

---

### Tahap 2 — PDF & Excel

**Header dokumen:**
- Judul: `REKAP ENTRY BY NUMBER — POPDA 2026`
- Nama kontingen
- Tanggal cetak
- Status

**Tabel data (dikelompokkan per cabor):**

| No | Cabor | Nomor Pertandingan | Jenis Kelamin | Tipe |
|---|---|---|---|---|
| 1 | Bulutangkis | Tunggal | PUTRA | INDIVIDU |
| 2 | Bulutangkis | Ganda | PUTRA | BEREGU |
| 3 | Atletik | 100M | PUTRA | INDIVIDU |

**Ringkasan:**
- Total nomor terdaftar: N

---

### Tahap 3 — PDF

**PDF berisi 3 section berurutan:**

**Section 1 — Daftar Atlet**

| No | Nama Lengkap | JK | NISN | Sekolah | Kab/Kota |
|---|---|---|---|---|---|
| 1 | Ahmad Fauzi | L | 1234567890 | SMP N 1 Tangerang | Kab. Tangerang |

**Section 2 — Daftar Pelatih**

| No | Nama Lengkap | JK | Jabatan | No. HP | Kab/Kota |
|---|---|---|---|---|---|
| 1 | Budi Pelatih | L | Pelatih Kepala | 08111111111 | Kota Tangerang |

**Section 3 — Daftar Official**

| No | Nama Lengkap | JK | Jabatan | No. HP | Kab/Kota |
|---|---|---|---|---|---|
| 1 | Siti Official | P | Manajer Tim | 08222222222 | Kota Serang |

---

### Tahap 3 — Excel

**3 sheet terpisah dalam satu file `.xlsx`:**

**Sheet "Atlet":**

| No | Nama Lengkap | JK | Tanggal Lahir | NISN | Sekolah | Kelas | Kab/Kota | No. HP |
|---|---|---|---|---|---|---|---|---|

**Sheet "Pelatih":**

| No | Nama Lengkap | JK | Tanggal Lahir | NIK | Jabatan | No. HP | Email | Kab/Kota |
|---|---|---|---|---|---|---|---|---|

**Sheet "Official":**

| No | Nama Lengkap | JK | Tanggal Lahir | NIK | Jabatan | No. HP | Email | Kab/Kota |
|---|---|---|---|---|---|---|---|---|

---

## 9. Error Handling

### Error umum

| Status | Kondisi | Response Body |
|---|---|---|
| `400` | Superadmin tidak kirim `territory_id` | `{ "success": false, "message": "Superadmin wajib kirim query parameter territory_id" }` |
| `400` | Tahap 1 belum submit (untuk export Tahap 2) | `{ "success": false, "message": "tahap 1 belum disubmit" }` |
| `404` | Tidak ada data untuk di-export | `{ "success": false, "message": "Tidak ada data untuk di-export" }` |
| `404` | Territory tidak punya kontingen | `{ "success": false, "message": "Kontingen untuk territory ini tidak ditemukan" }` |
| `500` | Gagal generate file | `{ "success": false, "message": "Gagal membuat file export" }` |

> Saat error, backend mengembalikan JSON (bukan binary) dengan `Content-Type: application/json`.

### Handling di frontend

```typescript
const handleExport = async (format: "pdf" | "excel") => {
  setExporting(true);
  try {
    await downloadExport(url, token!, filename);
  } catch (err: any) {
    alert("Export gagal: " + (err.message || "Coba lagi"));
  } finally {
    setExporting(false);
  }
};
```

---

## 10. Ringkasan Endpoint

### Tahap 1

| Method | URL | Format | Keterangan |
|---|---|---|---|
| `GET` | `/admin/tahap1/export/pdf` | PDF | Rekap cabor terdaftar |
| `GET` | `/admin/tahap1/export/excel` | XLSX | Rekap cabor terdaftar |

### Tahap 2

| Method | URL | Format | Keterangan |
|---|---|---|---|
| `GET` | `/admin/tahap2/export/pdf` | PDF | Daftar nomor terdaftar |
| `GET` | `/admin/tahap2/export/excel` | XLSX | Daftar nomor terdaftar |

### Tahap 3

| Method | URL | Format | Keterangan |
|---|---|---|---|
| `GET` | `/admin/tahap3/export/pdf` | PDF | Semua: atlet + pelatih + official |
| `GET` | `/admin/tahap3/export/excel` | XLSX | 3 sheet: Atlet, Pelatih, Official |
| `GET` | `/admin/tahap3/export/atlet/pdf` | PDF | (Opsional) Hanya atlet |
| `GET` | `/admin/tahap3/export/atlet/excel` | XLSX | (Opsional) Hanya atlet |
| `GET` | `/admin/tahap3/export/pelatih/pdf` | PDF | (Opsional) Hanya pelatih |
| `GET` | `/admin/tahap3/export/pelatih/excel` | XLSX | (Opsional) Hanya pelatih |
| `GET` | `/admin/tahap3/export/official/pdf` | PDF | (Opsional) Hanya official |
| `GET` | `/admin/tahap3/export/official/excel` | XLSX | (Opsional) Hanya official |

> Semua endpoint: `territory_id` wajib untuk superadmin, tidak perlu untuk admin biasa.

---

## Catatan Implementasi Backend (Go)

Untuk generate PDF dan Excel di Go, berikut library yang disarankan:

**PDF:**
- [`github.com/jung-kurt/gofpdf`](https://github.com/jung-kurt/gofpdf) — library PDF populer untuk Go
- [`github.com/signintech/gopdf`](https://github.com/signintech/gopdf) — alternatif dengan dukungan font Unicode lebih baik (penting untuk nama Indonesia)

**Excel:**
- [`github.com/xuri/excelize`](https://github.com/xuri/excelize) — library Excel paling lengkap untuk Go, support `.xlsx`

**Contoh handler Go (skeleton):**

```go
// GET /admin/tahap1/export/pdf
func (h *Handler) ExportTahap1PDF(c *gin.Context) {
    // 1. Ambil kontingen_id dari token (atau territory_id untuk superadmin)
    kontingenID := getKontingenIDFromContext(c)

    // 2. Query data dari DB
    caborList, err := h.repo.GetTahap1ByCabor(kontingenID)
    if err != nil {
        c.JSON(500, gin.H{"success": false, "message": "Gagal mengambil data"})
        return
    }
    if len(caborList) == 0 {
        c.JSON(404, gin.H{"success": false, "message": "Tidak ada data untuk di-export"})
        return
    }

    // 3. Generate PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    // ... isi konten PDF ...

    // 4. Set headers dan tulis response
    kontigenName := strings.ReplaceAll(kontingen.NamaKontingen, " ", "_")
    tanggal := time.Now().Format("2006-01-02")
    filename := fmt.Sprintf("tahap1_%s_%s.pdf", kontigenName, tanggal)

    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Header("Content-Type", "application/pdf")
    c.Header("Cache-Control", "no-store")

    if err := pdf.Output(c.Writer); err != nil {
        c.JSON(500, gin.H{"success": false, "message": "Gagal membuat PDF"})
    }
}

// GET /admin/tahap1/export/excel
func (h *Handler) ExportTahap1Excel(c *gin.Context) {
    kontingenID := getKontingenIDFromContext(c)
    caborList, _ := h.repo.GetTahap1ByCabor(kontingenID)

    f := excelize.NewFile()
    sheet := "Tahap 1"
    f.NewSheet(sheet)

    // Header baris
    headers := []string{"No", "Cabang Olahraga", "Atlet Putra", "Atlet Putri",
                         "Pelatih", "Total Atlet", "Total Personel"}
    for i, h := range headers {
        cell, _ := excelize.CoordinatesToCellName(i+1, 1)
        f.SetCellValue(sheet, cell, h)
    }

    // Isi data
    for rowIdx, cabor := range caborList {
        row := rowIdx + 2
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rowIdx+1)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), cabor.Nama)
        // dst...
    }

    filename := fmt.Sprintf("tahap1_%s_%s.xlsx", kontigenName, tanggal)
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

    f.Write(c.Writer)
}
```
