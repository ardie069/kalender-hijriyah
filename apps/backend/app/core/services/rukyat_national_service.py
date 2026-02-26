import pytz
from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)
from ..calendar.julian import jd_from_datetime
from ..config import REGIONAL_RUKYAT_CONFIG


class RukyatNationalService:

    def evaluate(self, context, region: str):

        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local,
            context.timezone,
            context.ts,
        )

        if baseline["day"] != 29:
            return {
                "is_rukyat_day": False,
                "hijri_date": baseline,
            }

        region_config = REGIONAL_RUKYAT_CONFIG.get(region)

        if not region_config:
            raise ValueError("Invalid region")

        criteria = region_config["criteria"]
        sites = region_config["sites"]

        any_visible = False
        best_visibility = None
        best_score = -999

        for site in sites:

            sunset_local = calculate_sunset(
                context.now_local.date(),
                site["lat"],
                site["lon"],
                site["timezone"],
                context.ts,
                context.eph,
            )

            if not sunset_local:
                continue

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
                site["lat"],
                site["lon"],
                conj_jd,
                context.ts,
                context.sun,
                context.moon,
                context.earth,
                criteria=criteria,
            )

            score = vis["moon_altitude"] + vis["elongation"]

            if score > best_score:
                best_score = score
                best_visibility = {
                    **vis,
                    "site": site["name"],
                }

            if vis["is_visible"]:
                any_visible = True

        return {
            "is_rukyat_day": True,
            "is_visible_national": any_visible,
            "criteria_used": criteria,
            "best_visibility": best_visibility,
        }
