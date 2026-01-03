from datetime import datetime
import pytz  # type: ignore

from core.hijri_calculator import get_hijri_date


def test_global_hijri_date_fixed_time(astro):
    now = datetime(2026, 1, 3, 1, 0, tzinfo=pytz.utc)

    result = get_hijri_date(
        lat=21.4225,
        lon=39.8262,
        method="global",
        timezone="Asia/Riyadh",
        now_local=now,  # type: ignore
        **astro,
    )

    assert "year" in result
    assert "month" in result
    assert "day" in result


def test_hisab_day_29_decision_path(astro):
    now = datetime(2026, 1, 18, 18, 0, tzinfo=pytz.utc)

    result = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=now, # type: ignore
        **astro,
    )

    assert result["day"] in (29, 30, 1)


def test_hisab_vs_rukyat_can_differ(astro):
    now = datetime(2026, 1, 3, 16, 0, tzinfo=pytz.utc)

    hisab = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=now,
        **astro,
    )

    rukyat = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="rukyat",
        timezone="Asia/Jakarta",
        now_local=now,
        **astro,
    )

    assert hisab["day"] >= rukyat["day"] # type: ignore

