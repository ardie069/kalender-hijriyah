from core.conjunction import get_conjunction_time
from core.ephemeris import ts, earth, sun, moon


def test_conjunction_exists():
    jd_start = 2460680.5

    conj = get_conjunction_time(
        jd_start,
        ts=ts,
        earth=earth,
        sun=sun,
        moon=moon,
    )

    assert conj is not None
    assert abs(conj - jd_start) < 1
