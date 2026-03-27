# 🕌 Kalender Hijriyah Digital API Berbasis Golang & NASA SPICE 🌙

[![Go Release](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Framework: Gin](https://img.shields.io/badge/Framework-Gin-059669?style=flat-square&logo=gin)](https://gin-gonic.com/)
[![Engine: SPICE Toolkit](https://img.shields.io/badge/Engine-SPICE%20C--Kernel-10b981?style=flat-square)](https://naif.jpl.nasa.gov/naif/toolkit.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

Aplikasi ini berfokus murni pada penyediaan REST API untuk perhitungan **tanggal Hijriyah secara akurat dan konsisten** berdasarkan lokasi geografis, posisi matahari terbenam (Maghrib), dan berbagai metode penetapan global maupun lokal menggunakan standar data ephemeris **NASA JPL SPICE**.

---

## 🧠 Perencanaan

Aplikasi ini dibangun menggunakan **Golang** untuk memaksimalkan performa dan konkurensi, dengan integrasi langsung ke pustaka **NASA SPICE (CSPICE)** melalui CGO. Hal ini memungkinkan perhitungan astronomi dilakukan dengan standar ketelitian tingkat wahana antariksa.

1. Data astronomi dikalkulasi dari ephemeris riil NASA JPL `de440s.bsp` untuk tingkat ketelitian luar angkasa.
2. Menjembatani perbedaan metode penetapan Hijriyah (KHGT, Umm al-Qura, MABIMS, Wujudul Hilal) di dalam satu platform.
3. Memastikan pergantian hari Hijriyah **tepat saat Matahari terbenam** (Maghrib lokal pengamat) dengan pencarian *bisection* presisi.

---

## ✨ Fitur Utama

- 📡 **Real-time Lunar Telemetry**: Sinkronisasi dan perhitungan lintasan matahari dan bulan dengan CSPICE.
- 🌅 **Precise Solar Events**: Deteksi *Maghrib* (Sunset) dan *Subuh* (Fajr) menggunakan metode *bisection search* 24 jam untuk akurasi tinggi di berbagai lintang.
- 🕋 **Prayer Times Engine**: Perhitungan jadwal shalat multi-metode (Kemenag, MWL, ISNA, Umm Al-Qura, dll) dengan dukungan koreksi lintang tinggi.
- 🧪 **Multi-Method Engine**:
  - 🌍 **KHGT**: Kalender Hijriyah Global Tunggal (Turki 2016) dengan pemindaian visibilitas global.
  - 🕋 **Umm al-Qura**: Standar kalender resmi Arab Saudi (Makkah).
  - 🔢 **Hisab**: Berdasarkan Wujudul Hilal.
  - 🔭 **Rukyatul Hilal**: Berdasarkan kriteria MABIMS (2022).
- 🔮 **Global Scan Optimizations**: Scan visibilitas komprehensif dari barat (Benua Amerika) untuk hisab global dengan koreksi ekuinoks/midnight UTC.

---

## 🔧 Teknologi & Arsitektur

### Backend API Server

- **Golang**: Bahasa yang sangat cepat dengan concurrency yang handal.
- **Gin**: HTTP Web framework yang simpel dan cepat.
- **CGO & CSPICE**: Menggunakan `spice_bridge.go` untuk mengakses fungsi CSPICE secara aman dengan mekanisme `sync.Mutex` karena keterbatasan *thread-safety* pada pustaka C asli.
- **Docker**: Siap deploy dimana saja secara konsisten dengan multi-stage build.

---

## 🚀 Deployment & Build

API v4 ini dapat dijalankan menggunakan Docker untuk kemudahan dan keandalan pustaka CSPICE *shared object*.

```bash
git clone https://github.com/ardie069/kalender-hijriyah.git
cd kalender-hijriyah

# Build & Run dengan Docker Compose
docker-compose up --build -d
```

API Server akan menyala secara lokal pada `http://localhost:8080`.

## 📁 Struktur Proyek

```plaintext
kalender-hijriyah/
├── cmd/
│   └── api/              # Entry point Gin API Server
├── data/                 # SPICE Kernels (de440s.bsp, naif0012.tls, dll)
├── core/
│   ├── api/              # HTTP Handlers
│   ├── astronomy/        # SPICE CGO Bindings & Kalkulator Orbit
│   ├── calendar/         # Logika KHGT, Umm al-Qura, MABIMS
│   ├── models/           # Skema logika Golang
│   ├── prayer/           # Konfigurasi & Kalkulasi Jadwal Shalat
│   └── services/         # Orkestrasi Kalender Utama
├── docs/                 # Dokumentasi Detail & Teori Falak
└── README.md
```

🌐 API Endpoints Utama

| Endpoint                 | Method | Fungsi                                      | Parameter Query                    |
| ------------------------ | ------ | ------------------------------------------- | ---------------------------------- |
| `/api/v4/hijri/date`     | GET    | Konversi tanggal ke Hijriyah (Multi-metode) | lat, lon, method, date             |
| `/api/v4/prayer/times`   | GET    | Jadwal Salat harian presisi                 | lat, lon, method, asr_method, date |
| `/api/v4/moon/telemetry` | GET    | Data astronomi Bulan real-time              | lat, lon, date                     |

---

## 📜 Lisensi

Lisensi MIT. Bebas dikembangkan untuk ilmu pengetahuan dan kemaslahatan umat.

Dikembangkan oleh [Ardiansyah](https://github.com/ardie069).
