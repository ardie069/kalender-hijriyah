# 🕌 Kalender Hijriyah Digital API Berbasis Golang & NASA SPICE 🌙

[![Go Release](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Framework: Gin](https://img.shields.io/badge/Framework-Gin-009688?style=flat-square&logo=gin)](https://gin-gonic.com/)
[![Engine: SPICE Toolkit](https://img.shields.io/badge/Engine-SPICE%20C--Kernel-10b981?style=flat-square)](https://naif.jpl.nasa.gov/naif/toolkit.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

Aplikasi ini berfokus murni pada penyediaan REST API untuk perhitungan **tanggal Hijriyah secara akurat dan konsisten** berdasarkan lokasi geografis, posisi matahari terbenam (Maghrib), dan berbagai metode penetapan global maupun lokal menggunakan standar data ephemeris **NASA JPL SPICE**.

---

## 🧠 Perencanaan

Proyek ini telah dire-write dari Python/FastAPI ke **Golang** untuk memaksimalkan performa dan konkurensi, menggunakan binding ke pustaka **NASA SPICE (CSPICE)**.

1. Data astronomi dikalkulasi dari ephemeris riil NASA JPL `de440s.bsp` untuk tingkat ketelitian luar angkasa.
2. Menjembatani perbedaan metode penetapan Hijriyah (KHGT, Umm al-Qura, MABIMS, Wujudul Hilal) di dalam satu platform.
3. Memastikan pergantian hari Hijriyah **tepat saat Matahari terbenam** (Maghrib lokal pengamat).

---

## ✨ Fitur Utama

- 📡 **Real-time Lunar Telemetry**: Sinkronisasi dan perhitungan lintasan matahari dan bulan dengan CSPICE.
- 🌅 **Precise Sunset/Fajr Bisection**: Bisection search presisi tinggi 24 jam untuk mendeteksi maghrib dan subuh lintas batasan zona waktu.
- 🧪 **Multi-Method Engine**:
  - 🌍 **KHGT**: Kalender Hijriyah Global Tunggal (Turki 2016).
  - 🕋 **Umm al-Qura**: Standar otoritas Makkah.
  - 🔢 **Hisab**: Berdasarkan Wujudul Hilal.
  - 🔭 **Rukyatul Hilal**: Berdasarkan kriteria MABIMS (2022).
- 🔮 **Global Scan Optimizations**: Scan visibilitas komprehensif dari barat (Benua Amerika) untuk hisab global dengan koreksi ekuinoks/midnight UTC.

---

## 🔧 Teknologi & Arsitektur

### Backend API Server

- **Golang**: Bahasa yang sangat cepat dengan concurrency yang handal.
- **Gin**: HTTP Web framework yang simpel dan cepat.
- **CGO & CSPICE**: Toolkit geometri navigasi antariksa standar NASA yang ditautkan via CGO.
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
├── internal/
│   ├── api/              # HTTP Handlers
│   ├── astronomy/        # SPICE CGO Bindings & Kalkulator Orbit
│   ├── calendar/         # Logika KHGT, Umm al-Qura, MABIMS
│   ├── models/           # Skema logika Golang
│   └── services/         # Orkestrasi Kalender Utama
├── docs/                 # Dokumentasi Metode Hijriyah
└── README.md
```

🌐 API Endpoints Utama

| Endpoint                 | Method | Fungsi                                      | Parameter Query                                                                 |
| ------------------------ | ------ | ------------------------------------------- | ------------------------------------------------------------------------------- |
| `/api/v4/hijri/date`     | GET    | Perhitungan prediksi bulan                  | date (ISO), lat, lon                                                            |
| `/api/v4/moon/telemetry` | GET    | Informasi bulan berdasarkan lokasi pengguna | lat, lon, altitude, elongation, azimuth, distance km, illumination, phase, time |

---

## 📜 Lisensi

Lisensi MIT. Bebas dikembangkan untuk ilmu pengetahuan dan kemaslahatan umat.

Dikembangkan oleh [Ardiansyah](https://github.com/ardie069).
