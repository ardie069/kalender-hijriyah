import type { Method as HijriMethod } from "@/types/hijri";

export const HIJRI_MONTHS_INDONESIA_GRAMMAR = [
  { id: 1, name: "Muharam", desc: "Awal Tahun" },
  { id: 2, name: "Safar", desc: "Bulan Perjalanan" },
  { id: 3, name: "Rabiulawal", desc: "Maulid Nabi" },
  { id: 4, name: "Rabiulakhir", desc: "Musim Dingin" },
  { id: 5, name: "Jumadilawal", desc: "Awal Kemarau" },
  { id: 6, name: "Jumadilakhir", desc: "Akhir Kemarau" },
  { id: 7, name: "Rajab", desc: "Bulan Mulia" },
  { id: 8, name: "Syakban", desc: "Persiapan" },
  { id: 9, name: "Ramadan", desc: "Puasa Wajib", isSpecial: true },
  { id: 10, name: "Syawal", desc: "Idul Fitri", isSpecial: true },
  { id: 11, name: "Zulkaidah", desc: "Bulan Tenang" },
  { id: 12, name: "Zulhijah", desc: "Ibadah Haji", isSpecial: true },
];

export type HijriMonth = (typeof HIJRI_MONTHS_INDONESIA_GRAMMAR)[number];

export const METHODS: { id: HijriMethod; label: string; icon: string }[] = [
  { id: "ughc", label: "KHGT", icon: "🌍" },
  { id: "umm_al_qura", label: "Umm al-Qura", icon: "🕋" },
  { id: "local_hisab", label: "Hisab", icon: "🔢" },
  { id: "local_rukyat", label: "Rukyat", icon: "🔭" },
];
