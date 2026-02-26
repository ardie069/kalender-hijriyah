from pydantic import BaseModel
from typing import Optional

from .hijri import NumPyBaseModel


class MoonPositionSchema(NumPyBaseModel):
    """Posisi bulan saat sunset pada hari ke-29 Hijriyah."""

    altitude: float
    elongation: float
    moon_age_hours: float


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
    moon_position: Optional[MoonPositionSchema] = None
    is_visible: Optional[bool] = None
    criteria_used: Optional[str] = None

    # National-mode only
    is_visible_national: Optional[bool] = None
    best_site: Optional[str] = None

    error: Optional[str] = None
