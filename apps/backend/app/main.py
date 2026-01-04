from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from slowapi.errors import RateLimitExceeded
from slowapi.middleware import SlowAPIMiddleware
from slowapi import _rate_limit_exceeded_handler

from app.deps.rate_limit import limiter
from app.api.hijri import router as hijri_router

app = FastAPI(
    title="Kalender Hijriyah API",
    version="3.0.0-beta",
)

app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)  # type: ignore
app.add_middleware(SlowAPIMiddleware)

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://127.0.0.1:5173",
        "http://localhost:5173",
        "http://kalender-hijriyah.netlify.app/.functions",
        "http://kalender-hijriyah-api.vercel.app",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(hijri_router, prefix="/api")


@app.get("/")
def health():
    return {"status": "OK"}
