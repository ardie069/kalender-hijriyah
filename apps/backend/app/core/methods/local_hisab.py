from skyfield.api import wgs84
from skyfield import almanac
from datetime import timedelta
import pytz

from ..services.engine import (
    calculate_conjunction,
    calculate_visibility,
)
from ..calendar.hijri_date import start_new_month
from ..calendar.julian import jd_from_datetime
from ._shared import prepare_hijri_baseline, handle_before_sunset, handle_normal_day
from .base import BaseHijriMethod, HijriResult
from app.deps.astronomy import get_provider


class LocalHisabMethod(BaseHijriMethod):
    def calculate(self, context):
        # 1. Baseline & Status Maghrib
        baseline, after_sunset, sunset_local = prepare_hijri_baseline(
            context,
            criteria="Wujudul Hilal",
        )

        if not after_sunset:
            return handle_before_sunset(baseline, "local_hisab")

        if baseline["day"] != 29:
            return handle_normal_day(baseline, "local_hisab")

        # 2. Evaluasi Hari ke-29 (Wujudul Hilal)
        p = get_provider()
        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc)

        # A. Cek Ijtima' (Konjungsi)
        conj_jd = calculate_conjunction(sunset_jd)
        is_conj_before_sunset = conj_jd < sunset_jd

        # B. Cek Selisih Terbenam (Moonset - Sunset)
        # Scan 2 jam di sekitar sunset untuk mencari waktu bulan terbenam
        t0 = p.ts.from_datetime(sunset_utc - timedelta(hours=1))
        t1 = p.ts.from_datetime(sunset_utc + timedelta(hours=2))

        f = almanac.risings_and_settings(
            p.eph, p.moon, wgs84.latlon(context.lat, context.lon)
        )
        times, events = almanac.find_discrete(t0, t1, f)

        moonset_jd = None
        for t, event in zip(times, events):
            if event == 0:  # 0 = Set (Terbenam)
                moonset_jd = t.tt
                break

        # C. Hitung Lag Time
        if moonset_jd:
            # Selisih dalam menit: (JD_moonset - JD_sunset) * 24 jam * 60 menit
            lag_time_min = round((moonset_jd - sunset_jd) * 1440, 2)
            is_moonset_after_sunset = moonset_jd > sunset_jd
        else:
            # Fallback ke altitude check jika moonset tidak ditemukan di range scan
            vis = calculate_visibility(
                sunset_utc, context.lat, context.lon, conj_jd, criteria="Wujudul Hilal"
            )
            is_moonset_after_sunset = vis["moon_altitude"] > 0
            lag_time_min = 0  # Tidak terdeteksi terbenam di range ini

        # 3. Keputusan Final (Wujudul Hilal)
        # Syarat: Ijtima' sudah terjadi DAN Bulan terbenam setelah Matahari
        if is_conj_before_sunset and is_moonset_after_sunset:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            final_date = {**baseline, "day": 30}  # Istikmal
            decision = "istikmal_30"

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": "local_hisab",
                "decision": decision,
                "lag_time_minutes": lag_time_min,
                "is_conjunction_before_sunset": is_conj_before_sunset,
                "is_moonset_after_sunset": is_moonset_after_sunset,
                "site_name": "User Local Coordinates",
            },
        )
