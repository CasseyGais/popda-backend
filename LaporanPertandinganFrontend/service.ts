import api from "../../lib/api";

const BASE_URL = "http://localhost:8000";
const getToken = () => localStorage.getItem("token") ?? "";

// ─── Enums ────────────────────────────────────────────────

export type Babak =
  | "PENYISIHAN"
  | "8_BESAR"
  | "PEREMPAT_FINAL"
  | "SEMIFINAL"
  | "FINAL"
  | "PEREBUTAN_TEMPAT_3"
  | "LAINNYA";

export type Pemenang = "TIM_A" | "TIM_B" | "DRAW";

export const BABAK_OPTIONS: { value: Babak; label: string }[] = [
  { value: "PENYISIHAN",         label: "Penyisihan" },
  { value: "8_BESAR",            label: "8 Besar" },
  { value: "PEREMPAT_FINAL",     label: "Perempat Final" },
  { value: "SEMIFINAL",          label: "Semifinal" },
  { value: "FINAL",              label: "Final" },
  { value: "PEREBUTAN_TEMPAT_3", label: "Perebutan Tempat 3" },
  { value: "LAINNYA",            label: "Lainnya" },
];

export const PEMENANG_OPTIONS: { value: Pemenang; label: string }[] = [
  { value: "TIM_A", label: "Tim A" },
  { value: "TIM_B", label: "Tim B" },
  { value: "DRAW",  label: "Seri" },
];

// ─── Types ────────────────────────────────────────────────

export interface AtletSisiItem {
  id: number;
  atlet_id: number;
  nama_lengkap: string;
  urutan: number;
}

export interface LaporanDetail {
  id: number;
  tanggal_pertandingan: string;   // YYYY-MM-DD
  waktu_pertandingan: string;     // HH:MM:SS
  venue: string;
  cabor_id: number;
  nomor_id: number;
  babak: Babak;
  kontingen_a_id: number;
  kontingen_b_id: number | null;
  hasil_pertandingan: string;
  pemenang: Pemenang;
  juara_ke: number | null;
  wasit: string;
  catatan_khusus: string | null;
  foto_bukti: string | null;
  video_bukti: string | null;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  // Join fields
  nama_cabor: string;
  nama_nomor: string;
  nama_kontingen_a: string;
  nama_kontingen_b: string | null;
  atlet_a: AtletSisiItem[];
  atlet_b: AtletSisiItem[];
}

/** Payload POST — jangan kirim created_by, diisi otomatis JWT */
export interface CreateLaporanPayload {
  tanggal_pertandingan: string;   // YYYY-MM-DD, wajib
  waktu_pertandingan: string;     // HH:MM:SS, wajib
  venue: string;                  // wajib
  cabor_id: number;               // wajib
  nomor_id: number;               // wajib
  babak: Babak;                   // wajib
  kontingen_a_id: number;         // wajib
  kontingen_b_id?: number;        // opsional
  hasil_pertandingan: string;     // wajib
  pemenang: Pemenang;             // wajib
  juara_ke?: number;              // opsional
  wasit: string;                  // wajib
  catatan_khusus?: string;        // opsional
  atlet_a?: number[];             // ordered list atlet_id sisi A
  atlet_b?: number[];             // ordered list atlet_id sisi B
}

/** Payload PUT — partial update */
export interface UpdateLaporanPayload {
  tanggal_pertandingan?: string;
  waktu_pertandingan?: string;
  venue?: string;
  cabor_id?: number;
  nomor_id?: number;
  babak?: Babak;
  kontingen_a_id?: number;
  kontingen_b_id?: number | null;
  hasil_pertandingan?: string;
  pemenang?: Pemenang;
  juara_ke?: number | null;
  wasit?: string;
  catatan_khusus?: string | null;
  /**
   * Replace semua atlet sisi A. Tidak dikirim = tidak berubah.
   * Array kosong [] = hapus semua atlet sisi A.
   */
  atlet_a?: number[];
  /**
   * Replace semua atlet sisi B. Tidak dikirim = tidak berubah.
   * Array kosong [] = hapus semua atlet sisi B.
   */
  atlet_b?: number[];
}

export interface LaporanFilter {
  tanggal?: string;     // YYYY-MM-DD
  cabor_id?: number;
  nomor_id?: number;
  babak?: Babak;
  pemenang?: Pemenang;
}

/** Data satu penandatangan untuk export PDF */
export interface TTDData {
  jabatan: string;
  nama_tercetak: string;
  nip?: string;
  signature_b64?: string; // base64 PNG dari signature pad
}

/** Payload export PDF — batch maupun single */
export interface ExportPDFPayload {
  tanggal?: string;         // filter per hari (batch only)
  cabor_id?: number;        // filter cabor (batch only)
  nomor_id?: number;        // filter nomor (batch only)
  penandatangan?: TTDData[];
}

// ─── Service ─────────────────────────────────────────────

/**
 * GET /admin/laporan-pertandingan
 * List semua laporan, diurutkan tanggal DESC, waktu DESC.
 */
export const getAllLaporan = (
  filters?: LaporanFilter
): Promise<LaporanDetail[]> =>
  api.get("/admin/laporan-pertandingan", { params: filters }).then(r => r.data.data ?? []);

/**
 * GET /admin/laporan-pertandingan/:id
 * Detail satu laporan lengkap termasuk atlet_a dan atlet_b.
 */
export const getLaporanById = (id: number): Promise<LaporanDetail> =>
  api.get(`/admin/laporan-pertandingan/${id}`).then(r => r.data.data);

/**
 * POST /admin/laporan-pertandingan
 * Buat laporan baru. created_by otomatis dari JWT — tidak perlu dikirim.
 */
export const createLaporan = (
  payload: CreateLaporanPayload
): Promise<LaporanDetail> =>
  api.post("/admin/laporan-pertandingan", payload).then(r => r.data.data);

/**
 * PUT /admin/laporan-pertandingan/:id
 * Partial update. Untuk atlet: jika field tidak ada di body → tidak berubah.
 * Kirim [] → hapus semua atlet sisi tersebut.
 */
export const updateLaporan = (
  id: number,
  payload: UpdateLaporanPayload
): Promise<LaporanDetail> =>
  api.put(`/admin/laporan-pertandingan/${id}`, payload).then(r => r.data.data);

/**
 * DELETE /admin/laporan-pertandingan/:id
 * Hard delete. Data atlet (laporan_pertandingan_atlet) ikut terhapus via CASCADE.
 */
export const deleteLaporan = (
  id: number
): Promise<{ success: boolean; message: string }> =>
  api.delete(`/admin/laporan-pertandingan/${id}`).then(r => r.data);

/**
 * PUT /admin/laporan-pertandingan/:id/foto
 * Upload foto bukti. Field name di FormData harus "foto".
 */
export const uploadFotoLaporan = async (
  id: number,
  file: File
): Promise<{ success: boolean; message: string; path: string }> => {
  const fd = new FormData();
  fd.append("foto", file);
  const r = await fetch(`${BASE_URL}/admin/laporan-pertandingan/${id}/foto`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${getToken()}` },
    body: fd,
  });
  const data = await r.json().catch(() => ({}));
  if (!r.ok || data.success === false)
    throw new Error(data.message || `Error ${r.status}`);
  return data;
};

/**
 * PUT /admin/laporan-pertandingan/:id/video
 * Upload video bukti. Field name di FormData harus "video".
 */
export const uploadVideoLaporan = async (
  id: number,
  file: File
): Promise<{ success: boolean; message: string; path: string }> => {
  const fd = new FormData();
  fd.append("video", file);
  const r = await fetch(`${BASE_URL}/admin/laporan-pertandingan/${id}/video`, {
    method: "PUT",
    headers: { Authorization: `Bearer ${getToken()}` },
    body: fd,
  });
  const data = await r.json().catch(() => ({}));
  if (!r.ok || data.success === false)
    throw new Error(data.message || `Error ${r.status}`);
  return data;
};

/**
 * POST /admin/laporan-pertandingan/:id/export/pdf
 * Export PDF satu pertandingan. Body opsional (tanpa body = tanpa tanda tangan).
 * Response: file PDF binary (blob).
 */
export const exportSatuPDF = async (
  id: number,
  payload: ExportPDFPayload = {}
): Promise<void> => {
  const res = await fetch(`${BASE_URL}/admin/laporan-pertandingan/${id}/export/pdf`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getToken()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || `Error ${res.status}`);
  }
  const blob = await res.blob();
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  a.download = `laporan_${id}.pdf`;
  a.click();
  URL.revokeObjectURL(a.href);
};

/**
 * POST /admin/laporan-pertandingan/export/pdf
 * Export PDF batch (semua / per hari / per cabor).
 * Pakai POST karena body bisa berisi data tanda tangan base64.
 * PENTING: route statis /export/pdf terdaftar sebelum /:id di backend.
 */
export const exportBatchPDF = async (
  payload: ExportPDFPayload = {}
): Promise<void> => {
  const res = await fetch(`${BASE_URL}/admin/laporan-pertandingan/export/pdf`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getToken()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || `Error ${res.status}`);
  }
  const blob = await res.blob();
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  const suffix = payload.tanggal ?? "semua";
  a.download = `laporan_pertandingan_${suffix}.pdf`;
  a.click();
  URL.revokeObjectURL(a.href);
};

// ─── Service object ───────────────────────────────────────

export const laporanPertandinganService = {
  getAll:       getAllLaporan,
  getById:      getLaporanById,
  create:       createLaporan,
  update:       updateLaporan,
  delete:       deleteLaporan,
  uploadFoto:   uploadFotoLaporan,
  uploadVideo:  uploadVideoLaporan,
  exportSatu:   exportSatuPDF,
  exportBatch:  exportBatchPDF,
};
