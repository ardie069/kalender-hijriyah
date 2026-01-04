# 🔌 API Design

API dirancang agar:

- deterministik
- eksplisit
- bisa dijelaskan

## Endpoint Utama

### `/hijri-date`

Mengembalikan tanggal Hijriyah hari ini.

Response:

```json
{
  "hijriDate": {
    "year": 1447,
    "month": 7,
    "day": 15
  }
}
```

### `/hijri-end-month`

Prediksi akhir bulan Hijriyah (29 atau 30).

### `/hijri-explain`

Menjelaskan mengapa tanggal tersebut dipilih.

Tujuan endpoint ini adalah transparansi logika, bukan pembenaran.
