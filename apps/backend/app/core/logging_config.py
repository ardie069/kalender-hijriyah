"""
Logging configuration — Menulis log ke folder logs/ dengan rotating file handler.
"""

import os
import logging
from logging.handlers import RotatingFileHandler
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent.parent  # apps/backend/
LOG_DIR = BASE_DIR / "logs"
LOG_DIR.mkdir(exist_ok=True)

LOG_FILE = LOG_DIR / "app.log"


def setup_logging(level: str = "INFO") -> None:
    """
    Setup root logger dengan:
    - RotatingFileHandler → logs/app.log (max 5MB, 3 backup)
    - StreamHandler → console (untuk development)
    """
    log_level = getattr(logging, level.upper(), logging.INFO)

    formatter = logging.Formatter(
        fmt="[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    # File handler — rotating
    file_handler = RotatingFileHandler(
        LOG_FILE,
        maxBytes=5 * 1024 * 1024,  # 5 MB
        backupCount=3,
        encoding="utf-8",
    )
    file_handler.setLevel(log_level)
    file_handler.setFormatter(formatter)

    # Console handler
    console_handler = logging.StreamHandler()
    console_handler.setLevel(log_level)
    console_handler.setFormatter(formatter)

    # Root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(log_level)

    # Hindari duplicate handler jika dipanggil ulang (reload)
    if not root_logger.handlers:
        root_logger.addHandler(file_handler)
        root_logger.addHandler(console_handler)

    # Redam log noisy dari library pihak ketiga
    logging.getLogger("uvicorn.access").setLevel(logging.WARNING)
    logging.getLogger("skyfield").setLevel(logging.WARNING)
