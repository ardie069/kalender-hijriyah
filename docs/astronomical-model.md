# 🌌 Astronomical Model Hilal Scope (v4)

Hilal Scope menggunakan model astronomi *real-time* berbasis data ephemeris standar tinggi untuk menghitung posisi Matahari dan Bulan menggunakan pustaka **NASA SPICE**.

## Komponen Engine Utama

- Algoritma Dasar J2000 (Julian Date Reference)
- Ephemeris JPL **DE440s** (`de440s.bsp`)
- Planetary Constants (`pck00011.tpc`) dan Leapseconds (`naif0012.tls`)
- **Golang CSPICE Bindings** (via CGO)

## Perhitungan Kunci (Falak Engine)

- **Global Scan KHGT Check**: Bisection Search untuk mendeteksi Sunset & Fajr dalam lingkup *sliding window* 24-jam di seluruh dunia (batas ±65° lintang).
- **Topocentric AltAz**: Proyeksi vektor Ecliptic ke geometri Zenith observer lokal, digunakan untuk kriteria **MABIMS** (dengan referensi Sabang) dan **Umm al-Qura**.
- **Geocentric Altitude**: Menggunakan frame Body-Fixed (IAU_EARTH) untuk kebutuhan hisab global (**KHGT**).
- **Topocentric Elongation**: Sejak v4.1.0, perhitungan elongasi (jarak sudut Bulan-Matahari) diubah dari **Geosentris** ke **Toposentris**. Hal ini memberikan data yang lebih akurat sesuai dengan apa yang benar-benar dilihat oleh pengamat di permukaan Bumi, terutama untuk kriteria visibilitas hilal yang sangat sensitif.
- **Refraksi Atmosfer**: Koreksi refraksi (model Bennett/Sæmundsson) diterapkan pada kriteria toposentris untuk akurasi visual yang lebih baik.

## 🛠 Arsitektur Teknis: SPICE & CGO

Sistem ini mengintegrasikan pustaka **CSPICE** (versi C dari toolkit SPICE NASA) ke dalam ekosistem **Golang** menggunakan CGO.

### 1. Keamanan Thread (Thread-Safety)

Pustaka CSPICE asli **tidak bersifat thread-safe**. Karena menggunakan variabel statis internal, Hilal Scope menerapkan **Global Lock** menggunakan `sync.Mutex` (`spiceMu`) di setiap pemanggilan fungsi jembatan CGO. Ini memastikan stabilitas engine pada lingkungan server yang melayani banyak request secara konkuren.

### 2. Deployment Serverless (Vercel)

Tantangan utama menjalankan engine ini di Vercel adalah ketergantungan pada pustaka C dan file kernel yang besar.

- **Vercel Build Step**: Kami menggunakan skrip `vercel-build.sh` untuk mengunduh binary `libcspice.a` yang sudah dikompilasi secara statis untuk lingkungan Amazon Linux 2 (arsitektur Vercel).
- **Static Linking**: Library di-link secara statis saat proses build Go di Vercel, sehingga binary akhir bersifat mandiri (*standalone*).
- **Kernel Placement**: File-file kernel (`.bsp`, `.tls`, `.tpc`) diletakkan di folder `data/` yang di-copy ke output deployment untuk diakses saat runtime melalui `furnsh_c`.

## 🏗 Setup Lokal

Untuk pengembangan di mesin lokal (Linux/WSL):

1. Pastikan GCC atau Clang terpasang.
2. Letakkan file library di `core/astronomy/lib/` dan header di `core/astronomy/include/`.
3. Set `CGO_ENABLED=1` saat menjalankan `go build`.

```bash
CGO_ENABLED=1 go run cmd/api/main.go
```

## Catatan Penting

- **Akurasi**: Data posisi memiliki tingkat ketelitian mili-arcsecond, namun hasil akhir sangat bergantung pada kriteria fikih yang dipilih.
- **Limitasi**: Pada lintang ekstrem (di atas ±65°), perhitungan Sunset/Fajr menggunakan metode koreksi khusus untuk menangani fenomena matahari tengah malam atau malam kutub.
- **Visibilitas**: Model ini menghitung kriteria visibilitas spesifik (threshold), bukan simulasi penampakan visual optik (teleskopik).
