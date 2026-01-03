from datetime import date
import pytz  # type: ignore

from app.core.sun_times import get_sunset_time


def test_get_sunset_time_basic(astro):
    sunset = get_sunset_time(
        date(2026, 1, 3),
        lat=-6.2,
        lon=106.8,
        timezone="Asia/Jakarta",
        ts=astro["ts"],
        eph=astro["eph"],
    )

    assert sunset is not None
    assert sunset.tzinfo is not None
    assert sunset.tzinfo.zone == "Asia/Jakarta"


def test_get_sunset_time_cache_hit(astro):
    d = date(2026, 1, 3)

    a = get_sunset_time(
        d, -6.2, 106.8, "Asia/Jakarta", ts=astro["ts"], eph=astro["eph"]
    )

    b = get_sunset_time(
        d, -6.2, 106.8, "Asia/Jakarta", ts=astro["ts"], eph=astro["eph"]
    )

    assert a == b


def test_get_sunset_time_lat_lon_rounding(astro):
    d = date(2026, 1, 3)

    a = get_sunset_time(
        d, -6.2000001, 106.8000001, "Asia/Jakarta", ts=astro["ts"], eph=astro["eph"]
    )

    b = get_sunset_time(
        d, -6.2, 106.8, "Asia/Jakarta", ts=astro["ts"], eph=astro["eph"]
    )

    assert a == b


def test_get_sunset_time_extreme_latitude(astro):
    sunset = get_sunset_time(
        date(2026, 6, 21),
        lat=69.6,
        lon=18.9,
        timezone="Europe/Oslo",
        ts=astro["ts"],
        eph=astro["eph"],
    )

    assert sunset is None or sunset.tzinfo is not None
