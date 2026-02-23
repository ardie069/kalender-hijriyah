import os
from dotenv import load_dotenv

from fastapi import FastAPI, Response
from fastapi.middleware.cors import CORSMiddleware
from slowapi.errors import RateLimitExceeded
from slowapi.middleware import SlowAPIMiddleware
from slowapi import _rate_limit_exceeded_handler
from fastapi.responses import FileResponse

from app.deps.rate_limit import limiter
from app.api.hijri import router as hijri_router

load_dotenv()

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
ICON_PATH = os.path.join(BASE_DIR, "static", "icon.png")


def get_cors_origins() -> list[str]:
    raw = os.getenv("CORS_ORIGINS", "")
    origins = [o.strip() for o in raw.split(",") if o.strip()]

    return origins if origins else ["http://localhost:3000"]


CORS_ORIGINS = get_cors_origins()

app = FastAPI(
    title="Kalender Hijriyah API",
    version="3.2.0",
)

app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)
app.add_middleware(SlowAPIMiddleware)

app.add_middleware(
    CORSMiddleware,
    allow_origins=CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(hijri_router, prefix="/api")


@app.get("/")
def health():
    return {"status": "OK"}


@app.get("/favicon.ico", include_in_schema=False)
async def favicon():
    if os.path.exists(ICON_PATH):
        return FileResponse(ICON_PATH)
    return Response(status_code=204)
