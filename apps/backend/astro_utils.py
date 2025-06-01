import math
from ephemeris import ts, sun, moon

def deg_to_rad(deg):
    return deg * math.pi / 180.0

def rad_to_deg(rad):
    return rad * 180.0 / math.pi

def angular_separation(ra1_deg, dec1_deg, ra2_deg, dec2_deg):
    """
    Menghitung jarak sudut antara dua titik di langit.
    Input dan output dalam derajat.
    """
    ra1 = deg_to_rad(ra1_deg)
    dec1 = deg_to_rad(dec1_deg)
    ra2 = deg_to_rad(ra2_deg)
    dec2 = deg_to_rad(dec2_deg)

    delta_ra = ra2 - ra1
    cos_sep = math.sin(dec1) * math.sin(dec2) + math.cos(dec1) * math.cos(dec2) * math.cos(delta_ra)
    cos_sep = min(1.0, max(-1.0, cos_sep))
    return rad_to_deg(math.acos(cos_sep))

def normalize_angle(angle):
    """
    Menormalkan sudut ke rentang [0, 360] derajat.
    """
    return angle % 360.0

def hours_to_deg(hours):
    """Konversi jam (RA dalam jam) ke derajat"""
    return hours * 15.0

def deg_to_hours(deg):
    """Konversi derajat ke jam (untuk RA)"""
    return deg / 15.0

def get_moon_equatorial(jd):
    """
    Mengambil posisi ekuatorial bulan pada Julian Date tertentu.
    """
    t = ts.tt(jd)  # Waktu dalam Julian Date
    astrometric = (moon.at(t)).apparent()  # Posisi bulan yang dapat diamati
    return astrometric

def get_sun_equatorial(jd):
    """
    Mengambil posisi ekuatorial matahari pada Julian Date tertentu.
    """
    t = ts.tt(jd)  # Waktu dalam Julian Date
    astrometric = (sun.at(t)).apparent()  # Posisi matahari yang dapat diamati
    return astrometric

DEFAULT_LOCATION = {
    'global': (21.422487, 39.826206),  # Koordinat Mekkah
    'zone': 'Asia/Riyadh'  # Zona waktu Mekkah
}