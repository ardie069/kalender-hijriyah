# 🌙 Hisab vs Rukyat di API v4

Kalender Hijriyah API v4 tidak menyatukan Hisab dan Rukyat, melainkan menampilkannya secara transparan dalam objek `methods`.

## Hisab (Wujudul Hilal)

- Menggunakan syarat *Sunset harus terjadi sesudah Ijtima (konjungsi)*.
- Altitude > 0° saat Maghrib.
- Deterministik murni geometris bumi-bulan-matahari.

## Rukyat (Imkanur Rukyat / MABIMS Baru 2021)

- Berbasis kriteria ekspektasi visibilitas optik.
- Kriteria ketat yang dianut negara-negara MABIMS:
  - **Ketinggian (Altitude) ≥ 3°**
  - **Elongasi ≥ 6.4°**

## Kalender Hijriyah Global Tunggal (KHGT / UGHC)

- Diadopsi dari **Kongres Turki 2016**.
- Bertujuan menyatukan kalender Islam sedunia menggunakan Hisab Imkanur Rukyat *Global*.
- Kriteria global:
  - Hilal harus mencapai **Elongasi ≥ 8°** dan **Altitude Geosentrik ≥ 5°** di *mana pun* di bumi (parameter optimal diukur di pantai barat benua Amerika).
  - Terdapat limitasi 24:00 UTC dan *override* jam fajar di Selandia Baru untuk memastikan tidak terjadi keterlambatan masuk bulan kalender di zona timur awal.

## Implikasi API v4

- Jika metode berbeda hasil, *field* `is_new_month` akan menampilkan hasil divergen secara lugas.
- API mengkalkulasi transisi murni pada Maghrib *lokal* pengamat untuk menghindari kebingungan perubahan kalender jam 12 malam.
