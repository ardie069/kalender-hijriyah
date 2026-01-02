from datetime import datetime
import pytz

from core.julian import jd_from_datetime
from core.ephemeris import ts


def test_jd_monotonic():
    dt1 = datetime(2026, 1, 3, 0, 0, tzinfo=pytz.utc)
    dt2 = datetime(2026, 1, 3, 1, 0, tzinfo=pytz.utc)

    jd1 = jd_from_datetime(dt1, ts)
    jd2 = jd_from_datetime(dt2, ts)

    assert jd2 > jd1
