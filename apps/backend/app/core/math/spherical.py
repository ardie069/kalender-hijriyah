from .angles import deg_to_rad, rad_to_deg
import math


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

    cos_sep = min(1.0, max(-1.0, cos_sep))
    return rad_to_deg(math.acos(cos_sep))
