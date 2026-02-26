from pydantic import BaseModel
from typing import Optional

from . import NumPyBaseModel


class HijriDateBrief(BaseModel):
    year: int
    month: int
    day: int
    month_name: str


class RukyatEvaluateResponse(NumPyBaseModel):
    """Response unified dari endpoint /rukyat/evaluate."""

    is_rukyat_day: bool
    hijri_date: Optional[HijriDateBrief] = None
    sunset_time: Optional[str] = None
    message: Optional[str] = None

    # Data untuk Simulator V4 — sesuai output RukyatService
    altitude_at_sunset: Optional[float] = None
    azimuth_at_sunset: Optional[float] = None
    azimuth_diff: Optional[float] = None
    elongation_at_sunset: Optional[float] = None
    moon_age_hours: Optional[float] = None
    is_visible: Optional[bool] = None
    criteria_used: Optional[str] = None
    site_name: Optional[str] = None

    # National-mode only
    is_visible_national: Optional[bool] = None
    best_site: Optional[str] = None

    error: Optional[str] = None
