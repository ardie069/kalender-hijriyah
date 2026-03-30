# 🔌 API Design Hilal Scope (v4)

API Hilal Scope dirancang untuk memberikan transparansi penuh terhadap data astronomi yang digunakan dalam penentuan tanggal Hijriyah. Semua perhitungan astronomi bersifat **toposentrik** (berbasis posisi pengamat di permukaan bumi) untuk akurasi maksimal bagi pengguna lokal.

## Endpoint Utama

### `GET /api/v4/hijri/date`

Mengembalikan tanggal Hijriyah berdasarkan berbagai metode penetapan (KHGT, MABIMS, Umm al-Qura) beserta data prediksi hilal untuk bulan berikutnya.

**Parameter Query:**

- `date` (opsional): Tanggal Gregorian dalam format ISO 8601 atau `YYYY-MM-DD`. Default: waktu saat ini (UTC).
- `lat` (**wajib**): Latitude lokasi pengamat (Contoh: `-6.17`).
- `lon` (**wajib**): Longitude lokasi pengamat (Contoh: `106.82`).

**Struktur Response (Detail):**

```json
{
  "status": "success",
  "data": {
    "gregorian_date": "2026-03-20T11:00:00Z",
    "location": { "latitude": -6.17, "longitude": 106.82, "timezone": "Asia/Jakarta" },
    "methods": {
      "KHGT": {
        "hijri_date": { "day": 1, "month": 10, "year": 1447 },
        "prediction": { 
          "is_new_month": true, 
          "altitude": 5.23, 
          "elongation": 8.45,
          "khgt_global_valid": true,
          "khgt_america_exception": false
        }
      },
      "UMM_AL_QURA": {
        "hijri_date": { "day": 1, "month": 10, "year": 1447 },
        "reference_altitude": 4.12,
        "reference_elongation": 7.89
      }
    }
  }
}
```

**Response Fields (`MethodResult`):**
- `hijri_date`: Objek berisi `day`, `month`, `year`, dan `month_name`.
- `current_altitude` / `current_elongation`: Data posisi bulan saat ini di posisi pengamat.
- `reference_altitude` / `reference_elongation`: Data astronomi di titik referensi metode (Makkah untuk Umm al-Qura, Sabang untuk MABIMS).
- `prediction`: Detail prediksi visibilitas hilal pada petang hari ke-29.

---

### `GET /api/v4/hijri/calendar`

Memberikan simulasi kalender tahunan berdasarkan metode tertentu.

**Parameter Query:**

- `year` (wajib): Tahun yang ingin dicari (Contoh: `2026`).
- `lat` (wajib): Latitude lokasi pengamat.
- `lon` (wajib): Longitude lokasi pengamat.
- `method` (opsional): Metode yang digunakan (`KHGT`, `MABIMS`, `UMM_AL_QURA`). Default: `KHGT`.

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
- `altitude`: Tinggi Bulan dari ufuk (derajat).
- `altitude_apparent`: Tinggi Bulan setelah koreksi refraksi atmosfer (derajat).
- `azimuth`: Sudut horizontal Bulan (derajat).
- `elongation`: Jarak sudut Bulan-Matahari (**Topozentrik**).
- `illumination`: Persentase cahaya piringan Bulan (0-1).
- `distance_km`: Jarak Bumi-Bulan (km).
- `phase_name`: Fase Bulan (New Moon, Waxing Crescent, dll).
- `moonrise` / `moonset`: Waktu terbit/terbenam Bulan di lokasi lokal.

---

### `GET /api/v4/prayer/times`

Menghitung waktu sholat harian berdasarkan koordinat.

**Parameter Query:**
- `lat` (wajib): Latitude lokasi pengamat.
- `lon` (wajib): Longitude lokasi pengamat.
- `date` (opsional): Tanggal format `YYYY-MM-DD`.
- `method` (opsional): Standar hitungan (`KEMENAG`, `MWL`, `UMM_AL_QURA`, dll).

---

## Prinsip Desain

1. **Toposentrik & Wajib Koordinat**: Penggunaan data NASA SPICE yang presisi mewajibkan parameter `lat` dan `lon` karena semua perhitungan (altitude, sunset, fajr) didasarkan pada posisi pengamat di permukaan bumi.
2. **Deterministik**: Input yang sama selalu menghasilkan output yang sama karena menggunakan kernel NASA SPICE.
3. **Eksplisit**: Semua parameter yang digunakan dalam keputusan (seperti altitude dan elongation) disertakan dalam respons.
4. **Sinkronisasi Fase**: Penentuan status `New Moon` disinkronkan secara global untuk memastikan konsistensi konjungsi.
