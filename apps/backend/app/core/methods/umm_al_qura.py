from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
)
from ..calendar.hijri_date import increment_hijri_day

from .base import BaseHijriMethod, HijriResult
from ..config import AL_HARAM_LOCATION


class UmmAlQuraMethod(BaseHijriMethod):
    def calculate(self, context):
        baseline, _ = calculate_baseline_hijri(
            context.now_local,
            context.timezone,
        )

        ref_lat, ref_lon, ref_tz = (
            AL_HARAM_LOCATION["lat"],
            AL_HARAM_LOCATION["lon"],
            AL_HARAM_LOCATION["timezone"],
        )

        sunset_local = calculate_sunset(
            context.now_local.date(),
            ref_lat,
            ref_lon,
            ref_tz,
        )

        after_sunset = (
            context.now_local.astimezone(sunset_local.tzinfo) >= sunset_local
            if sunset_local
            else False
        )

        if not after_sunset:
            return HijriResult(
                hijri_date=baseline,
                metadata={
                    "type": "umm_al_qura",
                    "decision": "before_maghrib",
                    "basis": "arithmetic_global",
                },
            )

        final_date = increment_hijri_day(baseline)

        return HijriResult(
            hijri_date=final_date,
            explanation={
                "method": "umm_al_qura",
                "after_sunset": True,
                "criteria_used": "Arithmetic (Umm al-Qura)",
                "reasoning": [
                    "Menggunakan kalender aritmatika resmi Umm al-Qura.",
                    "Tidak mempertimbangkan visibilitas hilal lokal.",
                ],
                "decision": "global_standard",
                "astronomical_data": None,
            },
            metadata={"type": "arithmetic_global"},
        )
