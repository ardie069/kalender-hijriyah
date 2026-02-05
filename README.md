# 📅 Kalender Hijriyah 🌙

Kalender Hijriyah adalah aplikasi web untuk menampilkan **tanggal Hijriyah secara akurat dan konsisten**
berdasarkan **lokasi geografis**, **waktu Matahari terbenam (Maghrib)**, dan **metode penetapan bulan Hijriyah**
(Global, Hisab, dan Rukyat).

Aplikasi ini memisahkan **logika falak (backend)** dan **tampilan (frontend)** secara tegas
untuk menghindari inkonsistensi perhitungan.

---

## ✨ Fitur Utama

- ✅ **Tanggal Hijriyah real-time** berbasis lokasi dan zona waktu
- ✅ **Pergantian hari Hijriyah saat Maghrib**, bukan tengah malam
- ✅ **Metode Perhitungan**:
  - 🌍 Global (Umm al-Qura)
  - 🧮 Hisab Astronomis atau Wujudul Hilal
  - 🌙 Rukyat Hilal (Imkanur Rukyat)
- ✅ **Prediksi Akhir Bulan Hijriyah** (29 atau 30 hari)
- ✅ **Validasi visibilitas hilal** (usia bulan, ketinggian, elongasi)
- ✅ **Penjelasan keputusan kalender (explain endpoint)**
- ✅ **Nama Weton** (khusus wilayah Jawa)
- ✅ **Dark / Light Mode**
- ✅ **Jam real-time sesuai zona waktu pengguna**

---

## 🧠 Prinsip Perhitungan

- Hari Hijriyah **dimulai saat Matahari terbenam (Maghrib)**
- Hisab dan Rukyat **dipisahkan secara logika**
- Tidak ada loncatan tanggal (29 → 1 → 2)
- Tidak ada perubahan tanggal di tengah malam
- Rukyat **lebih ketat** daripada Hisab

> Fokus utama proyek ini adalah **konsistensi kalender**, bukan sekadar konversi tanggal.

---

## 🔧 Teknologi yang Digunakan

### Backend (API)

- ⚡ **FastAPI (Python)**
- 🌌 **Skyfield** + Ephemeris JPL (DE421)
- 🧮 Astronomi berbasis Julian Day
- 🕰️ pytz & timezone-aware datetime
- 🚦 Rate limiting (SlowAPI)
- ☁️ Deploy: **Vercel (Serverless Python)**

### Frontend

- 🌐 **Next.js + Tailwind CSS**
- 🟦 **TypeScript**
- 🎨 **Tailwind CSS + DaisyUI**
- ☁️ Deploy: **Vercel**

---

## 🚀 Cara Menjalankan (Development)

### 1️⃣ Clone repository

```sh
git clone https://github.com/ardie069/kalender-hijriyah.git
cd kalender-hijriyah
```

---

2️⃣ Jalankan Backend (FastAPI)

```sh
cd apps/backend
python -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install .
uvicorn app.main:app --reload
```

Backend akan berjalan di:

```sh
http://127.0.0.1:8000
```

---

3️⃣ Jalankan Frontend (Next.js)

```sh
cd apps/frontend
npm install
npm run dev
```

Frontend akan berjalan di:

```sh
http://localhost:3000
```

---

🌐 Endpoint Utama API

| Endpoint | Keterangan |
| -------- | ---------- |
| /hijri-date | Tanggal Hijriyah hari ini |
| /hijri-end-month | Prediksi akhir bulan |
| /hijri-explain | Penjelasan keputusan kalender |
| /health | Health check |

---

📁 Struktur Proyek

```plaintext
kalender-hijriyah/
├── apps/
│   ├── backend/          # FastAPI + Astronomi
│   │   ├── app/
│   │   ├── data/         # de421.bsp
│   │   └── pyproject.toml
│   │
│   └── web/         # Next.js + Tailwind CSS
│       ├── app/
│       ├── public/
│       └── next.config.ts
│
├── docs/                 # Dokumentasi falak & arsitektur
├── netlify.toml
└── README.md

```

---

📜 Lisensi

MIT License — bebas digunakan, dimodifikasi, dan dikembangkan.

---

Dikembangkan oleh [Ardiansyah](https://github.com/ardie069)
🧠 Fokus: falak, konsistensi kalender, dan logika yang bisa dipertanggungjawabkan.

---

<details>
<summary><strong>📘 English Documentation (Click to expand)</strong></summary>

## 📅 Hijri Calendar 🌙

Hijri Calendar is a web application designed to provide **accurate and consistent Hijri dates**
based on **geographical location**, **local sunset (Maghrib)**, and **Hijri month determination methods**
(Global, Hisab, and Rukyat).

The project strictly separates **astronomical calculation logic (backend)**
from **presentation logic (frontend)** to avoid calendar inconsistencies.

---

## ✨ Key Features

- ✅ **Real-time Hijri date** based on location and timezone
- ✅ **Hijri day transition at sunset**, not at midnight
- ✅ **Calculation methods**:
  - 🌍 Global (Umm al-Qura)
  - 🧮 Astronomical Hisab
  - 🌙 Hilal Observation (Rukyat)
- ✅ **End-of-month prediction** (29 or 30 days)
- ✅ **Hilal visibility validation** (moon age, altitude, elongation)
- ✅ **Decision explanation endpoint**
- ✅ **Javanese Weton support** (regional)
- ✅ **Dark / Light mode**
- ✅ **Real-time clock**

---

## 🧠 Calculation Principles

- Hijri day **starts at sunset (Maghrib)**
- Hisab and Rukyat are **logically separated**
- No date skipping (29 → 1 → 2)
- No date change at midnight
- Rukyat criteria are **stricter** than Hisab

> The main goal of this project is **calendar consistency**, not mere date conversion.

---

## 🔧 Technology Stack

### Backend (API)

- ⚡ FastAPI (Python)
- 🌌 Skyfield + JPL Ephemeris (DE421)
- 🧮 Julian Day based astronomy
- 🕰️ Timezone-aware datetime
- 🚦 Rate limiting (SlowAPI)
- ☁️ Deployment: Vercel (Serverless)

### Frontend

- 🌐 Next.js
- 🟦 TypeScript
- 🎨 Tailwind CSS + DaisyUI
- ☁️ Deployment: Vercel

---

## 🌐 Main API Endpoints

| Endpoint | Description |
| -------- | ------------- |
| `/hijri-date` | Current Hijri date |
| `/hijri-end-month` | End-of-month prediction |
| `/hijri-explain` | Decision explanation |
| `/health` | Health check |

---

## 📜 License

MIT License — free to use and modify.

---

Developed by **Ardi**  
Focus: astronomical correctness, calendar consistency, and accountable logic.

</details>
