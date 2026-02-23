import json
import redis
import hashlib

r = redis.Redis(host="localhost", port=6379, db=0)


def _make_key(*args):
    raw = "|".join(str(a) for a in args)
    return "hijri:" + hashlib.sha256(raw.encode()).hexdigest()


def get_cache(key):
    value = r.get(key)
    if value:
        return json.loads(value)
    return None


def set_cache(key, value, ttl=86400):
    r.setex(key, ttl, json.dumps(value, default=str))
