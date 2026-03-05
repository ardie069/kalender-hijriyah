import os
import json
import hashlib
from pathlib import Path
import logging

logger = logging.getLogger(__name__)

# Cek apakah di Vercel
IS_VERCEL = os.getenv("VERCEL") == "1"

# Di Vercel, folder yang bisa ditulis cuma /tmp/
# Tapi /tmp/ itu volatile (ilang pas restart). Lebih baik pake Redis.
BASE_DIR = Path(__file__).resolve().parent.parent.parent.parent
CACHE_DIR = Path("/tmp/hijri_cache") if IS_VERCEL else BASE_DIR / "cache"

if not IS_VERCEL:
    CACHE_DIR.mkdir(exist_ok=True)


def _make_key(*args):
    raw = "|".join(str(a) for a in args)
    return hashlib.sha256(raw.encode()).hexdigest()


def get_cache(key, ttl_seconds=86400):
    # Kalau di Vercel, kita skip disk cache dan biarkan dia pakai Redis/Compute
    # Kecuali lu maksa mau pake /tmp/
    file_path = CACHE_DIR / f"{key}.json"

    if not file_path.exists():
        return None

    try:
        # Cek umur file
        import time

        mtime = file_path.stat().st_mtime
        if (time.time() - mtime) > ttl_seconds:
            # JANGAN UNLINK DI VERCEL (kecuali di /tmp)
            if not IS_VERCEL:
                file_path.unlink(missing_ok=True)
            return None

        with open(file_path, "r") as f:
            return json.load(f)
    except Exception as e:
        logger.error(f"Disk Cache Error: {e}")
        return None


def set_cache(key, value):
    # Di Vercel, jangan nulis ke disk app.
    # Kalau mau nulis, harus ke /tmp/
    if IS_VERCEL:
        # Opsi: biarkan kosong agar fallback ke Redis
        return

    file_path = CACHE_DIR / f"{key}.json"
    try:
        with open(file_path, "w") as f:
            json.dump(value, f)
    except Exception as e:
        logger.error(f"Disk Set Cache Error: {e}")
