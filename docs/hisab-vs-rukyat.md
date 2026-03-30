# 🌙 Hisab Kontemporer vs Rukyat di API v4

Kalender Hijriyah API v4 tidak menyatukan Hisab dan Rukyat, melainkan menampilkannya secara transparan dalam objek `methods`.

## Rukyat (Imkanur Rukyat / MABIMS Baru 2022)

- Berbasis kriteria ekspektasi visibilitas optik.
- Menggunakan proyeksi **Toposentris** (koreksi lokasi permukaan bumi dan refraksi atmosfer).
- Kriteria ketat yang dianut negara-negara MABIMS:
  - **Ketinggian (Altitude Geometric) ≥ 3°** (Ketinggian nyata/apparent ~3.45° dengan refraksi)
  - **Elongasi ≥ 6.4°**

## Kalender Hijriyah Global Tunggal (KHGT)

- Diadopsi dari **Kongres Turki 2016**.
- Bertujuan menyatukan kalender Islam sedunia menggunakan Hisab Imkanur Rukyat *Global*.
- Kriteria global (Turki 2016):
  - Hilal harus mencapai **Elongasi ≥ 8°** dan **Altitude Geosentrik ≥ 5°** di belahan bumi *mana pun* sebelum pukul **24:00 UTC**.
  - **Pengecualian Rollover**: Jika kriteria di atas hanya terpenuhi *setelah* 24:00 UTC, bulan baru tetap masuk jika:
    1. Kriteria terpenuhi di wilayah **daratan Benua Amerika**.
    2. **Ijtimak** (konjungsi) terjadi sebelum subuh di **Selandia Baru**.

## Implikasi API v4

- Jika metode berbeda hasil, *field* `is_new_month` akan menampilkan hasil divergen secara lugas.
- API mengkalkulasi transisi murni pada Maghrib *lokal* pengamat untuk menghindari kebingungan perubahan kalender jam 12 malam.
