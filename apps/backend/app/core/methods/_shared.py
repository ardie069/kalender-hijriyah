"""
Shared logic for local Hijri methods (local_hisab & local_rukyat).

Menghindari duplikasi kode antara kedua method tersebut.
"""

from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    check_historical_lag,
)
from ..calendar.hijri_date import (
    increment_hijri_day,
    decrement_hijri_day,
)
from .base import HijriResult


def prepare_hijri_baseline(context, criteria: str):
    """
    Shared pipeline: Calculate baseline, apply lag correction, check sunset.

    Returns:
        (baseline_date, after_sunset, sunset_local)
    """
    baseline, noon_jd = calculate_baseline_hijri(
        context.now_local,
        context.timezone,
    )

    is_lagging = check_historical_lag(
        baseline,
        noon_jd,
        context.lat,
        context.lon,
        context.timezone,
        criteria=criteria,
    )

    if is_lagging:
        baseline = decrement_hijri_day(baseline)

    sunset_local = calculate_sunset(
        context.now_local.date(),
        context.lat,
        context.lon,
        context.timezone,
    )

    after_sunset = context.now_local >= sunset_local if sunset_local else False

    return baseline, after_sunset, sunset_local


def handle_before_sunset(baseline, method_type: str) -> HijriResult:
    """Return early result if before sunset."""
    return HijriResult(
        hijri_date=baseline,
        metadata={
            "type": method_type,
            "decision": "before_maghrib",
        },
    )


def handle_normal_day(baseline, method_type: str) -> HijriResult:
    """Handle non-29th day increment."""
    return HijriResult(
        hijri_date=increment_hijri_day(baseline),
        metadata={
            "type": method_type,
            "decision": "normal_increment",
        },
    )
