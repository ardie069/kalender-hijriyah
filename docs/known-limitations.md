# ⚠️ Known Limitations & Disclaimers

Meskipun Hilal Scope menggunakan data tingkat wahana antariksa, terdapat beberapa batasan ilmiah dan teknis yang perlu dipahami:

## 🔭 Batasan Astronomi & Falak

1. **Prediksi vs Observasi**: Sistem ini memberikan prediksi visibilitas (*imkan rukyat*), bukan konfirmasi visual lapangan. Kondisi lokal seperti awan, polusi cahaya, dan hamburan atmosfer aktual tidak diperhitungkan.
2. **Refraksi Atmosfer**: Koreksi refraksi menggunakan model standar Bennett/Sæmundsson. Pada kondisi cuaca ekstrem, pembiasan nyata bisa sedikit berbeda.
3. **Lintang Ekstrem**: Pada wilayah di atas ±65° (dekat kutub), fenomena matahari tengah malam atau malam kutub membuat penentuan waktu Sunset/Fajr menjadi sangat kompleks. Hilal Scope menggunakan pendekatan durasi proporsional pada wilayah ini, namun akurasi tetap memiliki batasan teoretis.
4. **Resolusi Waktu**: Akurasi bisection search untuk Sunset adalah ±1 detik, yang sangat memadai untuk kalender namun mungkin memiliki deviasi sangat kecil untuk kebutuhan sains murni.

## 💻 Batasan Teknis & Deployment

1. **Vercel Execution Limit**: Pada model deployment serverless, fungsi *Global Scan* (KHGT) yang memindai seluruh dunia secara intensif kadang mendekati batas waktu eksekusi (timeout) jika terjadi beban server yang tinggi.
2. **Kernel Updates**: Sistem bergantung pada file kernel NASA (`.bsp`). Jika ada pembaruan kernel dari NASA (misalnya lompatan detik baru), sistem perlu diperbarui secara manual dengan file kernel terbaru.
3. **Thread-Safety Lock**: Karena adanya *Global Lock* pada engine CSPICE, request dalam jumlah yang sangat masif dalam satu waktu yang bersamaan dapat mengalami sedikit antrean (latensi) pada level engine.

## ⚖️ Penafian (Disclaimer)

Hilal Scope adalah alat bantu ilmiah untuk komparasi dan edukasi. Penetapan resmi hari raya atau awal bulan Hijriyah tetap menjadi otoritas lembaga keagamaan atau pemerintah di masing-masing negara.
