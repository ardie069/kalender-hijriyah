# 🔌 API Design v4

API dirancang agar:

- deterministik
- eksplisit (dikomputasi pada runtime menggunakan CSPICE)
- mendukung kalkulasi global vs lokal

## Endpoint Utama

### `GET /api/v4/hijri/date`

Mengembalikan tanggal Hijriyah hari ini beserta prediksi pergantian bulan dari berbagai metode penetapan.

Parameter Query Opsional:
- `date`: Tanggal ISO 8601 (Contoh: `2026-03-18T12:00:00Z`). Default: current time.
- `lat`: Latitude observer (Contoh: `-7.98`). Default: Malang.
- `lon`: Longitude observer (Contoh: `112.63`). Default: Malang.

Response Format:

```json
{
  "target_date_utc": "2026-03-18T12:00:00Z",
  "methods": {
    "UGHC_KHGT": {
      "hijri_date": {
        "day": 30,
        "month": 9,
        "month_name": "Ramadan",
        "year": 1447
      },
      ...
    },
    ...
  }
}
```

Tujuan endpoint ini adalah memberikan komparasi langsung hasil Hisab Global, Umm al-Qura, dan berbagai metode Hisab Lokal dalam satu *single source of truth*.
