from datetime import datetime, timedelta
import pytz
from .hijri_calculator import get_hijri_date
from .sun_times import get_sunset_time
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .julian import jd_from_datetime


def predict_end_of_month(lat, lon, method, timezone, *, ts, eph, sun, moon, earth):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    today = get_hijri_date(
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

    if today["day"] != 29:
        return {
            "method": method,
            "today": today,
            "estimated_end_of_month": None,
            "visibility": None,  # Pastikan ada key ini biar gak error di frontend
            "message": "Prediksi akhir bulan hanya tersedia pada tanggal 29 Hijriyah.",
        }

    # Logic Penentuan Estimasi Tanggal
    tomorrow_local = now_local + timedelta(days=1)
    tomorrow = get_hijri_date(
        lat,
        lon,
        method,
        timezone,
        now_local=tomorrow_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    # Ini cuma datanya saja (year, month, day)
    estimated_date = (
        today if tomorrow["month"] != today["month"] else {**today, "day": 30}
    )

    # Inisialisasi visibility_data
    visibility_data = None

    if method != "global":
        sunset_local = get_sunset_time(now_local.date(), lat, lon, timezone, ts, eph)
        if sunset_local:
            sunset_utc = sunset_local.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc, ts)
            conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)

            criteria = "MABIMS" if method == "rukyat" else "Wujudul Hilal"
            vis = evaluate_visibility(
                sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth, criteria=criteria
            )

            # MASUKKAN KE TOP LEVEL (Bukan di dalam estimated_end_of_month)
            visibility_data = {
                "moon_age": float((sunset_jd - conj_jd) * 24),
                "moon_altitude": float(vis["moon_altitude"]),
                "elongation": float(vis["elongation"]),
                "is_visible": bool(vis["is_visible"]),
            }

    return {
        "method": method,  # WAJIB ada buat isGlobal di frontend
        "today": today,
        "estimated_end_of_month": estimated_date,
        "visibility": visibility_data,  # Jalur distribusinya sekarang bener
        "message": (
            f"Berdasarkan kriteria {method.upper()}, bulan ini diperkirakan berumur "
            f"{estimated_date['day']} hari."
        ),
    }
