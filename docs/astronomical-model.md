# 🌌 Astronomical Model v4

Sistem ini menggunakan model astronomi *real-time* berbasis data ephemeris standar tinggi
untuk menghitung posisi Matahari dan Bulan menggunakan binding ke pustaka NASA SPICE.

## Komponen Engine Utama

- Algoritma Dasar J2000 (Julian Date Reference)
- Ephemeris JPL **DE440s** (de440s.bsp)
- Planetary Constants (pck00011.tpc) dan Leapseconds (naif0012.tls)
- **Golang CSPICE Bindings** (via CGO)

## Perhitungan Kunci (Falak Engine)

- **Global Scan KHGT Check**: Bisection Search untuk mendeteksi terbenamnya matahari (Sunset) & fajar (Fajr) dalam lingkup *sliding window* 24-jam di lintang dan bujur belahan dunia manapun (65° lintang utara/selatan).
- **Topocentric AltAz**: Proyeksi vektor Ecliptic (ECLIPJ2000) ke geometri Zenith observer lokal, digunakan untuk kriteria MABIMS dan Umm al-Qura.
- **Geocentric Altitude**: Menggunakan frame Body-Fixed (IAU_EARTH) untuk kebutuhan hisab murni (seperti KHGT dan Wujudul Hilal).
- Analisis Elongasi Geosentrik (Sudut pisah pusat Bumi-Bulan-Matahari) dan Toposentrik.
- Prediksi Konjungsi Presisi (Ijtima).

## 🌦 Refraksi Atmosfer

Untuk kriteria berbasis Toposentrik (seperti MABIMS), sistem menerapkan koreksi refraksi atmosfer pada fungsi `ApplyRefraction`.
Kami menggunakan model aproksimasi Bennett/Sæmundsson yang menghitung deviasi cahaya berdasarkan ketinggian geometris hilal. Hal ini sangat krusial karena hilal yang secara geometris berada di bawah ufuk bisa jadi terlihat secara visual karena pembiasan atmosfer.

## 🛠 Arsitektur Teknis: SPICE, CGO, & Thread-Safety

Sistem ini mengintegrasikan pustaka **CSPICE** (versi C dari toolkit SPICE NASA) ke dalam ekosistem **Golang** untuk mendapatkan akurasi tingkat tinggi yang setara dengan navigasi wahana antariksa.

### 1. Integrasi CGO

Karena CSPICE adalah pustaka berbasis C, kami menggunakan **CGO** sebagai jembatan. File `spice_bridge.go` bertindak sebagai *wrapper* yang memetakan fungsi-fungsi C (seperti `spkpos_c`, `str2et_c`, `furnsh_c`) ke fungsi Go. Hal ini memungkinkan logika bisnis kalender tetap berada di level tinggi (Go) sementara perhitungan geometri astronomi yang kompleks dilakukan oleh mesin SPICE yang sudah teruji.

### 2. Keterbatasan Thread-Safety

Pustaka CSPICE asli **tidak bersifat thread-safe**. Pustaka ini menggunakan variabel statis dan global internal untuk manajemen state (seperti error handling dan pool kernel). Jika dua goroutine mencoba memanggil fungsi SPICE secara bersamaan, hal ini dapat menyebabkan *memory corruption*, hasil kalkulasi yang tidak konsisten, atau aplikasi berhenti mendadak (*crash*).

### 3. Mekanisme `sync.Mutex`

Untuk mengatasi masalah keamanan thread tersebut, aplikasi ini menerapkan kebijakan **Global Lock** menggunakan `sync.Mutex` (didefinisikan sebagai `spiceMu`).

Setiap kali fungsi jembatan astronomi dipanggil:

1. `spiceMu.Lock()` dijalankan untuk memastikan hanya satu goroutine yang memiliki akses ke engine SPICE pada satu waktu.
2. Fungsi internal SPICE dieksekusi melalui CGO.
3. `spiceMu.Unlock()` (biasanya melalui `defer`) dijalankan segera setelah operasi selesai untuk memberikan giliran kepada goroutine lain.

**Dampak Performa:**
Meskipun akses ke SPICE dilakukan secara berurutan (*serialized*), performa sistem tetap sangat tinggi. Hal ini dikarenakan:

- Operasi matematika SPICE sangat efisien dan berbasis memori.
- Bottleneck aplikasi biasanya berada pada I/O jaringan, bukan pada kecepatan CPU dalam menghitung orbit.
- Mekanisme ini jauh lebih aman dan stabil untuk lingkungan server yang melayani banyak permintaan secara konkuren.

## 🏗 Setup Environment CGO & CSPICE

Untuk mengompilasi proyek ini, sistem Anda harus memiliki compiler C (GCC/Clang) dan library CSPICE yang sudah terpasang di struktur folder yang benar.

### 1. Persyaratan Sistem

- **Linux/macOS**: GCC atau Clang terpasang.
- **Windows**: MSYS2 atau MinGW-w64.
- **Go**: Versi 1.25+ dengan `CGO_ENABLED=1`.

### 2. Download CSPICE Toolkit

Unduh toolkit sesuai dengan arsitektur sistem Anda dari NASA NAIF Website. [Klik di sini](https://naif.jpl.nasa.gov/pub/naif/toolkit/C/).

### 3. Struktur Direktori Library

Berdasarkan arahan `#cgo` di `core/astronomy/spice_bridge.go`, library harus diletakkan sebagai berikut:

```plaintext
core/astronomy/
├── include/           <-- Salin semua file .h dari direktori 'include' CSPICE
│   ├── SpiceUsr.h
│   ├── SpiceZpl.h
│   └── ...
├── lib/               <-- Salin file library (.a atau .so) dari direktori 'lib'
│   ├── cspice.a
│   └── csupport.a
└── spice_bridge.go
```

### 4. Kompilasi

Pastikan variabel lingkungan `CGO_ENABLED` aktif saat melakukan build:

```bash
# Linux / macOS
CGO_ENABLED=1 go build -o hilal-engine ./cmd/api

# Windows (PowerShell)
$env:CGO_ENABLED="1"; go build -o hilal-engine.exe ./cmd/api
```

*Catatan: Jika Anda menggunakan Docker, proses ini sudah diotomatisasi di dalam Dockerfile menggunakan multi-stage build untuk memastikan library C ter-link secara statis atau tersedia di runtime.*

## Catatan

- Sistem ini **tidak menggantikan observasi lapangan** (Rukyatul Hilal aktual).
- Hasil visibilitas adalah perhitungan matematis astronomis menggunakan limitasi kriteria visibilitas spesifik, bukan penampakan visual (Optik).
