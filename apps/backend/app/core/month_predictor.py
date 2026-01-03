from datetime import datetime, timedelta
import pytz # type: ignore
from .hijri_calculator import get_hijri_date


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

    if today["day"] != 29:
        return {
            "today": today,
            "estimated_end_of_month": None,
            "message": "Prediksi akhir bulan hanya tersedia pada tanggal 29 Hijriyah.",
        }

    if tomorrow["month"] != today["month"]:
        estimated_end = today
    else:
        estimated_end = {**today, "day": 30}

    return {
        "today": today,
        "estimated_end_of_month": estimated_end,
        "message": (
            f"Perkiraan akhir bulan: "
            f"{estimated_end['day']}-{estimated_end['month']}-{estimated_end['year']}"
        ),
    }
