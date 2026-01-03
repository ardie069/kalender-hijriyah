from skyfield import almanac # type: ignore
from skyfield.api import wgs84 # type: ignore
import pytz # type: ignore


def get_sunset_time(date, lat, lon, timezone, ts, eph):
    observer = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)

    t0 = ts.utc(date.year, date.month, date.day, 0, 0)
    t1 = ts.utc(date.year, date.month, date.day, 23, 59)

    f = almanac.sunrise_sunset(eph, observer)
    times, events = almanac.find_discrete(t0, t1, f)

    for ti, ev in zip(times, events):
        if ev == 0:
            dt_utc = ti.utc_datetime().replace(tzinfo=pytz.utc)
            return dt_utc.astimezone(pytz.timezone(timezone))

    return None
