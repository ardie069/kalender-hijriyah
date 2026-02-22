import numpy as np
from pydantic import BaseModel, field_validator, ConfigDict
from datetime import datetime
from typing import Literal, Optional, Any

HijriMethod = Literal["umm_al_qura", "local_hisab", "local_rukyat", "ughc"]


def sanitize_numpy(v: Any) -> Any:
    """Fungsi pembantu untuk konversi rekursif tipe NumPy ke Python standar."""
    if isinstance(v, np.bool_):
        return bool(v)
    if isinstance(v, (np.floating, np.float64, np.float32)):
        return float(v)
    if isinstance(v, (np.integer, np.int64, np.int32)):
        return int(v)
    if isinstance(v, dict):
        return {k: sanitize_numpy(val) for k, val in v.items()}
    if isinstance(v, list):
        return [sanitize_numpy(i) for i in v]
    return v


class NumPyBaseModel(BaseModel):
    """Base model yang otomatis membersihkan tipe data NumPy."""

    model_config = ConfigDict(arbitrary_types_allowed=True)

    @field_validator("*", mode="before")
    @classmethod
    def validate_numpy(cls, v: Any) -> Any:
        return sanitize_numpy(v)


class HijriDateSchema(NumPyBaseModel):
    year: int
    month: int
    day: int


class LocationSchema(NumPyBaseModel):
    lat: float
    lon: float
    timezone: str


class HijriAstronomicalDataSchema(NumPyBaseModel):
    moon_altitude: float
    elongation: float
    moon_age: float
    is_visible: bool


class HijriExplanationSchema(NumPyBaseModel):
    method: HijriMethod
    after_sunset: bool
    criteria_used: str
    reasoning: list[str]
    decision: str
    astronomical_data: Optional[HijriAstronomicalDataSchema] = None


class HijriDateResponse(NumPyBaseModel):
    method: HijriMethod
    location: LocationSchema
    hijri_date: HijriDateSchema
    explanation: Optional[HijriExplanationSchema] = None
    generated_at: datetime


class HijriEndMonthEstimate(NumPyBaseModel):
    """Data astronomis untuk akhir bulan."""

    moon_age: float
    moon_altitude: float
    elongation: float
    is_visible: bool


class HijriEndMonthResponse(NumPyBaseModel):
    method: str
    location: LocationSchema
    generated_at: datetime
    today: HijriDateSchema
    estimated_end_of_month: Optional[HijriDateSchema] = None
    visibility: Optional[HijriEndMonthEstimate] = None
    message: Optional[str] = None
