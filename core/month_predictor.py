from datetime import datetime
import pytz
from julian import jd_from_datetime
from hijri_calculator import get_hijri_date
from visibility import evaluate_visibility
from conjunction import get_conjunction_time
from sun_times import get_sunset_time


def predict_end_of_month(lat, lon, method, timezone, *, ts, eph, sun, moon, earth):
    # Tanggal Hijriyah hari ini
    now_local = datetime.now(pytz.timezone(timezone))
    jd_today = jd_from_datetime(now_local, ts)
    hijri_today = get_hijri_date(lat, lon, method, timezone, ts=ts, eph=eph, sun=sun, moon=moon, earth=earth)

    # Cek tanggal 29
    if hijri_today["day"] != 29:
        return {
            "today": hijri_today,
            "estimated_end_of_month": None,
            "message": "Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.",
        }

    # Default: tetap tanggal 29
    estimated_end = {
        "day": 29,
        "month": hijri_today["month"],
        "year": hijri_today["year"],
    }

    # Hitung konjungsi
    conjunction_jd = get_conjunction_time(jd_today, ts, earth, sun, moon)

    # Prediksi akhir bulan
    is_new_month = False
    if method == "global":
        sunset_global = get_sunset_time(
            now_local.date(), 21.422487, 39.826206, "Asia/Riyadh", ts, eph
        )
        sunset_global_jd = jd_from_datetime(sunset_global.astimezone(pytz.utc), ts) # type: ignore
        if conjunction_jd < sunset_global_jd:
            is_new_month = True
    elif method == "hisab":
        sunset_user = get_sunset_time(now_local.date(), lat, lon, timezone, ts, eph)
        sunset_user_jd = jd_from_datetime(sunset_user.astimezone(pytz.utc), ts) # type: ignore
        if conjunction_jd < sunset_user_jd:
            is_new_month = True
    elif method == "rukyat":
        sunset_user = get_sunset_time(now_local.date(), lat, lon, timezone, ts, eph)
        visibility = evaluate_visibility(sunset_user, lat, lon, conjunction_jd, ts, sun, moon)
        if visibility["is_visible"]:
            is_new_month = True

    # Update jika bulan baru dimulai
    if is_new_month:
        estimated_end = {
            "day": 1,
            "month": hijri_today["month"] + 1, # type: ignore
            "year": hijri_today["year"],
        }
        if estimated_end["month"] > 12:
            estimated_end["month"] = 1
            estimated_end["year"] += 1

    return {
        "today": hijri_today,
        "estimated_end_of_month": estimated_end,
        "message": f"Prediksi bulan ini berakhir pada {estimated_end['day']}-{estimated_end['month']}-{estimated_end['year']}.",
    }
