import pytz

from .engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)

from ..calendar.julian import jd_from_datetime
from ..config import REGIONAL_RUKYAT_CONFIG


class RukyatService:
    def evaluate_local_rukyat(self, context, region: str = "indonesia"):
        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local, context.timezone, context.ts
        )

        if baseline["day"] != 29:
            return {"is_rukyat_day": False, "hijri_date": baseline}

        region_config = REGIONAL_RUKYAT_CONFIG[region]

        if not region_config:
            raise ValueError("Invalid region")

        criteria = region_config["criteria"]

        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
            context.ts,
            context.eph,
        )
        if not sunset_local:
            return {"is_rukyat_day": True, "error": "Sunset culd not be calculated"}

        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc, context.ts)

        conj_jd = calculate_conjunction(
            sunset_jd, context.ts, context.earth, context.sun, context.moon
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
            criteria=criteria,
        )

        return {
            "is_rukyat_day": True,
            "sunset_time": sunset_local.isoformat(),
            "altitude_at_sunset": vis["moon_altitude"],
            "elongation_at_sunset": vis["elongation"],
            "is_visible": vis["is_visible"],
            "criteria_used": criteria,
            "moon_age_hours": vis["moon_age_hours"],
        }
