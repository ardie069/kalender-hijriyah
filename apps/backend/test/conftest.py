import pytest
import pytz
from datetime import datetime
from app.deps.astronomy import ts, eph, sun, moon, earth


@pytest.fixture(scope="session")
def astro():
    return {
        "ts": ts,
        "eph": eph,
        "sun": sun,
        "moon": moon,
        "earth": earth,
    }


def dt(y, m, d, h, tz):
    return pytz.timezone(tz).localize(datetime(y, m, d, h, 0))


def test_global_maghrib_boundary(astro):
    from app.core.hijri_calculator import get_hijri_date

    now = dt(2026, 1, 18, 19, "Asia/Riyadh")  # malam
    h = get_hijri_date(
        21.4225, 39.8262, "global", "Asia/Riyadh", now_local=now, **astro
    )

    assert h["day"] in (29, 30)
    assert h["month"] == 7  # Rajab


def test_hisab_istikmal_on_29th(astro):
    from app.core.hijri_calculator import get_hijri_date

    now = dt(2026, 1, 19, 5, "Asia/Jakarta")
    h = get_hijri_date(-6.2, 106.8, "hisab", "Asia/Jakarta", now_local=now, **astro)

    assert h["day"] == 30
    assert h["month"] == 7  # Rajab


def test_rukyat_not_29_no_evaluation(astro):
    from app.core.hijri_calculator import get_hijri_date

    now = dt(2026, 1, 18, 17, "Asia/Jakarta")
    h = get_hijri_date(-6.2, 106.8, "rukyat", "Asia/Jakarta", now_local=now, **astro)

    assert h["day"] != 29  # baru masuk 29
    assert h["month"] == 7


@pytest.mark.parametrize("method", ["global", "hisab", "rukyat"])
def test_all_methods_converge_on_1_syaban(astro, method):
    from app.core.hijri_calculator import get_hijri_date

    now = dt(2026, 1, 20, 10, "Asia/Jakarta")
    h = get_hijri_date(-6.2, 106.8, method, "Asia/Jakarta", now_local=now, **astro)

    assert h["day"] == 1
    assert h["month"] == 8  # Sya’ban


def test_no_day_jump_over_1(astro):
    from app.core.hijri_calculator import get_hijri_date

    dates = []
    for d in [18, 19, 20]:
        h = get_hijri_date(
            -6.2,
            106.8,
            "rukyat",
            "Asia/Jakarta",
            now_local=dt(2026, 1, d, 19, "Asia/Jakarta"),
            **astro,
        )
        dates.append(h)

    days = [x["day"] for x in dates]
    assert days != [29, 1, 2]  # forbidden


def test_midnight_transition_consistency(astro):
    """Memastikan tidak ada perubahan tanggal Hijriyah antara jam 23:59 dan 00:01"""
    from app.core.hijri_calculator import get_hijri_date

    # 1. Cek sebelum tengah malam (23:55)
    before_midnight = dt(2026, 1, 18, 23, "Asia/Jakarta")
    h_before = get_hijri_date(
        -6.2, 106.8, "global", "Asia/Jakarta", now_local=before_midnight, **astro
    )

    # 2. Cek setelah tengah malam (00:05)
    after_midnight = dt(2026, 1, 19, 0, "Asia/Jakarta")  # menggunakan jam 0
    h_after = get_hijri_date(
        -6.2, 106.8, "global", "Asia/Jakarta", now_local=after_midnight, **astro
    )

    # Keduanya harus menghasilkan tanggal Hijriyah yang SAMA
    # karena keduanya berada di antara Maghrib tgl 18 dan Maghrib tgl 19.
    assert h_before["day"] == h_after["day"]
    assert h_before["month"] == h_after["month"]
