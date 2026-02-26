"""
Schemas shared base — NumPyBaseModel dan utilitas sanitasi.

Digunakan oleh semua schema Pydantic di project ini.
"""

import numpy as np
from pydantic import BaseModel, field_validator, ConfigDict
from typing import Any


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
