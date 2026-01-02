from datetime import datetime
import pytz

from core.visibility import evaluate_visibility
from core.ephemeris import ts, sun, moon, earth


def test_visibility_structure():
    sunset_dt_utc = datetime(2026, 1, 3, 10, 0, tzinfo=pytz.utc)
    conjunction_jd = 2460680.0

    result = evaluate_visibility(
        sunset_dt_utc=sunset_dt_utc,
        lat=21.4225,
        lon=39.8262,
        conjunction_jd=conjunction_jd,
        ts=ts,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    assert "is_visible" in result
    assert "moon_altitude" in result
    assert "elongation" in result
    assert "moon_age" in result
