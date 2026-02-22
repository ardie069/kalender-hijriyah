from datetime import datetime
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

    if decision == "new_month" or decision == "new_year":
        estimated_date = today
        msg = "Kriteria terpenuhi. Besok diperkirakan masuk bulan baru."
    elif decision == "istikmal_30":
        estimated_date = today
        msg = "Hilal tidak memenuhi kriteria. Bulan digenapkan 30 hari."
    else:
        return {
            "method": method,
            "today": today,
            "estimated_end_of_month": None,
            "visibility": meta.get("visibility") or meta.get("visibility_data"),
            "message": "Prediksi akhir bulan hanya tersedia pada fase evaluasi.",
        }

    return {
        "method": method,
        "today": today,
        "estimated_end_of_month": estimated_date,
        "visibility": meta.get("visibility") or meta.get("visibility_data"),
        "message": msg,
    }
