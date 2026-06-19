import { useState, useEffect } from "react";
import { Modal } from "../../../components/ui/modal";
import Button from "../../../components/ui/button/Button";
import Input from "../../../components/form/input/InputField";
import Label from "../../../components/form/Label";
import {
  laporanPertandinganService,
  BABAK_OPTIONS,
  PEMENANG_OPTIONS,
  type LaporanDetail,
  type Babak,
  type Pemenang,
  type CreateLaporanPayload,
  type UpdateLaporanPayload,
} from "../service";
import AtletSisiInput from "./AtletSisiInput";
import MediaUploadRow from "./MediaUploadRow";

// ─── Dropdown helpers — ambil dari endpoint master ────────

interface DropdownItem { id: number; nama: string }

async function fetchCabor(): Promise<DropdownItem[]> {
  try {
    const r = await fetch("http://localhost:8000/admin/master/cabor", {
      headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
    });
    const d = await r.json();
    return d.data ?? [];
  } catch { return []; }
}

async function fetchNomor(caborId: number): Promise<DropdownItem[]> {
  try {
    const r = await fetch(`http://localhost:8000/admin/master/nomor?cabor_id=${caborId}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
    });
    const d = await r.json();
    return d.data ?? [];
  } catch { return []; }
}

async function fetchKontingen(): Promise<DropdownItem[]> {
  try {
    const r = await fetch("http://localhost:8000/admin/kontingen", {
      headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
    });
    const d = await r.json();
    return d.data ?? [];
  } catch { return []; }
}

async function fetchAtletByCabor(caborId: number): Promise<{ id: number; nama_lengkap: string }[]> {
  try {
    const r = await fetch(`http://localhost:8000/admin/master/atlet?cabor_id=${caborId}`, {
      headers: { Authorization: `Bearer ${localStorage.getItem("token") ?? ""}` },
    });
    const d = await r.json();
    return d.data ?? [];
  } catch { return []; }
}

// ─── Props ────────────────────────────────────────────────

interface Props {
  isOpen: boolean;
  onClose: () => void;
  mode: "create" | "edit" | "view";
  data: LaporanDetail | null;
  onSuccess: (result: LaporanDetail) => void;
}

const TODAY = new Date().toISOString().slice(0, 10);

// ─── Component ───────────────────────────────────────────

export default function LaporanModal({ isOpen, onClose, mode, data, onSuccess }: Props) {
  const isView = mode === "view";

  // Dropdown options
  const [caborList, setCaborList]           = useState<DropdownItem[]>([]);
  const [nomorList, setNomorList]           = useState<DropdownItem[]>([]);
  const [kontingenList, setKontingenList]   = useState<DropdownItem[]>([]);
  const [atletList, setAtletList]           = useState<{ id: number; nama_lengkap: string }[]>([]);
  const [loadingDropdown, setLoadingDropdown] = useState(false);

  // Form fields
  const [tanggal, setTanggal]       = useState(TODAY);
  const [waktu, setWaktu]           = useState("09:00");
  const [venue, setVenue]           = useState("");
  const [caborId, setCaborId]       = useState<number | "">("");
  const [nomorId, setNomorId]       = useState<number | "">("");
  const [babak, setBabak]           = useState<Babak>("PENYISIHAN");
  const [kontingenAId, setKontingenAId] = useState<number | "">("");
  const [kontingenBId, setKontingenBId] = useState<number | "">("");
  const [hasil, setHasil]           = useState("");
  const [pemenang, setPemenang]     = useState<Pemenang>("TIM_A");
  const [juaraKe, setJuaraKe]       = useState<number | "">("");
  const [wasit, setWasit]           = useState("");
  const [catatan, setCatatan]       = useState("");

  // Atlet sisi A & B (array atlet_id)
  const [atletA, setAtletA] = useState<number[]>([]);
  const [atletB, setAtletB] = useState<number[]>([]);

  // Upload state (edit mode)
  const [localData, setLocalData]   = useState<LaporanDetail | null>(null);
  const [anyUploading, setAnyUploading] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState("");

  const activeData = localData ?? data;

  // ── Load dropdown master saat modal dibuka ─────────────
  useEffect(() => {
    if (!isOpen || isView) return;
    setLoadingDropdown(true);
    Promise.all([fetchCabor(), fetchKontingen()])
      .then(([cabors, kontingens]) => {
        setCaborList(cabors);
        setKontingenList(kontingens);
      })
      .finally(() => setLoadingDropdown(false));
  }, [isOpen, isView]);

  // ── Load nomor saat cabor berubah ──────────────────────
  useEffect(() => {
    if (!caborId) { setNomorList([]); return; }
    fetchNomor(Number(caborId)).then(setNomorList);
  }, [caborId]);

  // ── Load atlet saat cabor berubah ──────────────────────
  useEffect(() => {
    if (!caborId) { setAtletList([]); return; }
    fetchAtletByCabor(Number(caborId)).then(setAtletList);
  }, [caborId]);

  // ── Populate form saat modal dibuka ───────────────────
  useEffect(() => {
    if (!isOpen) return;
    setError("");
    setLocalData(null);
    if (data && mode !== "create") {
      setTanggal(data.tanggal_pertandingan);
      setWaktu(data.waktu_pertandingan.slice(0, 5));
      setVenue(data.venue);
      setCaborId(data.cabor_id);
      setNomorId(data.nomor_id);
      setBabak(data.babak);
      setKontingenAId(data.kontingen_a_id);
      setKontingenBId(data.kontingen_b_id ?? "");
      setHasil(data.hasil_pertandingan);
      setPemenang(data.pemenang);
      setJuaraKe(data.juara_ke ?? "");
      setWasit(data.wasit);
      setCatatan(data.catatan_khusus ?? "");
      setAtletA(data.atlet_a.map(a => a.atlet_id));
      setAtletB(data.atlet_b.map(a => a.atlet_id));
    } else {
      setTanggal(TODAY);
      setWaktu("09:00");
      setVenue("");
      setCaborId("");
      setNomorId("");
      setBabak("PENYISIHAN");
      setKontingenAId("");
      setKontingenBId("");
      setHasil("");
      setPemenang("TIM_A");
      setJuaraKe("");
      setWasit("");
      setCatatan("");
      setAtletA([]);
      setAtletB([]);
    }
  }, [isOpen, data, mode]);

  // ── Validasi & Save ───────────────────────────────────
  const handleSave = async () => {
    if (!tanggal || !waktu || !venue || caborId === "" || nomorId === "" || !babak || kontingenAId === "" || !hasil || !pemenang || !wasit) {
      setError("Field bertanda * wajib diisi");
      return;
    }
    setLoading(true);
    setError("");
    try {
      let result: LaporanDetail;
      if (mode === "create") {
        const payload: CreateLaporanPayload = {
          tanggal_pertandingan: tanggal,
          waktu_pertandingan:   waktu + ":00",
          venue,
          cabor_id:             Number(caborId),
          nomor_id:             Number(nomorId),
          babak,
          kontingen_a_id:       Number(kontingenAId),
          ...(kontingenBId !== "" && { kontingen_b_id: Number(kontingenBId) }),
          hasil_pertandingan:   hasil,
          pemenang,
          ...(juaraKe !== "" && { juara_ke: Number(juaraKe) }),
          wasit,
          ...(catatan && { catatan_khusus: catatan }),
          atlet_a: atletA,
          atlet_b: atletB,
        };
        result = await laporanPertandinganService.create(payload);
      } else {
        const payload: UpdateLaporanPayload = {
          tanggal_pertandingan: tanggal,
          waktu_pertandingan:   waktu + ":00",
          venue,
          cabor_id:             Number(caborId),
          nomor_id:             Number(nomorId),
          babak,
          kontingen_a_id:       Number(kontingenAId),
          kontingen_b_id:       kontingenBId !== "" ? Number(kontingenBId) : null,
          hasil_pertandingan:   hasil,
          pemenang,
          juara_ke:             juaraKe !== "" ? Number(juaraKe) : null,
          wasit,
          catatan_khusus:       catatan || null,
          atlet_a:              atletA,
          atlet_b:              atletB,
        };
        result = await laporanPertandinganService.update(data!.id, payload);
      }
      onSuccess(result);
      onClose();
    } catch (e: any) {
      setError(e.message || "Gagal menyimpan laporan");
    } finally {
      setLoading(false);
    }
  };

  // ── Upload foto/video ─────────────────────────────────
  const handleUploadFoto = async (file: File) => {
    const id = (localData ?? data)!.id;
    setAnyUploading(true);
    try {
      const res = await laporanPertandinganService.uploadFoto(id, file);
      setLocalData(prev => ({ ...(prev ?? data!), foto_bukti: res.path }));
    } finally {
      setAnyUploading(false);
    }
  };

  const handleUploadVideo = async (file: File) => {
    const id = (localData ?? data)!.id;
    setAnyUploading(true);
    try {
      const res = await laporanPertandinganService.uploadVideo(id, file);
      setLocalData(prev => ({ ...(prev ?? data!), video_bukti: res.path }));
    } finally {
      setAnyUploading(false);
    }
  };

  // ── Helpers tampilan ──────────────────────────────────
  const fmtDate = (d: string) =>
    new Date(d + "T00:00:00").toLocaleDateString("id-ID", {
      day: "2-digit", month: "long", year: "numeric",
    });

  const BABAK_LABEL: Record<Babak, string> = Object.fromEntries(
    BABAK_OPTIONS.map(o => [o.value, o.label])
  ) as Record<Babak, string>;

  const PEMENANG_LABEL: Record<Pemenang, string> = Object.fromEntries(
    PEMENANG_OPTIONS.map(o => [o.value, o.label])
  ) as Record<Pemenang, string>;

  return (
    <Modal isOpen={isOpen} onClose={onClose} className="max-w-[680px] m-4">
      <div className="no-scrollbar relative w-full overflow-y-auto rounded-3xl bg-white dark:bg-gray-900 p-6 lg:p-8">
        {/* Header */}
        <div className="mb-5 pr-8">
          <h4 className="text-xl font-semibold text-gray-800 dark:text-white">
            {mode === "create" ? "Tambah Laporan Pertandingan"
              : mode === "edit" ? "Edit Laporan Pertandingan"
              : "Detail Laporan Pertandingan"}
          </h4>
          {activeData && (
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              #{activeData.id} · {activeData.nama_cabor} — {activeData.nama_nomor}
            </p>
          )}
        </div>

        {error && (
          <div className="mb-4 px-4 py-3 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800/40 text-sm text-red-600 dark:text-red-400">
            {error}
          </div>
        )}

        <div className="custom-scrollbar max-h-[70vh] overflow-y-auto space-y-6 pr-1">

          {/* ── Informasi Pertandingan ── */}
          <section>
            <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
              Informasi Pertandingan
            </h5>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label>Tanggal <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{fmtDate(activeData!.tanggal_pertandingan)}</p>
                  : <Input type="date" value={tanggal} onChange={e => setTanggal(e.target.value)} />
                }
              </div>
              <div>
                <Label>Waktu <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.waktu_pertandingan.slice(0,5)} WIB</p>
                  : <Input type="time" value={waktu} onChange={e => setWaktu(e.target.value)} />
                }
              </div>
              <div className="col-span-2">
                <Label>Venue / Lapangan <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.venue}</p>
                  : <Input type="text" value={venue} onChange={e => setVenue(e.target.value)} placeholder="GOR Pemuda Serang" />
                }
              </div>
            </div>
          </section>

          {/* ── Cabor & Nomor ── */}
          <section>
            <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
              Cabang Olahraga
            </h5>
            {isView ? (
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Cabor</Label>
                  <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.nama_cabor}</p>
                </div>
                <div>
                  <Label>Nomor / Kelas</Label>
                  <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.nama_nomor}</p>
                </div>
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Cabor <span className="text-red-500">*</span></Label>
                  {loadingDropdown ? (
                    <p className="text-xs text-gray-400 mt-1">Memuat...</p>
                  ) : (
                    <select
                      value={caborId}
                      onChange={e => { setCaborId(Number(e.target.value)); setNomorId(""); }}
                      className="w-full h-11 rounded-lg border border-gray-300 bg-transparent px-3 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900"
                    >
                      <option value="">-- Pilih Cabor --</option>
                      {caborList.map(c => <option key={c.id} value={c.id}>{c.nama}</option>)}
                    </select>
                  )}
                </div>
                <div>
                  <Label>Nomor / Kelas <span className="text-red-500">*</span></Label>
                  <select
                    value={nomorId}
                    onChange={e => setNomorId(Number(e.target.value))}
                    disabled={!caborId}
                    className="w-full h-11 rounded-lg border border-gray-300 bg-transparent px-3 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900 disabled:opacity-50"
                  >
                    <option value="">-- Pilih Nomor --</option>
                    {nomorList.map(n => <option key={n.id} value={n.id}>{n.nama}</option>)}
                  </select>
                </div>
              </div>
            )}
          </section>

          {/* ── Babak & Kontingen ── */}
          <section>
            <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
              Babak & Tim
            </h5>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label>Babak <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{BABAK_LABEL[activeData!.babak]}</p>
                  : (
                    <select
                      value={babak}
                      onChange={e => setBabak(e.target.value as Babak)}
                      className="w-full h-11 rounded-lg border border-gray-300 bg-transparent px-3 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900"
                    >
                      {BABAK_OPTIONS.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
                    </select>
                  )
                }
              </div>
              <div>
                <Label>Tim A (Kontingen) <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.nama_kontingen_a}</p>
                  : (
                    <select
                      value={kontingenAId}
                      onChange={e => setKontingenAId(Number(e.target.value))}
                      className="w-full h-11 rounded-lg border border-gray-300 bg-transparent px-3 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900"
                    >
                      <option value="">-- Pilih Kontingen A --</option>
                      {kontingenList.map(k => <option key={k.id} value={k.id}>{k.nama}</option>)}
                    </select>
                  )
                }
              </div>
              <div>
                <Label>Tim B (Kontingen)</Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.nama_kontingen_b ?? "—"}</p>
                  : (
                    <select
                      value={kontingenBId}
                      onChange={e => setKontingenBId(e.target.value ? Number(e.target.value) : "")}
                      className="w-full h-11 rounded-lg border border-gray-300 bg-transparent px-3 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900"
                    >
                      <option value="">-- Tanpa lawan --</option>
                      {kontingenList.map(k => <option key={k.id} value={k.id}>{k.nama}</option>)}
                    </select>
                  )
                }
              </div>
            </div>
          </section>

          {/* ── Atlet Sisi A & B ── */}
          <section>
            <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
              Atlet Bertanding
            </h5>
            {isView ? (
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Atlet Sisi A</Label>
                  {activeData!.atlet_a.length === 0
                    ? <p className="text-xs text-gray-400 mt-1">Tidak ada atlet (beregu)</p>
                    : (
                      <ul className="mt-1 space-y-0.5">
                        {activeData!.atlet_a.map(a => (
                          <li key={a.id} className="text-sm text-gray-800 dark:text-white">
                            {a.urutan}. {a.nama_lengkap}
                          </li>
                        ))}
                      </ul>
                    )
                  }
                </div>
                <div>
                  <Label>Atlet Sisi B</Label>
                  {activeData!.atlet_b.length === 0
                    ? <p className="text-xs text-gray-400 mt-1">Tidak ada atlet (beregu)</p>
                    : (
                      <ul className="mt-1 space-y-0.5">
                        {activeData!.atlet_b.map(a => (
                          <li key={a.id} className="text-sm text-gray-800 dark:text-white">
                            {a.urutan}. {a.nama_lengkap}
                          </li>
                        ))}
                      </ul>
                    )
                  }
                </div>
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-4">
                <AtletSisiInput
                  label="Atlet Sisi A"
                  value={atletA}
                  onChange={setAtletA}
                  atletList={atletList}
                  disabled={!caborId}
                />
                <AtletSisiInput
                  label="Atlet Sisi B"
                  value={atletB}
                  onChange={setAtletB}
                  atletList={atletList}
                  disabled={!caborId}
                />
              </div>
            )}
            <p className="text-xs text-gray-400 mt-2">
              Urutan array menentukan urutan dalam tim. Kosongkan untuk pertandingan beregu (cukup kontingen).
            </p>
          </section>

          {/* ── Hasil & Pemenang ── */}
          <section>
            <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
              Hasil Pertandingan
            </h5>
            <div className="space-y-4">
              <div>
                <Label>Hasil Pertandingan <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.hasil_pertandingan}</p>
                  : (
                    <Input
                      type="text"
                      value={hasil}
                      onChange={e => setHasil(e.target.value)}
                      placeholder="21-18, 18-21, 21-15"
                    />
                  )
                }
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Pemenang <span className="text-red-500">*</span></Label>
                  {isView
                    ? <p className="text-sm text-gray-800 dark:text-white mt-1">{PEMENANG_LABEL[activeData!.pemenang]}</p>
                    : (
                      <select
                        value={pemenang}
                        onChange={e => setPemenang(e.target.value as Pemenang)}
                        className="w-full h-11 rounded-lg border border-gray-300 bg-transparent px-3 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900"
                      >
                        {PEMENANG_OPTIONS.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
                      </select>
                    )
                  }
                </div>
                <div>
                  <Label>Juara Ke</Label>
                  {isView
                    ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.juara_ke ? `Juara ${activeData!.juara_ke}` : "—"}</p>
                    : (
                      <Input
                        type="number"
                        value={juaraKe}
                        onChange={e => setJuaraKe(e.target.value ? Number(e.target.value) : "")}
                        placeholder="1 / 2 / 3 (opsional)"
                      />
                    )
                  }
                </div>
              </div>
            </div>
          </section>

          {/* ── Wasit & Catatan ── */}
          <section>
            <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
              Wasit & Catatan
            </h5>
            <div className="space-y-4">
              <div>
                <Label>Wasit / Juri <span className="text-red-500">*</span></Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.wasit}</p>
                  : <Input type="text" value={wasit} onChange={e => setWasit(e.target.value)} placeholder="Nama wasit" />
                }
              </div>
              <div>
                <Label>Catatan Khusus</Label>
                {isView
                  ? <p className="text-sm text-gray-800 dark:text-white mt-1">{activeData!.catatan_khusus || "—"}</p>
                  : (
                    <textarea
                      value={catatan}
                      onChange={e => setCatatan(e.target.value)}
                      rows={2}
                      placeholder="Opsional"
                      className="w-full rounded-lg border border-gray-300 bg-transparent px-3 py-2 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900 resize-none"
                    />
                  )
                }
              </div>
            </div>
          </section>

          {/* ── Foto & Video Bukti ── */}
          {(mode === "edit" || mode === "view") && activeData && (
            <section>
              <h5 className="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500 mb-3 border-b border-gray-100 dark:border-gray-800 pb-1.5">
                Bukti Media
              </h5>
              <div className="space-y-3">
                <MediaUploadRow
                  label="Foto Bukti"
                  currentPath={activeData.foto_bukti}
                  accept="image/*"
                  readonly={isView}
                  onUpload={mode === "edit" ? handleUploadFoto : undefined}
                  uploading={anyUploading}
                />
                <MediaUploadRow
                  label="Video Bukti"
                  currentPath={activeData.video_bukti}
                  accept="video/*"
                  readonly={isView}
                  onUpload={mode === "edit" ? handleUploadVideo : undefined}
                  uploading={anyUploading}
                />
              </div>
              {mode === "create" && (
                <p className="text-xs text-gray-400 mt-2">
                  Foto dan video dapat diupload setelah laporan tersimpan melalui menu Edit.
                </p>
              )}
            </section>
          )}

        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-3 mt-6 pt-4 border-t border-gray-100 dark:border-gray-800">
          <Button size="sm" variant="outline" onClick={onClose} disabled={loading || anyUploading}>
            {isView ? "Tutup" : "Batal"}
          </Button>
          {!isView && (
            <Button
              size="sm"
              onClick={handleSave}
              disabled={loading || anyUploading}
              className="bg-brand-500 hover:bg-brand-600 text-white"
            >
              {loading ? "Menyimpan..." : mode === "create" ? "Simpan" : "Perbarui"}
            </Button>
          )}
        </div>
      </div>
    </Modal>
  );
}
