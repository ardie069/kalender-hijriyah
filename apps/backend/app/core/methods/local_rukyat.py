from ..astronomy_engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
    check_historical_lag,
    increment_hijri_day,
    decrement_hijri_day,
    start_new_month,
)

from ..julian import jd_from_datetime
import pytz

from .base import BaseHijriMethod, HijriResult


class LocalRukyatMethod(BaseHijriMethod):
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

        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc, context.ts)

        conj_jd = calculate_conjunction(
            sunset_jd,
            context.ts,
            context.earth,
            context.sun,
            context.moon,
        )

        vis = calculate_visibility(
            sunset_utc,
            context.lat,
            context.lon,
            conj_jd,
            context.ts,
            context.sun,
            context.moon,
            context.earth,
            criteria="MABIMS",
        )

        if vis["is_visible"]:
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
                "visibility": vis,
            },
        )
