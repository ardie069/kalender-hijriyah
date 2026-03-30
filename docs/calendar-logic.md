# 🧠 Calendar Logic Hilal Scope (v4)

Dokumen ini menjelaskan alur orkestrasi dan logika falak yang digunakan Hilal Scope untuk menentukan tanggal Hijriyah secara akurat.

## Prinsip Utama

1. **Maghrib sebagai Titik Nol**: Hari Hijriyah berganti tepat saat Matahari terbenam secara lokal.
2. **Evaluasi Hari ke-29**: Penentuan apakah bulan berjalan berjumlah 29 atau 30 hari dilakukan dengan mengevaluasi visibilitas hilal pada petang hari ke-29.
3. **NASA SPICE Precision**: Semua posisi benda langit dihitung menggunakan kernel ephemeris NASA untuk mendapatkan akurasi tingkat tinggi.

## Kriteria Penetapan

### 🌍 1. KHGT (Kalender Hijriyah Global Tunggal)
Berdasarkan Kongres Turki 2016, awal bulan dimulai secara global jika:
- Di bagian mana pun di bumi (**Global Scan**), tinggi hilal minimal **5°** dan elongasi minimal **8°** saat Matahari terbenam.
- Terjadi sebelum pukul 00:00 UTC.
- **Eksepsi Benua Amerika**: Jika syarat terpenuhi di benua Amerika (setelah 00:00 UTC), maka tetap berlaku untuk hari yang sama secara global jika kriteria terpenuhi sebelum fajar di Selandia Baru. Hilal Scope menggunakan pemindaian 24 jam yang deterministik untuk memastikan validitas eksepsi ini.

### 🇮🇩 2. MABIMS (New Criteria 2022)
Digunakan secara resmi di Indonesia, Malaysia, Brunei, dan Singapura:
- **Titik Referensi**: Sabang, Indonesia (titik paling barat untuk memaksimalkan peluang visibilitas).
- **Syarat**: Tinggi hilal minimal **3°** dan elongasi minimal **6.4°**.
- Jika syarat terpenuhi di Sabang, maka esok harinya adalah tanggal 1 Hijriyah untuk seluruh wilayah zona MABIMS.

### 🕋 3. Umm al-Qura
Kriteria resmi Arab Saudi:
- **Titik Referensi**: Ka'bah, Makkah.
- Konjungsi (Ijtima) terjadi sebelum Matahari terbenam di Makkah.
- Bulan terbenam setelah Matahari terbenam di Makkah.
- Hilal Scope menyediakan `reference_altitude` dan `reference_elongation` khusus untuk koordinat Makkah pada metode ini.

---

## Alur Orkestrasi (Service Layer)

1. **Inisialisasi**: Menentukan waktu target dan lokasi pengamat.
2. **Backtrack Ijtima**: Mencari waktu konjungsi Bulan-Matahari sebelumnya untuk menentukan batasan bulan Hijriyah.
3. **Sunset Detection**: Menggunakan *Bisection Search* untuk menemukan waktu terbenam Matahari yang presisi di lokasi referensi (Sabang untuk MABIMS, Makkah untuk Umm al-Qura, atau Lokal).
4. **Visibilitas Hilal**: Menentukan apakah hilal memenuhi kriteria spesifik metode pada waktu Maghrib tersebut.
5. **Dynamic Rollover**: Berdasarkan status `isNewMonth`, sistem menghitung `monthStartDate`. Selisih hari antara `target_date` dan `monthStartDate` akan menentukan angka hari Hijriyah, dengan penyesuaian otomatis jika terjadi pergantian bulan atau tahun di luar kalender tabular.

## Penentuan Status Indikator

Hilal Scope menyediakan indikator visual untuk membantu pengguna memahami status kriteria:
- 🟢 **Memenuhi**: Semua syarat (Tinggi & Elongasi) terpenuhi.
- 🟡 **Hampir**: Salah satu syarat terpenuhi atau sangat dekat dengan ambang batas.
- 🔴 **Belum**: Tidak memenuhi syarat minimum.
