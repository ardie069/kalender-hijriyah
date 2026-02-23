from typing import Any, Dict


def explain_hijri_decision(method: str, result) -> Dict[str, Any]:
    meta = result.metadata
    decision = meta.get("decision")
    after_sunset = decision != "before_maghrib"

    explanation = {
        "method": method,
        "after_sunset": after_sunset,
        "criteria_used": "Unknown",
        "reasoning": [],
        "astronomical_data": None,
        "decision": decision,
    }

    # ==========================
    # UMM AL QURA (Statis/Aritmatis)
    # ==========================
    if method == "umm_al_qura":
        explanation["criteria_used"] = "Arithmetic (Umm al-Qura)"
        explanation["reasoning"] = [
            "Menggunakan kalender aritmatika resmi Al-Haram, Makkah.",
            "Keputusan bersifat global-statis dan tidak bergantung pada posisi hilal lokal.",
        ]
        return explanation

    # ==========================
    # LOCAL HISAB (Wujudul Hilal)
    # ==========================
    if method == "local_hisab":
        explanation["criteria_used"] = "Wujudul Hilal"

        if decision == "new_month":
            explanation["reasoning"] = [
                "Ijtima' (konjungsi) terjadi sebelum matahari terbenam.",
                "Bulan sudah berada di atas ufuk saat Maghrib tiba.",
            ]
        elif decision == "istikmal_30":
            explanation["reasoning"] = [
                "Konjungsi belum terjadi atau posisi bulan masih di bawah ufuk.",
                "Bulan berjalan digenapkan menjadi 30 hari (Istikmal).",
            ]
        else:
            explanation["reasoning"] = [
                "Menunggu waktu matahari terbenam untuk evaluasi posisi hilal."
            ]

        return explanation

    # ==========================
    # LOCAL RUKYAT (MABIMS)
    # ==========================
    if method == "local_rukyat":
        explanation["criteria_used"] = "MABIMS (Baru)"
        vis = meta.get("visibility")

        if vis:
            explanation["astronomical_data"] = vis
            if decision == "new_month":
                explanation["reasoning"] = [
                    f"Tinggi hilal ({vis['moon_altitude']:.2f}°) dan elongasi ({vis['elongation']:.2f}°) memenuhi kriteria.",
                    "Secara astronomis, hilal mungkin untuk diamati (Imkanur Rukyat).",
                ]
            else:
                explanation["reasoning"] = [
                    "Posisi hilal belum memenuhi batas minimum MABIMS (Tinggi 3°, Elongasi 6.4°).",
                    "Bulan berjalan digenapkan menjadi 30 hari.",
                ]
        else:
            explanation["reasoning"] = [
                "Data visibilitas akan muncul setelah waktu Maghrib lokal."
            ]

        return explanation

    # ==========================
    # UGHC (Unified Global / KHGT)
    # ==========================
    if method == "ughc":
        explanation["criteria_used"] = "Turkey 2016 (Unified Global)"
        vis = meta.get("visibility_data")

        if vis:
            explanation["astronomical_data"] = vis
            reasoning = [
                "Menggunakan prinsip Matlak Global: satu tanggal untuk seluruh dunia tanpa sekat wilayah.",
                f"Hasil pemindaian titik optimal: Lat {vis.get('lat')}, Lon {vis.get('lon')}.",
            ]

            if meta.get("global_visible"):
                lon = vis.get("lon", 0)
                is_america = -170 <= lon <= -30

                if is_america:
                    reasoning.append(
                        "Kriteria terpenuhi di daratan Amerika. Sesuai protokol KHGT, "
                        "ini menjadi validasi global meskipun di titik lain sudah melewati 24:00 UTC."
                    )
                else:
                    reasoning.append(
                        "Kriteria (Elongasi 8°, Tinggi 5°) terpenuhi di belahan bumi lain sebelum pukul 24:00 UTC."
                    )

                reasoning.append(
                    "Kesimpulan: Besok secara universal dimulai bulan baru (Unified Hijri)."
                )
            else:
                reasoning.append(
                    "Tidak ditemukan koordinat di seluruh permukaan bumi yang memenuhi syarat "
                    "sebelum batas waktu internasional (UTC)."
                )
                reasoning.append(
                    "Pengecekan ijtima' di Selandia Baru sebelum fajar juga tidak memberikan hasil positif."
                )
                reasoning.append(
                    "Bulan berjalan digenapkan secara global (Istikmal Tunggal)."
                )

            explanation["reasoning"] = reasoning
        else:
            explanation["reasoning"] = [
                "Melakukan sinkronisasi data astronomi geosentris untuk pemindaian global."
            ]

        return explanation

    return explanation
