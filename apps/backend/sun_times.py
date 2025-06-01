from datetime import datetime
import pytz
from skyfield.api import load, Topos
from skyfield import almanac
from ephemeris import ts, eph


def get_sunset_time(date, latitude, longitude, timezone_name):
    observer = Topos(latitude_degrees=latitude, longitude_degrees=longitude)
    t0 = ts.utc(date.year, date.month, date.day, 0)
    t1 = ts.utc(date.year, date.month, date.day, 23, 59)

    f = almanac.sunrise_sunset(eph, observer)
    times, events = almanac.find_discrete(t0, t1, f)

    for t, event in zip(times, events):
        if event == 0:  # 0 berarti sunset
            sunset_utc = t.utc_datetime().replace(tzinfo=pytz.utc)
            local_timezone = pytz.timezone(timezone_name)
            return sunset_utc.astimezone(local_timezone)

    # Jika tidak ditemukan sunset
    print(f"❌ Tidak menemukan sunset di {latitude}, {longitude} pada {date}")
    return None
