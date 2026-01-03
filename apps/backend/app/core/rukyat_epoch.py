from dataclasses import dataclass
from typing import Dict, Tuple


@dataclass(frozen=True)
class RukyatMonthEpoch:
    year: int
    month: int
    jd_start: float


RUKYAT_EPOCHS: Dict[Tuple[int, int], RukyatMonthEpoch] = {
    (1447, 7): RukyatMonthEpoch(
        year=1447,
        month=7,
        jd_start=2461031.0,
    ),
}

def get_active_rukyat_epoch(current_jd):
    if not RUKYAT_EPOCHS:
        raise ValueError("No rukyat epoch defined")

    epochs = sorted(
        RUKYAT_EPOCHS.values(),
        key=lambda e: e.jd_start,
    )

    for epoch in reversed(epochs):
        if epoch.jd_start <= current_jd:
            return epoch

    raise ValueError("No active rukyat epoch found for current JD")


def get_rukyat_day(current_jd: float) -> int:
    epoch = get_active_rukyat_epoch(current_jd)

    return int(current_jd - epoch.jd_start) + 1


def get_rukyat_date(current_jd: float) -> dict:
    epoch = get_active_rukyat_epoch(current_jd)
    day = int(current_jd - epoch.jd_start) + 1

    return {
        "year": epoch.year,
        "month": epoch.month,
        "day": day,
    }
