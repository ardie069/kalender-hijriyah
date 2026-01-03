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


# Spherical Astronomy


def angular_separation(
    ra1_deg: float,
    dec1_deg: float,
    ra2_deg: float,
    dec2_deg: float,
) -> float:
    """
    Compute angular separation between two sky positions.
    Input & output in degrees.
    """

    ra1 = deg_to_rad(ra1_deg)
    dec1 = deg_to_rad(dec1_deg)
    ra2 = deg_to_rad(ra2_deg)
    dec2 = deg_to_rad(dec2_deg)

    delta_ra = ra2 - ra1
    cos_sep = math.sin(dec1) * math.sin(dec2) + math.cos(dec1) * math.cos(
        dec2
    ) * math.cos(delta_ra)

    # Clamp numerical noise
    cos_sep = min(1.0, max(-1.0, cos_sep))
    return rad_to_deg(math.acos(cos_sep))
