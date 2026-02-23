import os
import json
import hashlib
from pathlib import Path
from datetime import datetime

CACHE_DIR = Path("cache")
CACHE_DIR.mkdir(exist_ok=True)


def cleanup_old_cache(days=60):
    now = datetime.utcnow().timestamp()

    for file in CACHE_DIR.glob("*.json"):
        if now - file.stat().st_mtime > days * 86400:
            file.unlink()


def _make_key(*args):
    raw = "|".join(str(a) for a in args)
    return hashlib.sha256(raw.encode()).hexdigest()


def get_cache(key):
    file_path = CACHE_DIR / f"{key}.json"

    if not file_path.exists():
        return None

    try:
        with open(file_path, "r") as f:
            data = json.load(f)
        return data
    except Exception:
        return None


def set_cache(key, value):
    file_path = CACHE_DIR / f"{key}.json"

    with open(file_path, "w") as f:
        json.dump(value, f, default=str)
