import json
import time
import hashlib
import logging
from pathlib import Path
from datetime import datetime

logger = logging.getLogger(__name__)

CACHE_DIR = Path(__file__).resolve().parent.parent.parent.parent / "cache"
CACHE_DIR.mkdir(exist_ok=True)


def cleanup_old_cache(days=60):
    now = time.time()
    removed = 0
    for file in CACHE_DIR.glob("*.json"):
        if now - file.stat().st_mtime > days * 86400:
            file.unlink()
            removed += 1
    if removed:
        logger.info("Cache cleanup: %d file(s) dihapus (> %d hari)", removed, days)


def _make_key(*args):
    raw = "v2|" + "|".join(str(a) for a in args)
    return hashlib.sha256(raw.encode()).hexdigest()


def get_cache(key, ttl_seconds=86400):
    """
    Ambil cache dari disk. Return None jika tidak ada atau sudah expired.

    Args:
        key: Cache key (hash string)
        ttl_seconds: Time-to-live dalam detik (default 24 jam).
                     Set 0 untuk tidak ada expiry.
    """
    file_path = CACHE_DIR / f"{key}.json"

    if not file_path.exists():
        return None

    # Cek umur cache jika ttl_seconds > 0
    if ttl_seconds > 0:
        age = time.time() - file_path.stat().st_mtime
        if age > ttl_seconds:
            logger.debug(
                "Cache expired: key=%s (umur=%.0fs, ttl=%ds)",
                key[:12],
                age,
                ttl_seconds,
            )
            file_path.unlink(missing_ok=True)
            return None

    try:
        with open(file_path, "r") as f:
            data = json.load(f)
        logger.debug("Cache HIT: key=%s", key[:12])
        return data
    except Exception:
        logger.warning("Cache corrupt, dihapus: key=%s", key[:12])
        file_path.unlink(missing_ok=True)
        return None


def set_cache(key, value):
    file_path = CACHE_DIR / f"{key}.json"
    with open(file_path, "w") as f:
        json.dump(value, f, default=str)
    logger.debug("Cache SET: key=%s", key[:12])
