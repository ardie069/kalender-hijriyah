# 🌌 Astronomical Model v4

Sistem ini menggunakan model astronomi *real-time* berbasis data ephemeris standar tinggi
untuk menghitung posisi Matahari dan Bulan menggunakan binding ke pustaka NASA SPICE.

## Komponen Engine Utama

- Algoritma Dasar J2000 (Julian Date Reference)
- Ephemeris JPL **DE440** (de440.bsp)
- Planetary Constants (pck00011.tpc) dan Leapseconds (naif0012.tls)
- **Golang CSPICE Bindings** (via CGO)

## Perhitungan Kunci (Falak Engine)

- **Global Scan UGHC Check**: Bisection Search untuk mendeteksi terbenamnya matahari (Sunset) & fajar (Fajr) dalam lingkup *sliding window* 24-jam di lintang dan bujur belahan dunia manapun.
- **Topocentric AltAz**: Proyeksi vektor Ecliptic (ECLIPJ2000) ke geometri Zenith observer lokal.
- **Geocentric Altitude**: Koreksi elevasi pengamat untuk kebutuhan hisab murni (seperti UGHC).
- Analisis Elongasi Geosentrik (Sudut pisah pusat Bumi-Bulan-Matahari).
- Prediksi Konjungsi Presisi (Ijtima).

## Catatan

- Sistem ini **tidak menggantikan observasi lapangan** (Rukyatul Hilal aktual).
- Hasil visibilitas adalah perhitungan matematis astronomis menggunakan limitasi kriteria visibilitas spesifik, bukan penampakan visual (Optik).
