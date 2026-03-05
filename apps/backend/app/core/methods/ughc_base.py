
from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
)
from ..calendar.hijri_date import (
    increment_hijri_day,
    start_new_month,
)
from ..services.ughc_decision import evaluate_ughc_decision
from .base import BaseHijriMethod, HijriResult


class BaseUGHCMethod(BaseHijriMethod):
    CRITERIA = None
    TYPE = None

    def calculate(self, context):
        # 1. Ambil baseline hanya untuk struktur datanya (Year, Month)
        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local,
            context.timezone,
        )

        # 2. Status Maghrib Lokal
        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
        )
        after_sunset = context.now_local >= sunset_local if sunset_local else False

        if not after_sunset:
            return HijriResult(
                hijri_date=baseline, metadata={"decision": "before_maghrib"}
            )

        # 3. Hari biasa (bukan hari ke-29/30) → langsung increment, TANPA scan global
        if baseline["day"] not in (29, 30):
            return HijriResult(
                hijri_date=increment_hijri_day(baseline),
                metadata={
                    "type": self.TYPE,
                    "decision": "normal_increment",
                },
            )

        # 4. Hari ke-29/30: evaluasi scan global via shared function
        is_new_month, scan = evaluate_ughc_decision(
            context.now_local.date(), noon_jd, self.CRITERIA
        )

        if is_new_month:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            if baseline["day"] == 29:
                final_date = {**baseline, "day": 30}
                decision = "istikmal_30"
            else:
                final_date = increment_hijri_day(baseline)
                decision = "normal_increment"

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": self.TYPE,
                "decision": decision,
                "global_visible": is_new_month,
                "visibility_data": scan.get("best_visibility"),
            },
        )
