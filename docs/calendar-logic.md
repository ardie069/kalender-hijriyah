# 🧠 Calendar Logic v4 (Engine Falak)

Dokumen ini menjelaskan **alur penentuan tanggal Hijriyah**
yang dikomputasi oleh API Golang v4.

## Prinsip Dasar Logika Kalender

1. Hari Hijriyah **dimulai saat Matahari terbenam (Maghrib)** di lokasi pengamat.
2. Tidak ada pergantian hari Hijriyah di tengah malam waktu lokal (00:00).
3. Evaluasi akhir bulan (penentuan umur bulan 29 atau 30 hari) **hanya terjadi pada hari ke-29 Hijriyah** setelah Maghrib.
4. Perhitungan lintasan Astronomi menggunakan Julian Date Ephemeris Timed (ET) dari NASA SPICE, dikonversi kembali ke UTC timezone-aware.

## Alur Orkestrasi (Service Layer)

1. Tentukan target `Date` (UTC) dan lokasi `Lat/Lon`.
2. Hitung waktu sunset (Maghrib) lokal via Bisection Search 24H.
3. Sinkronisasi base kalender Umm al-Qura dari pre-calculated tabular engine.
4. Tentukan apakah waktu request *sebelum* atau *sesudah* Maghrib.
5. Jalankan Evaluasi Multi-Metode:
   - **UGHC (Global)**: Melakukan scan seluruh benua Amerika dari Timur ke Barat untuk mencari terpenuhinya kriteria Turki 2016 sebelum dan sesudah 24:00 UTC, termasuk *override* Ijtima Selandia Baru.
   - **Hisab Lokal**: Menggunakan proyeksi Toposentrik untuk mengevaluasi syarat MABIMS atau Wujudul Hilal lokal.
   - **Umm al-Qura**: Mencocokkan dengan standar observatorium Makkah.

## Mengatasi Limitasi "Hari H+1"

Sistem menggunakan fungsi pembantu `PredictHijriDate` untuk mengeliminasi statik +1 hari. Ini memastikan bulan ber-rollover dengan dinamis (menjadi hari ke-1 bulan selanjutnya) jika Hisab menyatakan wujud/imkanur rukyat, dan mengunci ke hari ke-30 (Istikmal) apabila tidak terlihat.
