import pytest
from core.ephemeris import ts, eph, sun, moon, earth


@pytest.fixture(scope="session")
def astro():
    return {
        "ts": ts,
        "eph": eph,
        "sun": sun,
        "moon": moon,
        "earth": earth,
    }
