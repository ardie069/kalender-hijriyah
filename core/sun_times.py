# sun_times.py
from skyfield import almanac
from skyfield.api import load, wgs84
import pytz

# Load ephemeris bawaan JPL
ts = load.timescale()
eph = load("de421.bsp")  # bisa diganti "de440s.bsp" untuk lebih presisi

def get_sunset_time(date, lat, lon, timezone):
    """
    Menghitung waktu matahari terbenam (sunset) untuk lat/lon pada tanggal tertentu.
    Return: datetime dengan timezone lokal.
    """
    # Lokasi observer di permukaan bumi (pakai WGS84, bukan .topos())
    observer = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)

    # Rentang waktu satu hari
    t0 = ts.utc(date.year, date.month, date.day, 0, 0)
    t1 = ts.utc(date.year, date.month, date.day, 23, 59)

    # Sunrise/Sunset function
    f = almanac.sunrise_sunset(eph, observer)

    times, events = almanac.find_discrete(t0, t1, f)

    # Cari event sunset (event=0 → sunset, event=1 → sunrise)
    for ti, ev in zip(times, events):
        if ev == 0:  # sunset
            dt_utc = ti.utc_datetime().replace(tzinfo=pytz.utc)
            return dt_utc.astimezone(pytz.timezone(timezone))

    return None