# 🕌 Hilal Scope: Kalender Hijriyah Digital API 🌙

[![Go Release](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Framework: Gin](https://img.shields.io/badge/Framework-Gin-059669?style=flat-square&logo=gin)](https://gin-gonic.com/)
[![Engine: SPICE Toolkit](https://img.shields.io/badge/Engine-SPICE%20C--Kernel-10b981?style=flat-square)](https://naif.jpl.nasa.gov/naic/toolkit.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

**Hilal Scope** adalah REST API berperforma tinggi untuk perhitungan kalender Hijriyah presisi tingkat wahana antariksa. Menggunakan data ephemeris **NASA JPL SPICE**, API ini menyediakan sinkronisasi real-time antara posisi benda langit dengan penetapan tanggal Hijriyah menggunakan berbagai kriteria global maupun lokal.

---

## 🚀 Fitur Unggulan

- 🔭 **NASA-Grade Precision**: Menggunakan kernel `de440s.bsp` untuk akurasi data posisi Bulan dan Matahari yang ekstrem.
- 🌍 **KHGT (Unified Global Hijri Calendar)**: Implementasi Kalender Hijriyah Global Tunggal dengan optimasi pemindaian visibilitas global (Scan 24 jam).
- 🇮🇩 **MABIMS (Sabang Reference)**: Penentuan awal bulan berdasarkan titik geografis paling barat Indonesia (Sabang) sesuai standar Kemenag 2022.
- 📊 **Rich Lunar Telemetry**: Data altitude, elongation, illumination, dan fase bulan secara real-time maupun prediktif.
- 🕋 **Multi-Location Hijri**: Mendukung perhitungan toposentris (lokal observer) dan geosentris.
- ⚡ **Serverless Ready**: Teroptimasi untuk dideploy di **Vercel** dengan penanganan khusus pustaka CGO/CSPICE.

---

## 🛠️ Teknologi & Arsitektur

- **Golang**: Bahasa utama untuk performa tinggi dan konkurensi.
- **Gin**: Framework HTTP yang efisien.
- **CSPICE (C-Kernel)**: Integrasi CGO ke pustaka NASA SPICE untuk kalkulasi falak tingkat lanjut. [Sumber CSPICE Toolkit](https://naif.jpl.nasa.gov/pub/naif/toolkit/)
- **Docker**: Mendukung multi-stage build untuk deployment yang bersih.

---

## 📁 Struktur Proyek

```plaintext
kalender-hijriyah/
├── cmd/
│   └── api/              # Entry point Gin API Server
├── data/                 # SPICE Kernels (de440s.bsp, naif0012.tls, dll)
├── core/
│   ├── api/              # Handlers & Routes
│   ├── astronomy/        # Bindings CSPICE & Orbit Engine
│   ├── calendar/         # Logika KHGT, MABIMS (Sabang), Umm al-Qura
│   ├── models/           # Skema Data (HijriDate, MoonTelemetry)
│   ├── prayer/           # Jadwal Shalat & Koreksi Lintang
│   └── services/         # Orkestrasi & Resolusi Bulan
├── docs/                 # Dokumentasi Teknis & Teori Falak
└── README.md
```

---

## 📡 API Endpoints (v4)

Semua endpoint tersedia di prefix `/api/v4/` atau `/v4/`.

| Endpoint | Method | Deskripsi | Parameter Utama |
|----------|--------|-----------|-----------------|
| `/hijri/date` | GET | Konversi & Prediksi Hijriyah (4 Metode) | `lat`, `lon`, `date` |
| `/hijri/calendar` | GET | Kalender Hijriyah Tahunan | `year`, `lat`, `lon` |
| `/moon/telemetry` | GET | Data Astronomi Bulan Real-time | `lat`, `lon` |
| `/prayer/times` | GET | Jadwal Shalat Presisi | `lat`, `lon`, `method` |

---

## 🚢 Deployment (Vercel)

API Hilal Scope dapat diakses secara publik melalui:  
**[https://kalender-hijriyah-api.vercel.app/](https://kalender-hijriyah-api.vercel.app/)**

Proyek ini telah dikonfigurasi untuk berjalan di Vercel dengan skrip build otomatis yang mengunduh dan mengonfigurasi `libcspice.a`.

1. Pastikan file `vercel.json` sudah ada di root.
2. Gunakan Environment Variables untuk konfigurasi tambahan jika diperlukan.
3. Push ke GitHub dan hubungkan dengan Vercel.

---

## 📚 Dokumentasi Lanjutan

Untuk selengkapnya mengenai API ini, silahkan kunjungi [Dokumentasi Lengkap](docs/README.md)

---

## 📜 Lisensi

Lisensi MIT. Dikembangkan oleh [Ardiansyah](https://github.com/ardie069).

## ⚠️ Disclaimer & Credits

- *Astronomical calculations are powered by NASA's SPICE Toolkit (NAIF).*
- *This project is an independent educational tool and is not affiliated with or endorsed by NASA.*
- *All calculations are for informational purposes only and should not be used for critical decision-making.*
