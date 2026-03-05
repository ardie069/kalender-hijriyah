"""
UGHC Decision — single source-of-truth untuk logika keputusan bulan baru UGHC/KHGT.

Butir 1: Hilal terlihat sebelum 24:00 UTC → bulan baru
Butir 2: Hilal terlihat setelah 24:00 UTC:
  - Jika di Amerika → bulan baru
  - Jika bukan Amerika → cek konjungsi < NZ fajr → bulan baru
"""

import pytz
from datetime import datetime, timedelta

from ..calendar.julian import jd_from_datetime, jd_to_datetime
from ..astronomy.sun_times import get_fajr_time
from ..config import NZ_LOCATION
from .visibility_scan import GlobalVisibilityRegistry
from .engine import calculate_conjunction


def evaluate_ughc_decision(date, noon_jd, criteria, lat_step=None, lon_step=None):
    """
    Evaluasi keputusan UGHC bulan baru berdasarkan scan global.

    Args:
        date: tanggal hari ke-29 (date object)
        noon_jd: Julian Date noon untuk konjungsi lookup
        criteria: nama criteria (e.g. "TURKEY_2016_TOPOCENTRIC")
        lat_step/lon_step: override grid step (optional, for predictor)

    Returns:
        (is_new_month: bool, scan_result: dict)
    """
    scan_kwargs = {"date": date, "criteria": criteria}
    if lat_step is not None:
        scan_kwargs["lat_step"] = lat_step
    if lon_step is not None:
        scan_kwargs["lon_step"] = lon_step

    scan = GlobalVisibilityRegistry.scan_global(**scan_kwargs)

    is_new_month = False

    # A. Butir 1: Terpenuhi sebelum 24:00 UTC
    if scan.get("anywhere_before_24utc", False):
        is_new_month = True

    # B. Butir 2: Terpenuhi setelah 24:00 UTC (Amerika atau NZ Fajr)
    elif scan.get("anywhere_after_24utc", False) or scan.get("america_visible", False):
        if scan.get("america_visible", False):
            is_new_month = True
        else:
            # NZ Fajr check — konjungsi harus terjadi sebelum Fajr di NZ
            conj_jd = calculate_conjunction(noon_jd)
            conj_dt_utc = jd_to_datetime(conj_jd)
            # Khusus untuk fajar NZ, kita harus berpatokan pada tanggal lokal di New Zealand
            # saat konjungsi terjadi, karena NZ ada di ujung timur (UTC+12/13).
            nz_tz = pytz.timezone(NZ_LOCATION["timezone"])
            conj_dt_nz = conj_dt_utc.astimezone(nz_tz)

            nz_fajr = get_fajr_time(
                conj_dt_nz.date(),
                NZ_LOCATION["lat"],
                NZ_LOCATION["lon"],
                NZ_LOCATION["timezone"],
            )
            if nz_fajr and conj_dt_utc < nz_fajr.astimezone(pytz.utc):
                is_new_month = True

    return is_new_month, scan
