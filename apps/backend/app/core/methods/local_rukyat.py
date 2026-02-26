from ..services.visibility_scan import GlobalVisibilityRegistry
from ..calendar.hijri_date import start_new_month
from ..config import REGIONAL_RUKYAT_CONFIG

from ._shared import prepare_hijri_baseline, handle_before_sunset, handle_normal_day
from .base import BaseHijriMethod, HijriResult


class LocalRukyatMethod(BaseHijriMethod):
    def __init__(self, mode="individual", region=None):
        self.mode = mode
        self.region = region

    def calculate(self, context):
        baseline, after_sunset, sunset_local = prepare_hijri_baseline(
            context, criteria="MABIMS",
        )

        if not after_sunset:
            return handle_before_sunset(baseline, "local_rukyat")

        if baseline["day"] != 29:
            return handle_normal_day(baseline, "local_rukyat")

        # Evaluasi hari ke-29
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
            criteria,
        )

        if result["visible"] if "visible" in result else result.get("anywhere_before_24utc", False):
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
                "visibility": result.get("best_visibility"),
            },
        )
