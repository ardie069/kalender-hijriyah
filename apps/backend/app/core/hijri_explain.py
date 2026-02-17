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

    # 1. Parameter Lokasi & Kriteria
    # Gunakan MABIMS untuk Rukyat, Wujudul Hilal untuk Hisab
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

    # 4. Konstruksi Penjelasan Dasar
    explanation: Dict[str, Any] = {
        "method": method,
        "after_sunset": bool(after_sunset),
        "criteria_used": SELECTED_CRITERIA,
        "reasoning": [],
    }

    # Logika Penjelasan Berdasarkan Metode
    if method == "global":
        explanation["reasoning"].append(
            "Mengikuti standar kalender aritmatika internasional (Umm al-Qura/Mekkah)."
        )
        explanation["reasoning"].append(
            "Pergantian bulan sudah ditetapkan berdasarkan siklus rata-rata bulan tanpa bergantung pada visibilitas lokal."
        )

    else:
        # Pengecekan Lagging (Historis) hanya untuk Rukyat/Hisab
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

        # 5. Uji Materi Akhir Bulan (Isbat Mode)
        # Hanya tampilkan data astronomis jika BUKAN global dan sudah mendekati tgl 29
        effective_day = baseline["day"]
        if effective_day >= 28:  # Mulai monitor sejak hari ke-28
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

            # Masukkan data teknis hanya jika user butuh (Hisab/Rukyat)
            explanation["astronomical_data"] = {
                "conjunction_before_sunset": bool(conj_jd < sunset_jd),
                "moon_altitude": float(round(vis["moon_altitude"], 2)),
                "elongation": float(round(vis["elongation"], 2)),
                "is_visible": bool(vis["is_visible"]),
            }

            # Detail Keputusan tgl 29
            if effective_day == 29 and after_sunset:
                if method == "hisab":
                    res = (
                        "Bulan Baru"
                        if conj_jd < sunset_jd
                        else "Istikmal (Genapkan 30 Hari)"
                    )
                    explanation["reasoning"].append(
                        f"Keputusan: {res} berdasarkan kriteria posisi konjungsi."
                    )
                elif method == "rukyat":
                    res = (
                        "Bulan Baru"
                        if vis["is_visible"]
                        else "Istikmal (Genapkan 30 Hari)"
                    )
                    explanation["reasoning"].append(
                        f"Keputusan: {res} berdasarkan kriteria visibilitas {SELECTED_CRITERIA}."
                    )

    return explanation
