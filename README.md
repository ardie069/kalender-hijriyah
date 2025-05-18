# 📅 Kalender Hijriyah Sederhana 🌙

Kalender Hijriyah Sederhana adalah aplikasi web yang menampilkan **tanggal Hijriyah secara real-time** berdasarkan lokasi pengguna, perhitungan astronomis, dan metode yang dipilih. menggunakan frontend **Vue 3** dengan dukungan **Express.js** pada backend.

## ✨ Fitur Utama

- ✅ **Konversi Kalender Hijriyah** dari Masehi secara real-time
- ✅ **Pilihan Metode Perhitungan**: Global, Hisab, Rukyat
- ✅ **Tema Dark & Light Mode** yang dapat diganti
- ✅ **Penyesuaian Otomatis** tanggal Hijriyah setelah matahari terbenam
- ✅ **Jam Realtime** dengan zona waktu lokal
- ✅ **Deteksi Lokasi** otomatis melalui Geolocation API
- ✅ **Prediksi Akhir Bulan Hijriyah** yang di mana setiap metode perhitungan akan divalidasi dengan **Imkanur Rukyat**
- ✅ **Nama Weton** yang tersedia hari weton di kawasan Jawa dan sekitarnya, seperti _Senin Pahing_

## 🔧 Teknologi yang Digunakan

- 🌐 **Vue 3** + Vite (Frontend)
- 🛠️ **Express.js** (Backend)
- 🎨 **Tailwind CSS** + **Daisy UI** (Styling)
- 🧮 **Luxon**, **SunCalc**, **Astronomia** untuk kalkulasi Hijriyah dan posisi matahari

## 🚀 Cara Menjalankan

### 1. Clone repository ini

```sh
git clone https://github.com/ardie069/kalender-hijriyah.git
```

### 2. Masuk ke folder proyek

```sh
cd kalender-hijriyah
```

### 3. Instalasi dependensi

```sh
npm install
```

### 4. Jalankan backend

```sh
node server/server.js
```

### 5. Jalankan frontend (development)

```sh
cd frontend
npm run dev
```

### 6. Akses di browser

```sh
http://localhost:5173
```

> 🧪 Untuk produksi, frontend akan dibuild dan disajikan oleh backend dari folder `dist`.

## 📁 Struktur Proyek

```plaintext
kalender-hijriyah/
├── apps/       
│   ├── backend/                # Backend Express
│   │    ├── server.js
│   │    └── hijriCalculator.js
│   │ 
│   └── frontend/               # Frontend Vue
│       ├── index.html
│       ├── vite.config.js
│       ├── src/
│       ├── App.vue
│       ├── main.js
│       ├── components/
│       │    └── HijriDate.vue
│       └── .env     
└── cli/                        # mode cli atau terminal
```

## 📝 Lisensi

MIT License. Silakan digunakan dan dikembangkan lebih lanjut! 🚀

---

Dikembangkan oleh [@ardie069](https://github.com/ardie069) 💫
