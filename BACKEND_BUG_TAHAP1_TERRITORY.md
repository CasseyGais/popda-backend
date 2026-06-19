# Bug Report: `GET /admin/tahap1` Mengabaikan `?territory_id`

**Tanggal ditemukan:** 2026-06-04  
**Modul:** Tahap I — Entry By Sport  
**Endpoint:** `GET /admin/tahap1`  
**Severity:** 🔴 Critical — Superadmin tidak bisa lihat data per kontingen

---

## Ringkasan Masalah

Ketika Superadmin mengirim request `GET /admin/tahap1?territory_id=9`, backend **mengabaikan** nilai `territory_id` dari query param dan tetap mengembalikan data kontingen milik JWT token (Superadmin sendiri), bukan data kontingen dari territory yang diminta.

---

## Bukti (Console Log Frontend)

```
// Superadmin pilih territory: Kota Cilegon (id: 9)

[getTahap1] REQUEST params: { territory_id: 9 }

[getTahap1] RESPONSE:
  data.kontingen_id : 2                  ← harusnya kontingen milik Kota Cilegon
  data.territory_id : 1                  ← SALAH: masih territory_id: 1
  data.nama         : "Kabupaten Tangerang"  ← SALAH: harusnya "Kota Cilegon"
  data.cabor_list   : [Bulutangkis 2/1/1]    ← data milik Kab. Tangerang, bukan Kota Cilegon
```

Frontend mengirim `territory_id=9` dengan benar, tapi backend return data `territory_id=1` (Kabupaten Tangerang).

---

## Cara Reproduce

1. Login sebagai Superadmin
2. Buka halaman **Tahap I: Entry By Sport**
3. Di territory selector (pojok kiri atas), pilih territory yang **berbeda** dari territory superadmin sendiri (misalnya: Kota Cilegon, id=9)
4. Perhatikan: nama header berubah ke "Kota Cilegon" tapi data cabor tetap sama — milik Kabupaten Tangerang

**Request yang dikirim frontend:**
```
GET http://localhost:8000/admin/tahap1?territory_id=9
Authorization: Bearer <token_superadmin>
```

**Response yang seharusnya:**
```json
{
  "success": true,
  "data": {
    "kontingen_id": <id_kontingen_kota_cilegon>,
    "territory_id": 9,
    "nama_kontingen": "Kota Cilegon",
    "tahap1_status": "DRAFT",
    "cabor_list": []
  }
}
```

**Response yang aktual dari backend:**
```json
{
  "success": true,
  "data": {
    "kontingen_id": 2,
    "territory_id": 1,
    "nama_kontingen": "Kabupaten Tangerang",
    "tahap1_status": "SUBMITTED",
    "cabor_list": [{ "cabor_id": 6, "putra": 2, "putri": 1, "pelatih": 1, ... }]
  }
}
```

---

## Perilaku yang Diharapkan (Sesuai Dokumentasi)

Mengacu ke **TAHAP1_DOCUMENTATION.md** — Bagian 8 (Superadmin — Akses via Territory):

> **Superadmin:** Backend **selalu** membaca `?territory_id=X` dari query param, kemudian lookup `kontingen_id` dari tabel `kontingen WHERE territory_id = X`.
>
> Nilai `kontingen_id` dari JWT **diabaikan** untuk superadmin — yang dipakai adalah `territory_id` dari query param.

Alur yang benar:
```
Superadmin request → GET /admin/tahap1?territory_id=9
                                ↓
         SELECT id FROM kontingen WHERE territory_id = 9
                                ↓
                  dapat kontingen_id = <id_kota_cilegon>
                                ↓
         SELECT * FROM trx_kontingen_cabor WHERE kontingen_id = <id_kota_cilegon>
                                ↓
                  return data milik Kota Cilegon
```

---

## Endpoint Lain yang Kemungkinan Sama

Semua endpoint Tahap 1 yang menerima `?territory_id` perlu dicek:

| Method | Endpoint | Status |
|--------|----------|--------|
| `GET`  | `/admin/tahap1?territory_id=X` | 🔴 **Bug dikonfirmasi** |
| `PUT`  | `/admin/tahap1?territory_id=X` | ⚠️ Belum dicek |
| `DELETE` | `/admin/tahap1/:cabor_id?territory_id=X` | ⚠️ Belum dicek |
| `POST` | `/admin/tahap1/submit?territory_id=X` | ⚠️ Belum dicek |
| `GET`  | `/admin/tahap1/export/pdf?territory_id=X` | ⚠️ Belum dicek |
| `GET`  | `/admin/tahap1/export/excel?territory_id=X` | ⚠️ Belum dicek |

Bug yang sama kemungkinan juga ada di:
- `GET /admin/tahap2?territory_id=X`
- `GET /admin/tahap3?territory_id=X`
- `GET /admin/master/pelatih?territory_id=X`
- `GET /admin/master/official?territory_id=X`
- `GET /admin/tahap3/atlet?territory_id=X`

---

## Root Cause (Hipotesis)

Backend mungkin melakukan salah satu dari ini di handler `GET /admin/tahap1`:

```go
// ❌ SALAH — pakai kontingen_id dari JWT saja, territory_id diabaikan
kontingenID := claims.KontingenID

// ✅ BENAR — cek dulu apakah superadmin, kalau iya pakai territory_id
if claims.KontingenID == 0 {
    // Superadmin: resolve dari territory_id query param
    territoryID := c.QueryParam("territory_id")
    if territoryID == "" {
        return error("Superadmin wajib kirim territory_id")
    }
    kontingen, _ := db.Where("territory_id = ?", territoryID).First(&kontingen)
    kontingenID = kontingen.ID
} else {
    // Admin biasa: pakai dari JWT
    kontingenID = claims.KontingenID
}
```

Kemungkinan kondisi check superadmin tidak berjalan, atau kondisinya salah (misalnya check `role == "superadmin"` tapi JWT superadmin tidak punya field role yang sesuai).

---

## Fix yang Diperlukan di Backend

1. **Cek apakah user adalah superadmin** — bisa dari `claims.KontingenID == 0` atau dari claims role
2. **Jika superadmin**: wajib baca `territory_id` dari query param → lookup `kontingen_id` dari tabel `kontingen`
3. **Jika admin biasa**: pakai `kontingen_id` dari JWT seperti biasa
4. Terapkan logika yang sama ke **semua endpoint** yang menerima `?territory_id`

---

## Catatan Frontend

Frontend **sudah benar** mengirim `territory_id`:

```typescript
// src/modules/Bysport/service.ts
export const getTahap1 = (territoryId?: number): Promise<Tahap1Response> => {
  const params = territoryId ? { territory_id: territoryId } : {};
  return api.get("/admin/tahap1", { params }).then(r => r.data);
};

// src/modules/Bysport/pages/MainPage.tsx
const isSuperAdmin = can("*");  // true jika permissions includes "*"
const territoryId = isSuperAdmin ? currentTerritory?.id : undefined;

// saat territory berubah, useEffect re-fetch dengan territoryId baru
useEffect(() => {
  if (isSuperAdmin && !currentTerritory?.id) return;
  fetchAll(territoryId);
}, [isSuperAdmin, currentTerritory?.id, fetchAll, territoryId]);
```

Dari console log terbukti:
- `territory_id: 9` terkirim dengan benar di request
- Backend tetap return `territory_id: 1` di response → bug di sisi backend
