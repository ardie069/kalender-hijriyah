# Changelog

Semua perubahan penting pada proyek ini akan didokumentasikan di file ini.

Format changelog ini berdasarkan [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), dan proyek ini mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.1.0] - 2026-03-28

### ✨ Fitur Baru (Added)
- Mendukung telemetri spesifik Makkah untuk metode **Umm al-Qura**.
- Penambahan data telemetri bulan untuk lokasi referensi (Topocentric).
- Sinkronisasi data fase bulan dengan akurasi toposentris dalam API response.

### 🚀 Peningkatan & Refactor (Improved/Refactored)
- Perubahan kalkulasi perpanjangan (**Elongation**) dari koordinat **Geosentris** ke **Toposentris** untuk akurasi pengamat lokal yang lebih tinggi.
- Penentuan titik visibilitas hilal pada metode **KHGT** sekarang lebih deterministik.
- Pembersihan fungsi `evaluateGeocentricPoint` yang tidak lagi digunakan.
- Optimasi struktur data telemetri bulan untuk konsistensi antar endpoint.

### 🐛 Perbaikan (Fixed)
- Perbaikan ambiguitas pemilihan titik visibilitas pada pemindaian global KHGT.
- Koreksi penamaan fase bulan agar sesuai dengan standar terminologi astronomi.

### 📚 Dokumentasi (Docs)
- Update terminologi **KHGT** pada `README.md` agar lebih jelas bagi pengguna baru.
- Penambahan dokumentasi desain API (v4) dengan struktur response yang lebih detail.
- Klarifikasi kebutuhan parameter koordinat (`lat`, `lon`) sebagai input wajib.

---

## [4.0.0] - 2026-03-25
*Rilis awal dengan integrasi NASA JPL SPICE Kernel.*
