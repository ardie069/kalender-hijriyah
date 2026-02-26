from ..services.engine import (
    calculate_conjunction,
    calculate_visibility,
)
from ..calendar.hijri_date import start_new_month
from ..calendar.julian import jd_from_datetime

from ._shared import prepare_hijri_baseline, handle_before_sunset, handle_normal_day
from .base import BaseHijriMethod, HijriResult

import pytz


class LocalHisabMethod(BaseHijriMethod):
    def calculate(self, context):
        baseline, after_sunset, sunset_local = prepare_hijri_baseline(
            context, criteria="Wujudul Hilal",
        )

        if not after_sunset:
            return handle_before_sunset(baseline, "local_hisab")

        if baseline["day"] != 29:
            return handle_normal_day(baseline, "local_hisab")

        # Evaluasi hari ke-29: cek konjungsi vs sunset
        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc)

        conj_jd = calculate_conjunction(sunset_jd)

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
