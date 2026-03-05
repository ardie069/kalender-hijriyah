import os
import logging
from logging.handlers import RotatingFileHandler
from pathlib import Path

IS_VERCEL = os.getenv("VERCEL") == "1"


def setup_logging(level: str = "INFO") -> None:
    log_level = getattr(logging, level.upper(), logging.INFO)

    formatter = logging.Formatter(
        fmt="[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    root_logger = logging.getLogger()
    root_logger.setLevel(log_level)

    if root_logger.handlers:
        root_logger.handlers.clear()

    console_handler = logging.StreamHandler()
    console_handler.setLevel(log_level)
    console_handler.setFormatter(formatter)
    root_logger.addHandler(console_handler)

    if not IS_VERCEL:
        try:
            BASE_DIR = Path(__file__).resolve().parent.parent.parent
            LOG_DIR = BASE_DIR / "logs"
            LOG_DIR.mkdir(exist_ok=True)
            LOG_FILE = LOG_DIR / "app.log"

            file_handler = RotatingFileHandler(
                LOG_FILE,
                maxBytes=5 * 1024 * 1024,
                backupCount=3,
                encoding="utf-8",
            )
            file_handler.setLevel(log_level)
            file_handler.setFormatter(formatter)
            root_logger.addHandler(file_handler)
        except Exception as e:
            print(f"Warning: Gagal inisialisasi FileHandler: {e}")

    # Redam log noisy dari library
    logging.getLogger("uvicorn.access").setLevel(logging.WARNING)
    logging.getLogger("skyfield").setLevel(logging.WARNING)
