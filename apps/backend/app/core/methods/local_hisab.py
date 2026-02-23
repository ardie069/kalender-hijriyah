from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
    check_historical_lag,
)
from ..calendar.hijri_date import (
    increment_hijri_day,
    decrement_hijri_day,
    start_new_month,
)

from ..calendar.julian import jd_from_datetime
import pytz

from .base import BaseHijriMethod, HijriResult


class LocalHisabMethod(BaseHijriMethod):
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
            criteria="Wujudul Hilal",
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
                    "type": "local_hisab",
                    "decision": "before_maghrib",
                },
            )

        if baseline["day"] != 29:
            final_date = increment_hijri_day(baseline)

            return HijriResult(
                hijri_date=final_date,
                metadata={"type": "local_hisab", "decision": "normal_increment"},
            )

        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc, context.ts)

        conj_jd = calculate_conjunction(
            sunset_jd, context.ts, context.earth, context.sun, context.moon
        )

        if conj_jd < sunset_jd:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            final_date = {**baseline, "day": 30}
            decision = "istikmal_30"

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": "local_hisab",
                "decision": decision,
                "conjunction_before_sunset": conj_jd < sunset_jd,
            },
        )
