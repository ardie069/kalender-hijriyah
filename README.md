# 🕌 Kalender Hijriyah Digital — Lunar Analytics 🌙

[![Framework: Next.js 15](https://img.shields.io/badge/Frontend-Next.js%2015-000000?style=flat-square&logo=nextdotjs)](https://nextjs.org/)
[![Backend: FastAPI](https://img.shields.io/badge/Backend-FastAPI-009688?style=flat-square&logo=fastapi)](https://fastapi.tiangolo.com/)
[![Engine: Skyfield](https://img.shields.io/badge/Engine-Skyfield%20API-10b981?style=flat-square)](https://rhodesmill.org/skyfield/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

> **"Sains adalah obor yang membakar kegelapan ghaib dengan cahaya logika materialistik."**
> Aplikasi ini adalah sintesis antara presisi astronomi modern dan kearifan lokal dalam satu ekosistem digital.

Aplikasi ini menampilkan **tanggal Hijriyah secara akurat dan konsisten** berdasarkan lokasi geografis, posisi matahari terbenam (Maghrib), dan berbagai metode penetapan global maupun lokal. Kami memisahkan secara tegas antara **Logika Falak (Backend)** dan **Estetika User Experience (Frontend)** untuk menjaga validitas data.

---

## 🧠 Filosofi Madilog & Dialektika Falak

Proyek ini dibangun di atas prinsip **Materialisme, Dialektika, dan Logika**:

1. **Materialisme**: Data astronomi didasarkan pada pergerakan benda langit riil menggunakan ephemeris NASA JPL DE421.
2. **Dialektika**: Menjembatani perbedaan metode (Hisab vs Rukyat) melalui transparansi data visibilitas hilal.
3. **Logika**: Memastikan tidak ada perubahan tanggal di tengah malam; hari baru dimulai tepat saat Matahari terbenam.

---

## ✨ Fitur Utama (v2.0 - 2026 Edition)

- 📡 **Real-time Hijri Telemetry**: Sinkronisasi otomatis berdasarkan koordinat GPS dan zona waktu pengguna.
- 🌅 **Sunset-Based Transition**: Logika pergantian hari saat Maghrib, bukan pukul 00:00.
- 🧪 **Multi-Method Engine**:
  - 🌍 **KHGT (Global)**: Kalender Hijriyah Global Tunggal.
  - 🕋 **Umm al-Qura**: Standar otoritas Mekkah.
  - 🔢 **Hisab Lokal**: Berdasarkan kriteria Wujudul Hilal.
  - 🔭 **Rukyat Lokal**: Simulasi visibilitas hilal berbasis kriteria MABIMS.
- 🔮 **Predictive Analytics**: Prediksi apakah bulan berjalan berjumlah 29 atau 30 hari (Istikmal).
- 🗺️ **Global Hilal Map**: Integrasi pelacakan posisi hilal di seluruh belahan dunia.
- ✨ **Local Wisdom**: Integrasi Nama Weton (khusus wilayah Jawa).
- 🎨 **Emerald-Scientific UI**: Desain "Lunar Aura" yang mendukung Dark/Light mode secara *seamless*.

---

## 🔧 Teknologi & Arsitektur

### Backend (The Engine)

- **FastAPI**: *High-performance* Python framework untuk pemrosesan data asinkron.
- **Skyfield**: Library astronomi untuk kalkulasi posisi planet dan satelit dengan tingkat akurasi tinggi.
- **Julian Day Calculation**: Standar internal perhitungan waktu astronomis.
- **Vercel Serverless**: Infrastruktur backend yang *scalable* dan enteng.

### Frontend (The Dashboard)

- **Next.js 15**: Menggunakan *App Router* dan *Server Components* untuk performa maksimal.
- **Tailwind CSS v4 + DaisyUI**: Estetika bento-grid dengan sentuhan *Glassmorphism*.
- **TypeScript**: Menjamin keamanan tipe data antara API dan UI.
- **Zustand/Context**: Manajemen state tema dan preferensi pengguna.

---

## 🚀 Instalasi & Development

### 1. Clone & Setup

```bash
git clone [https://github.com/ardie069/kalender-hijriyah.git](https://github.com/ardie069/kalender-hijriyah.git)
cd kalender-hijriyah
```

### 2. Backend (FastAPI)

```bash
cd apps/backend
python -m venv .venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
pip install -r requirements.txt
uvicorn app.main:app --reload
```

API akan tersedia di `http://127.0.0.1:8000`.

### 3. Frontend (Next.js)

```bash
cd apps/web
pnpm install
pnpm dev
```

Dashboard akan tersedia di `http://localhost:3000`.

## 📁 Struktur Proyek (Atomic Design)

```plaintext
kalender-hijriyah/
├── apps/
│   ├── backend/          # Python FastAPI (Falak Logic)
│   │   ├── app/          # Core API & Astronomical Service
│   │   └── data/         # Ephemeris Files (de421.bsp)
│   └── web/              # Next.js 15 (Emerald UI)
│       ├── components/   # Atomic UI Components
│       ├── context/      # Theme & Auth Context
│       └── hooks/        # Custom React Hooks (useHijri)
├── docs/                 # Documentation & Research
└── README.md
```

🌐 API Endpoints Utama

| Endpoint             | Fungsi                       | Parameter Utama     |
| -------------------- | ---------------------------- | --------------------|
| GET /hijri-date      |	Ambil tanggal Hijriyah riil |	method, lat, lon    |
| GET /hijri-end-month |	Prediksi akhir bulan        |	method, year, month |

---

## 📜 Lisensi

Lisensi MIT. Bebas dikembangkan untuk kemaslahatan umat dan ilmu pengetahuan.

Build with Logic & Passion by [Ardiansyah](https://github.com/ardie069). Fokus: Falakiyah, Konsistensi Kalender, dan Dialektika Logika.
