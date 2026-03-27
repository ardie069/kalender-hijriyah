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
2. Ambil tanggal Hijriyah **Tabular** sebagai pijakan awal untuk estimasi bulan dan tahun (`baseH`).
3. Cari waktu **Ijtima (Konjungsi)** sebelumnya dalam rentang 30 hari menggunakan `FindPreviousIjtima`.
4. Tentukan waktu evaluasi (Sunset/Maghrib) berdasarkan kriteria:
   - **UMM_AL_QURA**: Lokasi Makkah.
   - **MABIMS**: Lokasi Sabang (titik barat Indonesia).
   - **Lokal**: Lokasi observer (`lat/lon`).
   - **KHGT**: Pemindaian Global (tidak terikat satu lokasi).
5. Tentukan apakah waktu request *sebelum* atau *sesudah* Maghrib.
6. Tentukan `monthStartDate` (H+1 atau H+2 dari Ijtima) berdasarkan hasil evaluasi kriteria (`isNewMonth`).
7. Hitung selisih hari (`daysElapsed`) antara target dengan awal bulan.
8. Lakukan koreksi **Rollover** (jika `hDay < 1` atau `hDay > 30`) untuk menyesuaikan bulan dan tahun secara dinamis.

## Mengatasi Limitasi "Hari H+1"

Sistem secara dinamis menghitung hari Hijriyah (`hDay`) berdasarkan selisih waktu antara tanggal yang diminta dengan `monthStartDate` yang sudah dievaluasi. Logika ini memastikan pergantian bulan (rollover) terjadi secara akurat: menjadi hari ke-1 bulan selanjutnya jika kriteria hisab/rukyat terpenuhi, atau mengunci ke hari ke-30 (Istikmal) apabila hilal tidak terlihat. Penyesuaian bulan dan tahun juga dilakukan secara otomatis jika `hDay` berada di luar rentang bulan saat ini.
