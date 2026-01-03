from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.hijri import router as hijri_router

app = FastAPI(
    title="Kalender Hijriyah API",
    version="2.0.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://127.0.0.1:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(hijri_router, prefix="/api")


@app.get("/")
def health():
    return {"status": "OK"}
