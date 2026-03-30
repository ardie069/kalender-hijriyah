# 🔌 API Design Hilal Scope (v4)

API Hilal Scope dirancang untuk memberikan transparansi penuh terhadap data astronomi yang digunakan dalam penentuan tanggal Hijriyah. Semua perhitungan astronomi (selain metode Tabular) bersifat toposentrik-geosentrik dan wajib menyertakan koordinat pengamat.

## Endpoint Utama

### `GET /api/v4/hijri/date`

Mengembalikan tanggal Hijriyah berdasarkan berbagai metode penetapan (KHGT, MABIMS, Umm al-Qura, Wujudul Hilal) beserta data prediksi hilal untuk bulan berikutnya.

**Parameter Query:**

- `date` (opsional): Tanggal Gregorian dalam format ISO 8601 atau `YYYY-MM-DD`. Default: waktu saat ini (UTC).
- `lat` (wajib): Latitude lokasi pengamat (Contoh: `-6.17`).
- `lon` (wajib): Longitude lokasi pengamat (Contoh: `106.82`).

**Struktur Response (Ringkas):**

```json
{
  "status": "success",
  "data": {
    "gregorian_date": "2026-03-20T11:00:00Z",
    "location": { "latitude": -6.17, "longitude": 106.82, "timezone": "Asia/Jakarta" },
    "methods": {
      "KHGT": {
        "hijri_date": { "day": 1, "month": 10, "year": 1447 },
        "prediction": { "is_new_month": true, "altitude": 5.23, "elongation": 8.45 }
      },
      "MABIMS": {
        "hijri_date": { "day": 1, "month": 10, "year": 1447 },
        "reference_altitude": 3.45,
        "reference_elongation": 6.78
      }
    }
  }
}
```

---

### `GET /api/v4/hijri/calendar`

Memberikan simulasi kalender tahunan berdasarkan metode tertentu.

**Parameter Query:**

- `year` (wajib): Tahun yang ingin dicari (Contoh: `2026`).
- `lat` (wajib): Latitude lokasi pengamat.
- `lon` (wajib): Longitude lokasi pengamat.
- `method` (opsional): Metode yang digunakan (`KHGT`, `MABIMS`, `UMM_AL_QURA`, `TABULAR`). Default: `KHGT`.

---

### `GET /api/v4/hijri/search`

Mencari padanan tanggal Hijriyah menggunakan metode **Tabular (Aritmatika)**. Endpoint ini tidak memerlukan koordinat karena perhitungan bersifat statis.

**Parameter Query:**

- `date` (wajib): Tanggal Gregorian dalam format `YYYY-MM-DD`.

---

### `GET /api/v4/moon/telemetry`

Data astronomi Bulan real-time untuk lokasi pengamat.

**Parameter Query:**

- `lat` (wajib): Latitude lokasi pengamat.
- `lon` (wajib): Longitude lokasi pengamat.

**Response Fields:**

- `altitude`: Tinggi Bulan (derajat).
- `azimuth`: Azimuth Bulan (derajat).
- `elongation`: Jarak sudut Bulan-Matahari (derajat).
- `illumination`: Persentase cahaya Bulan (0-1).
- `age_hours`: Umur Bulan sejak Ijtima (jam).

---

### `GET /api/v4/prayer/times`

Menghitung waktu sholat harian berdasarkan koordinat dan standar otoritas tertentu.

**Parameter Query:**

- `lat` (wajib): Latitude lokasi pengamat.
- `lon` (wajib): Longitude lokasi pengamat.
- `date` (opsional): Tanggal format `YYYY-MM-DD`. Default: hari ini.
- `method` (opsional): Standar hitungan (`KEMENAG`, `MWL`, `ISNA`, `EGYPTIAN`, `UMM_AL_QURA`, `KARACHI`, `TEHRAN`, `JAKIM`, `MUIS`). Default: `KEMENAG`.

---

## Prinsip Desain

1. **Toposentrik & Wajib Koordinat**: Karena Hilal Scope menggunakan data dari NASA SPICE yang sangat presisi, setiap perhitungan (ketinggian hilal, waktu sunset, waktu sholat) didasarkan pada posisi pengamat di permukaan bumi. Oleh karena itu, parameter `lat` dan `lon` bersifat wajib untuk hampir seluruh endpoint.
2. **Deterministik**: Input yang sama selalu menghasilkan output yang sama karena menggunakan kernel NASA SPICE.
3. **Eksplisit**: Semua parameter yang digunakan dalam keputusan (seperti altitude dan elongation) disertakan dalam respons.
4. **Multi-Metode**: Memberikan fleksibilitas bagi pengguna untuk memilih standar kalender yang sesuai dengan keyakinan atau kebutuhan administratif mereka.
