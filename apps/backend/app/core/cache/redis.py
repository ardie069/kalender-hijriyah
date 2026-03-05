import json
import redis
import hashlib
import os
import logging

logger = logging.getLogger(__name__)

REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379")

try:
    r = redis.Redis.from_url(REDIS_URL, decode_responses=True)
except Exception as e:
    logger.error(f"Gagal konek ke Redis: {e}")
    r = None

CACHE_VERSION = "v3"


def _make_key(prefix, *args):
    """Bikin key unik dengan namespace dan versioning."""
    raw = "|".join(str(a) for a in args)
    hashed = hashlib.sha256(raw.encode()).hexdigest()
    return f"hijri:{CACHE_VERSION}:{prefix}:{hashed[:16]}"


def get_cache(key):
    if not r:
        return None
    try:
        value = r.get(key)
        if value:
            return json.loads(value)
    except Exception as e:
        logger.warning(f"Redis GET error: {e}")
    return None


def set_cache(key, value, ttl=86400):
    if not r:
        return
    try:
        r.setex(key, ttl, json.dumps(value, default=str))
    except Exception as e:
        logger.error(f"Redis SET error: {e}")
