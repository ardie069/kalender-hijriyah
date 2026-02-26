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

export const HIJRI_MONTHS_INTERNATIONAL_GRAMMAR = [
  { id: 1, name: "Muharram", desc: "Start of Year" },
  { id: 2, name: "Shafar", desc: "Travel Month" },
  { id: 3, name: "Rabiul Awwal", desc: "Birth of the Prophet" },
  { id: 4, name: "Rabiul Akhir", desc: "Winter Solstice" },
  { id: 5, name: "Jumadil Ula", desc: "Start of Dry Season" },
  { id: 6, name: "Jumadil Akhir", desc: "End of Dry Season" },
  { id: 7, name: "Rajab", desc: "Holy Month" },
  { id: 8, name: "Sha'ban", desc: "Preparation" },
  { id: 9, name: "Ramadan", desc: "Obligatory Fasting", isSpecial: true },
  { id: 10, name: "Syawwal", desc: "Eid al-Fitr", isSpecial: true },
  { id: 11, name: "Dhul Qa'dah", desc: "Calmiring Month" },
  { id: 12, name: "Dhul Hijjah", desc: "Hajj Month", isSpecial: true },
];

export type HijriMonth = (typeof HIJRI_MONTHS_INDONESIA_GRAMMAR)[number];
export type HijriMonthInternational = (typeof HIJRI_MONTHS_INTERNATIONAL_GRAMMAR)[number];

export const METHODS: { id: HijriMethod; label: string; icon: string }[] = [
  { id: "ughc", label: "KHGT", icon: "🌍" },
  { id: "umm_al_qura", label: "Umm al-Qura", icon: "🕋" },
  { id: "local_hisab", label: "Hisab", icon: "🔢" },
  { id: "local_rukyat", label: "Rukyat", icon: "🔭" },
];
