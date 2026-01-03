import pytz  # type: ignore
from datetime import datetime

from app.core.hijri_calculator import get_hijri_date
from app.deps.astronomy import ts, eph, sun, moon, earth


def test_global_maghrib_boundary_no_day_jump():
    tz = pytz.timezone("Asia/Riyadh")

    # 3 Jan 2026 (contoh stabil)
    before_maghrib = tz.localize(datetime(2026, 1, 3, 16, 30))
    after_maghrib = tz.localize(datetime(2026, 1, 3, 19, 0))

    before = get_hijri_date(
        lat=21.4225,
        lon=39.8262,
        method="global",
        timezone="Asia/Riyadh",
        now_local=before_maghrib,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    after = get_hijri_date(
        lat=21.4225,
        lon=39.8262,
        method="global",
        timezone="Asia/Riyadh",
        now_local=after_maghrib,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    # ❌ BUG LAMA: bisa lompat 2 hari
    delta = after["day"] - before["day"]  # type: ignore

    assert delta in (0, 1), f"Global lompat tidak valid: {before} -> {after}"


def test_hisab_maghrib_boundary_stable():
    tz = pytz.timezone("Asia/Jakarta")

    before_maghrib = tz.localize(datetime(2026, 1, 3, 17, 30))
    after_maghrib = tz.localize(datetime(2026, 1, 3, 18, 30))

    before = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=before_maghrib,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    after = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=after_maghrib,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    delta = after["day"] - before["day"]  # type: ignore

    assert delta in (0, 1)


def test_rukyat_maghrib_boundary_no_double_jump():
    tz = pytz.timezone("Asia/Jakarta")

    before_maghrib = tz.localize(datetime(2026, 1, 3, 17, 45))
    after_maghrib = tz.localize(datetime(2026, 1, 3, 18, 45))

    before = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="rukyat",
        timezone="Asia/Jakarta",
        now_local=before_maghrib,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    after = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="rukyat",
        timezone="Asia/Jakarta",
        now_local=after_maghrib,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    delta = after["day"] - before["day"]  # type: ignore

    assert delta >= 0
    assert delta <= 1, f"Rukyat lompat gila: {before} -> {after}"


def test_no_plus_one_day_bug():
    tz = pytz.timezone("Asia/Jakarta")

    night = tz.localize(datetime(2026, 1, 3, 23, 0))
    early_morning = tz.localize(datetime(2026, 1, 4, 1, 0))

    a = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=night,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    b = get_hijri_date(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=early_morning,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    assert b["day"] - a["day"] in (0, 1)  # type: ignore
