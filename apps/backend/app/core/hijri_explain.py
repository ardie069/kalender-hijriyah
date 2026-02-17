import pytz
from datetime import datetime, time
from typing import Any, Dict, Optional

from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .config import DEFAULT_LOCATION

from .hijri_calculator import get_hijri_date, _check_historical_lag


def explain_hijri_decision(
    lat: float,
    lon: float,
    method: str,
    timezone: str,
    *,
    now_local: datetime,
    ts: Any,
    eph: Any,
    sun: Any,
    moon: Any,
    earth: Any,
) -> Dict[str, Any]:

    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # 1. Parameter Lokasi & Kriteria
    ref_lat, ref_lon, ref_zone = (
        DEFAULT_LOCATION["global"] if method == "global" else (lat, lon, timezone)
    )

    SELECTED_CRITERIA = "MABIMS" if method == "rukyat" else "Wujudul Hilal"

    # 2. Ambil Sunset & Baseline Waktu
    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    if not sunset_local:
        return {"error": "sunset_not_available", "method": method}

    after_sunset = now_local >= sunset_local
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    # Baseline Arithmetic (Hisab Global) di tengah hari
    noon_local = tz.localize(datetime.combine(now_local.date(), time(12, 0)))
    noon_jd = jd_from_datetime(noon_local.astimezone(pytz.utc), ts)
    baseline = julian_to_hijri(noon_jd)

    # 3. Sinkronisasi dengan Engine Utama
    final_date = get_hijri_date(
        lat,
        lon,
        method,
        timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    # 4. Konstruksi Penjelasan (Structured Response)
    explanation: Dict[str, Any] = {
        "method": method,
        "baseline_date": baseline,
        "final_hijri_date": final_date,
        "after_sunset": bool(after_sunset),
        "criteria_used": SELECTED_CRITERIA,
        "reasoning": [],
    }

    # --- Metode Global ---
    if method == "global":
        explanation["decision"] = "global_standard"
        explanation["reasoning"].append(
            "Mengikuti standar kalender aritmatika internasional (Umm al-Qura/Mekkah)."
        )
        explanation["reasoning"].append(
            "Pergantian bulan sudah ditetapkan berdasarkan siklus rata-rata bulan tanpa bergantung pada visibilitas lokal."
        )
        return explanation

    # --- Metode Hisab & Rukyat ---

    # Pengecekan Lagging (Historis)
    is_lagging = _check_historical_lag(
        baseline,
        noon_jd,
        lat,
        lon,
        timezone,
        ts,
        eph,
        sun,
        moon,
        earth,
        criteria=SELECTED_CRITERIA,
    )

    if is_lagging:
        explanation["reasoning"].append(
            "Bulan berjalan lebih lambat (Lag -1) karena posisi hilal di awal bulan ini tidak memenuhi syarat visibilitas."
        )
    else:
        explanation["reasoning"].append(
            "Siklus bulan berjalan normal sesuai dengan kriteria astronomis yang ditetapkan."
        )

    # Tentukan effective_day berdasarkan baseline (sebelum adjust)
    effective_day = baseline["day"]

    # Jika bukan hari ke 28+, tidak perlu evaluasi astronomis detail
    if effective_day < 28:
        explanation["decision"] = "no_evaluation_needed"
        return explanation

    # 5. Evaluasi Astronomis (hari ke-28+)
    conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)
    vis = evaluate_visibility(
        sunset_utc,
        lat,
        lon,
        conj_jd,
        ts,
        sun,
        moon,
        earth,
        criteria=SELECTED_CRITERIA,
    )

    # Data konjungsi (selalu tersedia untuk hisab & rukyat pada hari 28+)
    explanation["conjunction"] = {
        "conjunction_jd": float(conj_jd),
        "conjunction_before_sunset": bool(conj_jd < sunset_jd),
    }

    # Data visibilitas (selalu tersedia untuk hisab & rukyat pada hari 28+)
    explanation["visibility"] = {
        "moon_altitude": float(round(vis["moon_altitude"], 2)),
        "elongation": float(round(vis["elongation"], 2)),
        "moon_age": float(round(vis["moon_age"], 2)),
        "is_visible": bool(vis["is_visible"]),
    }

    # Keputusan pada tanggal 29 setelah maghrib
    if effective_day == 29 and after_sunset:
        if method == "hisab":
            if conj_jd < sunset_jd:
                explanation["decision"] = "new_month"
                explanation["reasoning"].append(
                    "Keputusan: Bulan Baru — konjungsi terjadi sebelum matahari terbenam (kriteria Wujudul Hilal terpenuhi)."
                )
            else:
                explanation["decision"] = "istikmal_30"
                explanation["reasoning"].append(
                    "Keputusan: Istikmal (Genapkan 30 Hari) — konjungsi belum terjadi saat matahari terbenam."
                )
        elif method == "rukyat":
            if vis["is_visible"]:
                explanation["decision"] = "new_month"
                explanation["reasoning"].append(
                    f"Keputusan: Bulan Baru — hilal terlihat berdasarkan kriteria visibilitas {SELECTED_CRITERIA}."
                )
            else:
                explanation["decision"] = "istikmal_30"
                explanation["reasoning"].append(
                    f"Keputusan: Istikmal (Genapkan 30 Hari) — hilal tidak memenuhi kriteria visibilitas {SELECTED_CRITERIA}."
                )
    else:
        explanation["decision"] = "no_evaluation_needed"
        explanation["reasoning"].append(
            "Belum waktunya sidang isbat (menunggu sore hari tanggal 29)."
        )

    return explanation
