from skyfield.api import load, wgs84
from skyfield.almanac import fraction_illuminated
from skyfield import almanac
import numpy as np
from datetime import timedelta


class SkyfieldAdapter:
    def __init__(self, ephemeris_path="de421.bsp"):
        # Load ephemeris & timescale (Materi Utama)
        self.ts = load.timescale()
        self.eph = load(ephemeris_path)
        self.earth = self.eph["earth"]
        self.moon = self.eph["moon"]
        self.sun = self.eph["sun"]

    def get_observer(self, lat: float, lon: float, elevation: float = 0):
        """Membangun titik pengamatan di kerak bumi (WGS84)"""
        return self.earth + wgs84.latlon(lat, lon, elevation_m=elevation)

    def get_moon_telemetry(self, dt_utc, lat: float, lon: float):
        """
        Menghasilkan telemetry bulan toposentris (Real-world position)
        """
        t = self.ts.from_datetime(dt_utc)
        observer = self.get_observer(lat, lon)

        # 1. Posisi Toposentris (Altitude & Azimuth)
        # Dialektika: Dari pusat bumi ke mata pengamat
        astrometric = observer.at(t).observe(self.moon)
        apparent = astrometric.apparent()
        alt, az, distance = apparent.altaz()

        # 2. Iluminasi (Persentase Cahaya)
        illum = fraction_illuminated(self.eph, "moon", t) * 100

        # 3. Elongasi (Sudut Pisah Matahari-Bulan)
        # Penting buat kriteria MABIMS
        sun_pos = observer.at(t).observe(self.sun).apparent()
        elongation = apparent.separation_from(sun_pos).degrees

        return {
            "altitude": alt.degrees,
            "azimuth": az.degrees,
            "distance_km": distance.km,
            "illumination": illum,
            "elongation": elongation,
            "timestamp": dt_utc,
        }

    def get_moon_phase_angle(self, dt_utc):
        """Menghitung phase angle untuk menentukan Waxing/Waning"""
        t = self.ts.from_datetime(dt_utc)
        # Logika: Sudut fase dari pusat matahari-bulan-bumi
        s = self.sun.at(t)
        m = self.moon.at(t)
        e = self.earth.at(t)

        # Vektor posisi
        sm = m.position.au - s.position.au
        em = m.position.au - e.position.au

        # Dot product untuk cari sudut
        cos_angle = np.dot(sm, em) / (np.linalg.norm(sm) * np.linalg.norm(em))
        angle = np.degrees(np.arccos(np.clip(cos_angle, -1.0, 1.0)))

        return angle

    def get_last_conjunction(self, dt_utc):
        """Mencari waktu New Moon terdekat sebelum dt_utc"""
        t1 = self.ts.from_datetime(dt_utc - timedelta(days=31))
        t2 = self.ts.from_datetime(dt_utc)

        # Cari fase bulan di rentang 31 hari ke belakang
        f = almanac.moon_phases(self.eph)
        times, phases = almanac.find_discrete(t1, t2, f)

        # Cari yang fasenya 0 (New Moon)
        new_moons = [t for t, p in zip(times, phases) if p == 0]

        if new_moons:
            return new_moons[-1].utc_datetime()  # Ambil yang paling terakhir
        return dt_utc - timedelta(days=1)  # Fallback jika gagal
