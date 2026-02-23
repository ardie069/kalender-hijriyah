from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    check_historical_lag,
)
from ..calendar.hijri_date import (
    increment_hijri_day,
    decrement_hijri_day,
    start_new_month,
)

from ..services.visibility_scan import GlobalVisibilityRegistry
from ..config import REGIONAL_RUKYAT_CONFIG

from .base import BaseHijriMethod, HijriResult


class LocalRukyatMethod(BaseHijriMethod):
    def __init__(self, mode="individual", region=None):
        self.mode = mode
        self.region = region

    def calculate(self, context):
        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local,
            context.timezone,
            context.ts,
        )

        is_lagging = check_historical_lag(
            baseline,
            noon_jd,
            context.lat,
            context.lon,
            context.timezone,
            context.ts,
            context.eph,
            context.sun,
            context.moon,
            context.earth,
            criteria="MABIMS",
        )

        if is_lagging:
            baseline = decrement_hijri_day(baseline)

        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
            context.ts,
            context.eph,
        )

        after_sunset = context.now_local >= sunset_local if sunset_local else False

        if not after_sunset:
            return HijriResult(
                hijri_date=baseline,
                metadata={
                    "type": "local_rukyat",
                    "decision": "before_maghrib",
                },
            )

        if baseline["day"] != 29:
            final_date = increment_hijri_day(baseline)

            return HijriResult(
                hijri_date=final_date,
                metadata={
                    "type": "local_rukyat",
                    "decision": "normal_increment",
                },
            )

        if self.mode == "individual":
            sites = [
                {
                    "name": "User Location",
                    "lat": context.lat,
                    "lon": context.lon,
                    "timezone": context.timezone,
                }
            ]

            criteria = "MABIMS"

        elif self.mode == "national":
            if not self.region:
                raise ValueError("Region must be provided in national mode")

            region_config = REGIONAL_RUKYAT_CONFIG[self.region]
            sites = region_config["sites"]
            criteria = region_config["criteria"]
        
        result = GlobalVisibilityRegistry.scan_global(
            context.now_local.date(),
            sites,
            criteria,
            context.ts,
            context.eph,
            context.sun,
            context.moon,
            context.earth,
        )

        if result["visible"]:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            final_date = {**baseline, "day": 30}
            decision = "istikmal_30"

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": "local_rukyat",
                "decision": decision,
                "visibility": result["best_visibility"],
            },
        )
