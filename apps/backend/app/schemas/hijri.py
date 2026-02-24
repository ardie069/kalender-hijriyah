import numpy as np
from pydantic import BaseModel, field_validator, ConfigDict
from datetime import datetime
from typing import Literal, Optional, Any

HijriMethod = Literal["umm_al_qura", "local_hisab", "local_rukyat", "ughc"]


def sanitize_numpy(v: Any) -> Any:
    """Konversi rekursif tipe NumPy ke Python standar."""
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
    """Base model otomatis membersihkan tipe data NumPy."""

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
    """Data astronomis lengkap, termasuk koordinat referensi (untuk UGHC)."""

    moon_altitude: float
    elongation: float
    moon_age: Optional[float] = None
    is_visible: bool
    lat: Optional[float] = None
    lon: Optional[float] = None
    location_name: Optional[str] = None


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


class HijriEndMonthResponse(NumPyBaseModel):
    method: HijriMethod
    location: LocationSchema
    generated_at: datetime
    today: HijriDateSchema
    estimated_next_month_1: Optional[HijriDateSchema] = None
    estimated_gregorian: Optional[datetime] = None
    visibility: Optional[HijriAstronomicalDataSchema] = None
    message: Optional[str] = None
