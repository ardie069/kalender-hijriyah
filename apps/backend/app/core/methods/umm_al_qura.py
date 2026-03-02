from skyfield.api import wgs84
from skyfield import almanac
from datetime import timedelta
import pytz

from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
)
from ..calendar.hijri_date import increment_hijri_day, start_new_month
from ..calendar.julian import jd_from_datetime
from ..config import AL_HARAM_LOCATION
from app.deps.astronomy import get_provider
from .base import BaseHijriMethod, HijriResult


class UmmAlQuraMethod(BaseHijriMethod):
    def calculate(self, context):
        # 1. Baseline Aritmatika
        baseline, _ = calculate_baseline_hijri(context.now_local, context.timezone)

        ref_lat, ref_lon, ref_tz = (
            AL_HARAM_LOCATION["lat"],
            AL_HARAM_LOCATION["lon"],
            AL_HARAM_LOCATION["timezone"],
        )
        makkah_tz = pytz.timezone(ref_tz)
        local_now_makkah = context.now_local.astimezone(makkah_tz)

        # 2. Cari Sunset Makkah
        sunset_makkah = calculate_sunset(
            local_now_makkah.date(), ref_lat, ref_lon, ref_tz
        )

        if not sunset_makkah or local_now_makkah < sunset_makkah:
            return HijriResult(
                hijri_date=baseline,
                metadata={"type": "umm_al_qura", "decision": "before_maghrib"},
            )

        # 3. Logika Evaluasi
        if baseline["day"] == 29:
            p = get_provider()
            sunset_utc = sunset_makkah.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc)

            # A. Cek Ijtima'
            conj_jd = calculate_conjunction(sunset_jd)
            is_conj_before_sunset = conj_jd < sunset_jd

            # B. Cek Moonset vs Sunset (Umm Al-Qura Standard)
            t0 = p.ts.from_datetime(sunset_utc - timedelta(hours=1))
            t1 = p.ts.from_datetime(sunset_utc + timedelta(hours=1))

            # Pake search discrete buat nyari kejadian moonset yang presisi
            f = almanac.risings_and_settings(
                p.eph, p.moon, wgs84.latlon(ref_lat, ref_lon)
            )
            times, events = almanac.find_discrete(t0, t1, f)

            moonset_jd = None
            for t, event in zip(times, events):
                if event == 0:  # 0 = Set (Terbenam)
                    moonset_jd = t.tt
                    break

            # Jika moonset tidak terdeteksi di range 2 jam, pake altitude check (safety)
            if moonset_jd is None:
                t_sunset = p.ts.from_datetime(sunset_utc)
                observer = p.earth + wgs84.latlon(ref_lat, ref_lon)
                alt, _, _ = observer.at(t_sunset).observe(p.moon).apparent().altaz()
                is_moonset_after_sunset = (
                    alt.degrees + 0.25
                ) > 0  # Koreksi semidiameter
            else:
                is_moonset_after_sunset = moonset_jd > sunset_jd

            # Keputusan Final
            if is_conj_before_sunset and is_moonset_after_sunset:
                final_date = start_new_month(baseline)
                decision = "new_month"
            else:
                final_date = {**baseline, "day": 30}
                decision = "istikmal_30"

            metadata = {
                "type": "umm_al_qura",
                "decision": decision,
                "is_conjunction_before_sunset": is_conj_before_sunset,
                "is_moonset_after_sunset": is_moonset_after_sunset,
            }
        else:
            # Hari Biasa
            final_date = increment_hijri_day(baseline)
            metadata = {"type": "umm_al_qura", "decision": "normal_increment"}

        return HijriResult(hijri_date=final_date, metadata=metadata)
