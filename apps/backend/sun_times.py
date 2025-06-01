# sun_times.py
from datetime import datetime
from skyfield.api import load, wgs84
from ephemeris import ts, sun
import pytz

def get_sunset_time(date, lat, lon, timezone):
    """
    Menghitung waktu matahari terbenam (sunset) pada tanggal tertentu di lokasi yang diberikan.
    
    :param date: Tanggal untuk perhitungan (datetime.date)
    :param lat: Lintang lokasi (derajat)
    :param lon: Bujur lokasi (derajat)
    :param timezone: Zona waktu tempat (string, contoh: 'Asia/Riyadh')
    :return: Waktu matahari terbenam sebagai objek datetime dengan zona waktu
    """
    observer = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)
    
    # Waktu yang dimulai pada tengah malam UTC
    t = ts.utc(date.year, date.month, date.day)
    
    # Menghitung posisi matahari pada tanggal yang diberikan
    astrometric = (sun.at(t)).apparent()
    
    # Menghitung waktu matahari terbenam
    alt, az, d = observer.at(t).observe(sun).apparent().altaz()
    sunset_time = t + d
    
    # Mengubah ke waktu lokal
    sunset_time_utc = datetime.utcfromtimestamp(sunset_time.utc).replace(tzinfo=pytz.utc)
    local_tz = pytz.timezone(timezone)
    sunset_local = sunset_time_utc.astimezone(local_tz)
    
    return sunset_local
