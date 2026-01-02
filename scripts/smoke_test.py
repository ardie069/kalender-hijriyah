from core.ephemeris import ts, eph, sun, moon, earth
from core.hijri_calculator import get_hijri_date

result = get_hijri_date(
    lat=-7.9666,
    lon=112.6326,
    method="hisab",
    timezone="Asia/Jakarta",
    ts=ts,
    eph=eph,
    sun=sun,
    moon=moon,
    earth=earth,
)

print(result)