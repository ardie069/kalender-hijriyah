# app/core/services/month_predictor.py
import pytz
from datetime import datetime, timedelta, date
from app.core.methods.factory import get_method_instance

HIJRI_MONTH_NAMES = [
    "", "Muharam", "Safar", "Rabiulawal", "Rabiulakhir",
    "Jumadilawal", "Jumadilakhir", "Rajab", "Syakban",
    "Ramadan", "Syawal", "Zulkaidah", "Zulhijah",
]


class MonthPredictor:
    def __init__(self, factory):
        self.factory = factory

    def predict_full_year(self, hijri_year: int, lat, lon, timezone, method):
        method_instance = get_method_instance(method)
        calendar_data = []

        # Estimasi 1 Muharram
        current_gregorian = self._estimate_start_of_hijri_year(
            hijri_year, lat, lon, timezone, method
        )

        for m_idx in range(1, 13):
            day_29_gregorian = current_gregorian + timedelta(days=28)
            # PANGGIL FACTORY DI SINI
            context = self.factory.create_context(lat, lon, timezone, day_29_gregorian)

            result = method_instance.calculate(context)
            meta = result.metadata

            total_days = 30 if meta.get("decision") == "istikmal_30" else 29

            calendar_data.append(
                {
                    "month_id": m_idx,
                    "month_name": HIJRI_MONTH_NAMES[m_idx],
                    "total_days": total_days,
                    "start_gregorian": current_gregorian.isoformat(),
                    "day_1_weekday": current_gregorian.weekday(),
                    "visibility": meta.get("visibility"),
                }
            )
            current_gregorian += timedelta(days=total_days)

        return calendar_data

    def _estimate_start_of_hijri_year(self, year, lat, lon, tz, method):
        # Scan +/- 15 hari dari estimasi rata-rata
        total_days_approx = int((year - 1) * 354.367)
        base_date = date(622, 7, 19) + timedelta(days=total_days_approx)
        method_instance = get_method_instance(method)
        timezone = pytz.timezone(tz)

        current_test = base_date - timedelta(days=10)
        for _ in range(30):
            test_dt = datetime.combine(current_test, datetime.min.time()).replace(
                tzinfo=timezone
            )
            ctx = self.factory.create_context(lat, lon, tz, test_dt)
            res = method_instance.calculate(ctx)
            if res.hijri_date["month"] == 1 and res.hijri_date["day"] == 1:
                # Return datetime, bukan date, supaya kompatibel dengan create_context
                return datetime.combine(current_test, datetime.min.time()).replace(
                    tzinfo=timezone
                )
            current_test += timedelta(days=1)
        return datetime.combine(base_date, datetime.min.time()).replace(
            tzinfo=timezone
        )


# --- STANDALONE FUNCTION BUAT API ---
def predict_end_of_month(lat, lon, method, timezone, context_factory):
    # Kita panggil lewat method predictor
    predictor = MonthPredictor(context_factory)

    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    # Logic prediksi satu bulan
    method_instance = get_method_instance(method)
    context_today = context_factory.create_context(lat, lon, timezone, now_local)
    result_today = method_instance.calculate(context_today)
    hijri_now = result_today.hijri_date

    days_to_29 = 29 - hijri_now["day"]
    simulation_local = now_local + timedelta(days=days_to_29)

    context_sim = context_factory.create_context(lat, lon, timezone, simulation_local)
    result_sim = method_instance.calculate(context_sim)

    is_istikmal = result_sim.metadata.get("decision") == "istikmal_30"

    # ... return dictionary lu yang lama ...
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
