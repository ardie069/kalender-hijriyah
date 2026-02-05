import pytz
from datetime import datetime, time
from typing import Any, Dict

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

    # 1. Lokasi & Parameter Material
    ref_lat, ref_lon, ref_zone = (
        DEFAULT_LOCATION["global"] if method == "global" else (lat, lon, timezone)
    )

    # Gunakan kriteria yang sama dengan engine utama biar konsisten
    SELECTED_CRITERIA = "Wujudul Hilal"

    # 2. Ambil Sunset & Baseline
    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    if not sunset_local:
        return {"error": "sunset_not_available", "method": method}

    after_sunset = now_local >= sunset_local
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    # Baseline Arithmetic (Hisab Global)
    noon_local = tz.localize(datetime.combine(now_local.date(), time(12, 0)))
    noon_jd = jd_from_datetime(noon_local.astimezone(pytz.utc), ts)
    baseline = julian_to_hijri(noon_jd)

    # 3. Analisis Sejarah (Lagging Check)
    # Ini buat ngejelasin kenapa Rukyat bisa selisih -1 dari Global
    is_lagging = False
    if method == "rukyat":
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

    # Ambil hasil final dari engine utama
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

    # 4. Konstruksi Penjelasan (Dialektika Data)
    explanation: Dict[str, Any] = {
        "method": method,
        "after_sunset": after_sunset,
        "is_lagging_from_start": is_lagging,
        "criteria_used": SELECTED_CRITERIA,
        "final_hijri_date": final_date,
        "reasoning": [],
    }

    # Logika Penjelasan:
    if method == "global":
        explanation["reasoning"].append(
            "Mengikuti standar aritmatika global (Mekkah/Umm al-Qura)."
        )

    elif method == "rukyat":
        if is_lagging:
            explanation["reasoning"].append(
                "Bulan ini dimulai lebih lambat (Lag -1) karena hilal awal bulan tidak memenuhi kriteria visibilitas."
            )
        else:
            explanation["reasoning"].append(
                "Bulan ini berjalan serentak dengan Hisab karena hilal awal bulan berhasil teramati/memenuhi kriteria."
            )

    # 5. Uji Materi khusus hari ke-29 (Sidang Isbat)
    # Kita cek baseline day-nya.
    effective_day = baseline["day"]
    if after_sunset:
        effective_day += 1  # Pendekatan kasar buat cek tgl 29

    if effective_day >= 29:
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

        explanation["astronomical_data"] = {
            "conjunction_before_sunset": conj_jd < sunset_jd,
            "moon_altitude": round(vis["moon_altitude"], 2),
            "elongation": round(vis["elongation"], 2),
            "is_visible": vis["is_visible"],
        }

        if effective_day == 29 and after_sunset:
            if method == "hisab":
                res = "Bulan Baru" if conj_jd < sunset_jd else "Istikmal (30 Hari)"
                explanation["reasoning"].append(
                    f"Keputusan akhir bulan: {res} berdasarkan posisi konjungsi."
                )
            elif method == "rukyat":
                res = "Bulan Baru" if vis["is_visible"] else "Istikmal (30 Hari)"
                explanation["reasoning"].append(
                    f"Keputusan akhir bulan: {res} berdasarkan kriteria visibilitas {SELECTED_CRITERIA}."
                )

    return explanation
