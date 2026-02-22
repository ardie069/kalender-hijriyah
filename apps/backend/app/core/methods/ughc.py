import pytz
from datetime import datetime

from ..julian import jd_from_datetime, julian_to_hijri
from ..conjunction import get_conjunction_time
from ..sun_times import get_sunset_time
from ..visibility import evaluate_visibility

from .base import BaseHijriMethod, HijriResult


class UGHCMethod(BaseHijriMethod):
    def calculate(self, context):
        now_utc = context.now_local.astimezone(pytz.utc)

        end_of_day_utc = datetime(
            now_utc.year, now_utc.month, now_utc.day, 23, 59, 59, tzinfo=pytz.utc
        )

        candidate_longitudes = [-150, -120, -90, -60, -30, 0]
        candidate_latitudes = [-40, -20, 0, 20, 40]

        global_visible = False
        visibility_data = {}

        best_score = -999

        for lat in candidate_latitudes:
            for lon in candidate_longitudes:

                sunset_local = get_sunset_time(
                    context.now_local.date(), lat, lon, "UTC", context.ts, context.eph
                )

                if not sunset_local:
                    continue

                sunset_utc = sunset_local.astimezone(pytz.utc)

                if sunset_utc > end_of_day_utc:
                    continue

                sunset_jd = jd_from_datetime(sunset_utc, context.ts)

                conj_jd = get_conjunction_time(
                    sunset_jd,
                    context.ts,
                    context.earth,
                    context.sun,
                    context.moon,
                )

                moon_age_hours = (sunset_jd - conj_jd) * 24
                conjunction_before_sunset = moon_age_hours > 0

                vis = evaluate_visibility(
                    sunset_utc,
                    lat,
                    lon,
                    conj_jd,
                    context.ts,
                    context.sun,
                    context.moon,
                    context.earth,
                    criteria="TURKEY_2016_GEOCENTRIC",
                )

                score = vis["moon_altitude"] + vis["elongation"]

                if score > best_score:
                    best_score = score
                    visibility_data = {
                        **vis,
                        "lat": lat,
                        "lon": lon,
                        "moon_age_hours": moon_age_hours,
                        "conjunction_before_sunset": conjunction_before_sunset,
                    }

                if (
                    conjunction_before_sunset
                    and vis["moon_altitude"] >= 5.0
                    and vis["elongation"] >= 8.0
                ):
                    global_visible = True
                    break

            if global_visible:
                break

        noon_utc = datetime(
            now_utc.year,
            now_utc.month,
            now_utc.day,
            12,
            0,
            0,
            tzinfo=pytz.utc,
        )

        noon_jd = jd_from_datetime(noon_utc, context.ts)
        baseline = julian_to_hijri(noon_jd)

        year = baseline["year"]
        month = baseline["month"]
        day = baseline["day"]

        decision = "no_evaluation_needed"

        if day < 29:
            final_date = baseline
        elif day == 29:
            if global_visible:
                if month == 12:
                    decision = "new_year"
                    final_date = {
                        "year": year + 1,
                        "month": 1,
                        "day": 1,
                    }
                else:
                    decision = "new_month"
                    final_date = {
                        "year": year,
                        "month": month + 1,
                        "day": 1,
                    }
            else:
                decision = "istikmal_30"
                final_date = baseline
        else:
            if month == 12:
                decision = "new_year"
                final_date = {
                    "year": year + 1,
                    "month": 1,
                    "day": 1,
                }
            else:
                decision = "new_month"
                final_date = {
                    "year": year,
                    "month": month + 1,
                    "day": 1,
                }

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": "ughc_geocentric",
                "global_visible": global_visible,
                "visibility_data": visibility_data,
                "decision": decision,
                "baseline": baseline,
            },
        )
