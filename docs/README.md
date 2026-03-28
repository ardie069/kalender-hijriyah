# 📚 Documentation — Hilal Scope

Folder ini berisi dokumentasi teknis dan konseptual terkait **Hilal Scope**, sebuah engine falak modern berbasis NASA SPICE.

Dokumentasi ini ditujukan bagi pengembang, astronom amatir, dan siapa pun yang ingin memahami metodologi di balik sistem penentuan kalender Hijriyah ini.

---

## 📄 Daftar Dokumen

- [**api-design.md**](file:///home/ardie069/Projects/hilal-scope/kalender-hijriyah/docs/api-design.md)  
  Desain endpoint API v4, skema request/response, dan telemetri hilal.

- [**astronomical-model.md**](file:///home/ardie069/Projects/hilal-scope/kalender-hijriyah/docs/astronomical-model.md)  
  Detail teknis penggunaan NASA SPICE, integrasi CGO, dan strategi deployment di Vercel.

- [**calendar-logic.md**](file:///home/ardie069/Projects/hilal-scope/kalender-hijriyah/docs/calendar-logic.md)  
  Logika orkestrasi kalender (KHGT, MABIMS, Umm al-Qura) dan proses rollover tanggal.

- [**hisab-vs-rukyat.md**](file:///home/ardie069/Projects/hilal-scope/kalender-hijriyah/docs/hisab-vs-rukyat.md)  
  Diskusi mengenai pendekatan hisab (perhitungan) dan rukyat (observasi) dalam sistem ini.

- [**known-limitations.md**](file:///home/ardie069/Projects/hilal-scope/kalender-hijriyah/docs/known-limitations.md)  
  Batasan teknis, akurasi data, dan tantangan pada lintang ekstrem.

---

## 🌍 Quick Access

- **API Production**: [https://kalender-hijriyah-api.vercel.app/](https://kalender-hijriyah-api.vercel.app/)
- **Project Root**: [README.md](../README.md)

---

> [!NOTE]  
> Hilal Scope dirancang untuk menjembatani perhitungan astronomi presisi tinggi dengan kriteria fikih yang beragam, memberikan satu platform untuk komparasi kalender Hijriyah global.

## ⚠️ Disclaimer & Credits

- *Astronomical calculations are powered by NASA's SPICE Toolkit (NAIF).*
- *This project is an independent educational tool and is not affiliated with or endorsed by NASA.*
- *All calculations are for informational purposes only and should not be used for critical decision-making.*
