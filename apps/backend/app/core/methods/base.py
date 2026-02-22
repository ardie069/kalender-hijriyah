from __future__ import annotations


class HijriContext:
    def __init__(self, lat, lon, timezone, now_local, ts, eph, sun, moon, earth):
        self.lat = lat
        self.lon = lon
        self.timezone = timezone
        self.now_local = now_local
        self.ts = ts
        self.eph = eph
        self.sun = sun
        self.moon = moon
        self.earth = earth


class HijriResult:
    def __init__(self, hijri_date, explanation=None, metadata=None):
        self.hijri_date = hijri_date
        self.explanation = explanation or {}
        self.metadata = metadata or {}


class BaseHijriMethod:
    def calculate(self, context: HijriContext) -> HijriResult:
        raise NotImplementedError
