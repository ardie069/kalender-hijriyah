# 🧠 Calendar Logic

Dokumen ini menjelaskan **alur penentuan tanggal Hijriyah**
yang digunakan oleh sistem.

## Prinsip Dasar

1. Hari Hijriyah **dimulai saat Matahari terbenam (Maghrib)**
2. Tidak ada pergantian hari Hijriyah di tengah malam
3. Baseline tanggal diambil dari **siang hari (noon)** untuk stabilitas
4. Evaluasi akhir bulan **hanya terjadi pada hari ke-29 Hijriyah**

## Alur Umum

1. Tentukan lokasi dan zona waktu
2. Hitung waktu sunset lokal
3. Ambil baseline Hijriyah dari Julian Day siang hari
4. Tentukan apakah sekarang sebelum atau sesudah Maghrib
5. Terapkan logika sesuai metode:
   - Global
   - Hisab
   - Rukyat

## Mengapa Baseline Noon?

Menggunakan noon (±12:00) mencegah:

- loncatan tanggal di tengah malam
- inkonsistensi antara jam 23:59 dan 00:01
- efek samping dari pembulatan Julian Day

Baseline noon **tidak menentukan bulan baru**, hanya titik referensi aman.
