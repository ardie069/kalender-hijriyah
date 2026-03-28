# 🔌 API Design Hilal Scope (v4)

API Hilal Scope dirancang untuk memberikan transparansi penuh terhadap data astronomi yang digunakan dalam penentuan tanggal Hijriyah.

## Endpoint Utama

### `GET /api/v4/hijri/date`

Mengembalikan tanggal Hijriyah berdasarkan berbagai metode penetapan (KHGT, MABIMS, Umm al-Qura, Wujudul Hilal) beserta data prediksi hilal untuk bulan berikutnya.

**Parameter Query:**

- `date` (opsional): Tanggal Gregorian dalam format ISO 8601 atau `YYYY-MM-DD`. Default: waktu saat ini (UTC).
- `lat` (wajib): Latitude lokasi pengamat (Contoh: `-6.17`).
- `lon` (wajib): Longitude lokasi pengamat (Contoh: `106.82`).

**Struktur Response:**

```json
{
  "status": "success",
  "data": {
    "gregorian_date": "2026-03-20T11:00:00Z",
    "location": {
      "latitude": -6.17,
      "longitude": 106.82,
      "timezone": "Asia/Jakarta"
    },
    "methods": {
      "KHGT": {
        "hijri_date": { "day": 1, "month": 10, "month_name": "Syawal", "year": 1447 },
        "prediction": {
          "is_new_month": true,
          "altitude": 5.23,
          "elongation": 8.45,
          "khgt_global_valid": true,
          "khgt_america_exception": false
        }
      },
      "MABIMS": {
        "hijri_date": { "day": 1, "month": 10, "month_name": "Syawal", "year": 1447 },
        "reference_altitude": 3.45,
        "reference_elongation": 6.78,
        "prediction": {
          "is_new_month": true,
          "location": { "name": "Sabang, Indonesia" }
        }
      }
    }
  }
}
```

### `GET /api/v4/moon/telemetry`

Data astronomi Bulan real-time untuk lokasi pengamat.

**Response:**

- `altitude`: Tinggi Bulan (derajat).
- `azimuth`: Azimuth Bulan (derajat).
- `elongation`: Jarak sudut Bulan-Matahari (derajat).
- `illumination`: Persentase cahaya Bulan (0-1).
- `age_hours`: Umur Bulan sejak Ijtima (jam).

---

## Prinsip Desain

1. **Deterministik**: Input yang sama selalu menghasilkan output yang sama karena menggunakan kernel NASA SPICE.
2. **Eksplisit**: Semua parameter yang digunakan dalam keputusan (seperti altitude dan elongation) disertakan dalam respons.
3. **Multi-Metode**: Memberikan fleksibilitas bagi pengguna untuk memilih standar kalender yang sesuai dengan keyakinan atau kebutuhan administratif mereka.
