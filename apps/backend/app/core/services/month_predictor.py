import pytz
from datetime import datetime, timedelta, date
from app.core.methods.factory import get_method_instance
from ..config import AL_HARAM_LOCATION, HIJRI_MONTHS

from app.core.services.engine import (
    calculate_sunset,
    calculate_baseline_hijri,
)


class MonthPredictor:
    def __init__(self, factory):
        self.factory = factory

    def predict_full_year(self, hijri_year: int, lat, lon, timezone, method):
        # ── OVERRIDE LOKASI UNTUK UMM AL-QURA ──
        # Jika metodenya Umm Al-Qura, abaikan lat/lon user, pakai Makkah.
        if method == "umm_al_qura":
            lat = AL_HARAM_LOCATION["lat"]
            lon = AL_HARAM_LOCATION["lon"]
            timezone = AL_HARAM_LOCATION["timezone"]

        method_instance = get_method_instance(method)
        calendar_data = []

        current_maghrib = self._estimate_start_of_hijri_year(
            hijri_year, lat, lon, timezone, method
        )

        for m_idx in range(1, 13):
            day_29_date = current_maghrib.date() + timedelta(days=28)
            maghrib_29 = calculate_sunset(day_29_date, lat, lon, timezone)

            if not maghrib_29:
                total_days = 30
                visibility_data = None
            else:
                sim_time = maghrib_29
                context = self.factory.create_context(lat, lon, timezone, sim_time)
                result = method_instance.calculate(context)

                # Logika ganti bulan yang method-agnostic
                result_month = result.hijri_date.get("month", m_idx)
                total_days = 29 if result_month != m_idx else 30

                # Gunakan metadata ughc atau visibility standar
                visibility_data = result.metadata.get(
                    "visibility"
                ) or result.metadata.get("visibility_data")

            calendar_data.append(
                {
                    "month_id": m_idx,
                    "month_name": HIJRI_MONTHS[m_idx],
                    "total_days": total_days,
                    "start_gregorian": current_maghrib.date().isoformat(),
                    "day_1_weekday": current_maghrib.weekday(),
                    "visibility": visibility_data,
                    "decision": (
                        result.metadata.get("decision")
                        if "result" in locals()
                        else "fallback"
                    ),
                }
            )

            current_month_days = total_days
            next_month_date = current_maghrib.date() + timedelta(
                days=current_month_days
            )
            current_maghrib = calculate_sunset(next_month_date, lat, lon, timezone)

        return calendar_data

    def _estimate_start_of_hijri_year(self, year, lat, lon, tz, method):
        """Mencari Maghrib yang menandai 1 Muharram."""
        days_approx = int((year - 1) * 354.36707)
        base_date = date(622, 7, 19) + timedelta(days=days_approx)

        # Scan +/- 10 hari
        scan_start = base_date - timedelta(days=10)
        for i in range(20):
            test_date = scan_start + timedelta(days=i)
            sunset = calculate_sunset(test_date, lat, lon, tz)
            if not sunset:
                continue

            h_date, _ = calculate_baseline_hijri(sunset, tz, lat, lon)

            if h_date["month"] == 1 and h_date["day"] == 1:
                return sunset

        # Fallback ke Maghrib base_date
        return calculate_sunset(base_date, lat, lon, tz)


# --- STANDALONE FUNCTION BUAT API ---
def predict_end_of_month(lat, lon, method, timezone, context_factory):
    """
    Prediksi akhir bulan Hijriyah (istikmal atau bulan baru).

    Sudah benar: menggunakan method_instance.calculate() untuk delegasi.
    """
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    method_instance = get_method_instance(method)
    context_today = context_factory.create_context(lat, lon, timezone, now_local)
    result_today = method_instance.calculate(context_today)
    hijri_now = result_today.hijri_date

    days_to_29 = 29 - hijri_now["day"]
    simulation_local = now_local + timedelta(days=days_to_29)

    context_sim = context_factory.create_context(lat, lon, timezone, simulation_local)
    result_sim = method_instance.calculate(context_sim)

    is_istikmal = result_sim.metadata.get("decision") == "istikmal_30"

    return {
        "current_hijri": hijri_now,
        "prediction": {
            "is_istikmal": is_istikmal,
            "estimated_next_month_1": {
                "year": (
                    hijri_now["year"]
                    if hijri_now["month"] < 12
                    else hijri_now["year"] + 1
                ),
                "month": hijri_now["month"] + 1 if hijri_now["month"] < 12 else 1,
                "day": 1,
            },
        },
        "visibility_at_29": result_sim.metadata.get("visibility"),
        "message": "Prediksi berhasil",
    }
