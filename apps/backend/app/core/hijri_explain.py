from typing import Any, Dict


def explain_hijri_decision(method: str, result) -> Dict[str, Any]:

    meta = result.metadata
    decision = meta.get("decision")

    after_sunset = True
    if decision == "before_maghrib":
        after_sunset = False

    explanation = {
        "method": method,
        "after_sunset": after_sunset,
        "criteria_used": "Unknown",
        "reasoning": [],
        "astronomical_data": None,
        "decision": decision,
    }

    # ==========================
    # UMM AL QURA
    # ==========================
    if method == "umm_al_qura":
        explanation["criteria_used"] = "Arithmetic (Umm al-Qura)"
        explanation["reasoning"] = [
            "Mengikuti kalender aritmatika global.",
            "Tidak bergantung pada visibilitas hilal lokal.",
        ]
        return explanation

    # ==========================
    # LOCAL HISAB
    # ==========================
    if method == "local_hisab":
        explanation["criteria_used"] = "Wujudul Hilal"

        if meta.get("decision") == "new_month":
            explanation["reasoning"].append(
                "Konjungsi terjadi sebelum matahari terbenam."
            )
        elif meta.get("decision") == "istikmal_30":
            explanation["reasoning"].append(
                "Konjungsi belum terjadi sebelum matahari terbenam."
            )
        else:
            explanation["reasoning"].append("Belum masuk fase evaluasi bulan baru.")

        return explanation

    # ==========================
    # LOCAL RUKYAT
    # ==========================
    if method == "local_rukyat":
        explanation["criteria_used"] = "MABIMS"

        vis = meta.get("visibility")

        if vis:
            explanation["astronomical_data"] = vis

            if meta.get("decision") == "new_month":
                explanation["reasoning"].append(
                    "Hilal memenuhi kriteria visibilitas MABIMS."
                )
            elif meta.get("decision") == "istikmal_30":
                explanation["reasoning"].append(
                    "Hilal tidak memenuhi kriteria visibilitas MABIMS."
                )
        else:
            explanation["reasoning"].append("Belum masuk fase evaluasi visibilitas.")

        return explanation

    # ==========================
    # UGHC
    # ==========================
    if method == "ughc":
        explanation["criteria_used"] = "Turki 2016 (Geosentrik)"

        vis = meta.get("visibility_data")

        if vis:
            explanation["astronomical_data"] = vis

            if meta.get("global_visible"):
                explanation["reasoning"].append(
                    "Hilal memenuhi kriteria global sebelum 24:00 UTC."
                )
            else:
                explanation["reasoning"].append(
                    "Tidak ada wilayah global yang memenuhi kriteria sebelum 24:00 UTC."
                )

        return explanation

    return explanation
