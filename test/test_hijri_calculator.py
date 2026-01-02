from core.hijri_calculator import get_hijri_date
from core.ephemeris import ts, eph, sun, moon, earth


def test_hijri_date_global():
    result = get_hijri_date(
        lat=21.4225,
        lon=39.8262,
        method="global",
        timezone="Asia/Jakarta",
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    assert "year" in result
    assert "month" in result
    assert "day" in result
