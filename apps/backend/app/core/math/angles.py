import math


def deg_to_rad(deg: float) -> float:
    return deg * math.pi / 180.0


def rad_to_deg(rad: float) -> float:
    return rad * 180.0 / math.pi


def normalize_angle(angle: float) -> float:
    """Normalize angle to (0, 360) degrees."""
    return angle % 360.0


def hours_to_deg(hours: float) -> float:
    """Convert hours (RA) to degrees."""
    return hours * 15.0


def deg_to_hours(deg: float) -> float:
    """Convert degrees to hours (RA)."""
    return deg / 15.0