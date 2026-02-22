from datetime import datetime, timedelta
import pytz

from app.core.method_factory import get_method_instance
from app.core.methods.base import HijriContext


def predict_end_of_month(lat, lon, method, timezone, *, ts, eph, sun, moon, earth):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    # ==========================
    # 1. BUILD CONTEXT
    # ==========================
    context = HijriContext(
        lat=lat,
        lon=lon,
        timezone=timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    # ==========================
    # 2. EXECUTE METHOD
    # ==========================
    method_instance = get_method_instance(method)
    result = method_instance.calculate(context)

    today = result.hijri_date
    meta = result.metadata

    # ==========================
    # 3. PREDICTION LOGIC
    # ==========================
    decision = meta.get("decision")
    
    if decision == "no_evaluation_needed":
        return {
            "method": method,
            "today": today,
            "estimated_end_of_month": None,
            "visibility": None,
            "message": "Prediksi akhir bulan hanya tersedia pada fase evaluasi akhir bulan.",
        }

    if decision not in ("new_month", "new_year", "istikmal_30"):
        return {
            "method": method,
            "today": today,
            "estimated_end_of_month": None,
            "visibility": None,
            "message": "Prediksi akhir bulan tidak tersedia.",
        }

    baseline = meta.get("baseline") or today

    year = baseline["year"]
    month = baseline["month"]

    if month == 12:
        next_month = 1
        next_year = year + 1
    else:
        next_month = month + 1
        next_year = year

    estimated_hijri = {
        "year": next_year,
        "month": next_month,
        "day": 1,
    }
    
    if decision in ("new_month", "new_year"):
        offset_days = 1
        msg_suffix = "besok"
    else:
        offset_days = 2
        msg_suffix = "lusa"
    
    estimated_gregorian = now_local + timedelta(days=offset_days)
    
    msg = (
        f"Estimasi 1 {next_month} {next_year} H"
        f"jatuh pada {estimated_gregorian.strftime('%d %B %Y')} ({msg_suffix})"
    )

    return {
        "method": method,
        "today": today,
        "estimated_end_of_month": estimated_hijri,
        "estimated_end_of_month_gregorian": estimated_gregorian.isoformat(),
        "visibility": meta.get("visibility") or meta.get("visibility_data"),
        "message": msg,
    }
