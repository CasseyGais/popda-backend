/**
 * AtletSisiInput
 * Komponen untuk memilih atlet per sisi (A atau B) dalam form laporan.
 * Urutan item di array menentukan urutan dalam tim (index 0 = urutan 1).
 *
 * Mendukung:
 * - Tunggal 1v1: 1 atlet per sisi
 * - Ganda 2v2: 2 atlet per sisi
 * - Beregu: kosong (0 atlet), pertandingan cukup per kontingen
 */

interface AtletOption {
  id: number;
  nama_lengkap: string;
}

interface Props {
  label: string;
  value: number[];           // array atlet_id yang sudah dipilih (urutan = urutan bertanding)
  onChange: (ids: number[]) => void;
  atletList: AtletOption[];  // semua atlet yang bisa dipilih (filtered by cabor)
  disabled?: boolean;
}

export default function AtletSisiInput({ label, value, onChange, atletList, disabled }: Props) {
  // Tambah slot atlet baru
  const addSlot = () => onChange([...value, 0]); // 0 = belum dipilih

  // Update atlet di posisi tertentu
  const updateAt = (idx: number, atletId: number) => {
    const next = [...value];
    next[idx] = atletId;
    onChange(next);
  };

  // Hapus atlet di posisi tertentu
  const removeAt = (idx: number) => {
    const next = value.filter((_, i) => i !== idx);
    onChange(next);
  };

  // Atlet yang sudah dipilih di sisi ini (hindari duplikat)
  const selectedSet = new Set(value.filter(id => id !== 0));

  return (
    <div className="space-y-2">
      <p className="text-sm font-medium text-gray-700 dark:text-gray-300">{label}</p>

      {value.length === 0 && (
        <p className="text-xs text-gray-400 italic">
          Kosong — pertandingan beregu (hanya kontingen)
        </p>
      )}

      {value.map((atletId, idx) => (
        <div key={idx} className="flex items-center gap-2">
          <span className="text-xs text-gray-400 w-4 shrink-0">{idx + 1}.</span>
          <select
            value={atletId || ""}
            onChange={e => updateAt(idx, Number(e.target.value))}
            disabled={disabled}
            className="flex-1 h-9 rounded-lg border border-gray-300 bg-transparent px-2 text-sm text-gray-800 dark:border-gray-700 dark:text-white dark:bg-gray-900 disabled:opacity-50"
          >
            <option value="">-- Pilih Atlet --</option>
            {atletList.map(a => (
              <option
                key={a.id}
                value={a.id}
                // Disable jika atlet sudah dipilih di slot lain di sisi yang sama
                disabled={selectedSet.has(a.id) && a.id !== atletId}
              >
                {a.nama_lengkap}
              </option>
            ))}
          </select>
          <button
            type="button"
            onClick={() => removeAt(idx)}
            className="shrink-0 p-1 text-red-500 hover:text-red-700 dark:hover:text-red-400 transition-colors"
            title="Hapus atlet ini"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      ))}

      <button
        type="button"
        onClick={addSlot}
        disabled={disabled}
        className="inline-flex items-center gap-1.5 text-xs text-brand-600 dark:text-brand-400 hover:underline disabled:opacity-40 disabled:cursor-not-allowed"
      >
        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
        </svg>
        Tambah Atlet
      </button>

      {disabled && (
        <p className="text-xs text-amber-600 dark:text-amber-400">Pilih cabor terlebih dahulu</p>
      )}
    </div>
  );
}
