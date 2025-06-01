from flask import Flask, request, jsonify, send_from_directory
from flask_cors import CORS
from dotenv import load_dotenv
import os
from hijri_calculator import get_hijri_date, predict_end_of_month
from pathlib import Path

# Load environment variables
load_dotenv()

app = Flask(__name__, static_folder='../frontend/dist', static_url_path='/')
CORS(app, resources={r"/api/*": {"origins": ["http://127.0.0.1:5173", "http://127.0.0.1:5174"]}})

# Constants from environment or defaults
DEFAULT_LAT = float(os.getenv("DEFAULT_LAT", "0"))
DEFAULT_LON = float(os.getenv("DEFAULT_LON", "0"))
HIJRI_METHOD = os.getenv("HIJRI_METHOD", "global")
TIMEZONE = os.getenv("TIMEZONE", "UTC")

# Middleware-like behavior (cache control)
@app.after_request
def add_header(response):
    response.headers["Cache-Control"] = "no-store, no-cache, must-revalidate, proxy-revalidate"
    response.headers["Pragma"] = "no-cache"
    response.headers["Expires"] = "0"
    return response

# Endpoint: Tanggal Hijriyah hari ini
@app.route("/api/hijri-date", methods=["GET"])
def hijri_date():
    try:
        lat = float(request.args.get("lat", DEFAULT_LAT))
        lon = float(request.args.get("lon", DEFAULT_LON))
        method = request.args.get("method", HIJRI_METHOD)
        timezone = request.args.get("timezone", TIMEZONE)

        hijri_date = get_hijri_date(lat, lon, method, timezone)
        return jsonify({"hijriDate": hijri_date})
    except Exception as e:
        return jsonify({"error": "Terjadi kesalahan pada server"}), 500

# Endpoint: Prediksi akhir bulan Hijriyah
@app.route("/api/hijri-end-month", methods=["GET"])
def hijri_end_month():
    try:
        lat = float(request.args.get("lat", DEFAULT_LAT))
        lon = float(request.args.get("lon", DEFAULT_LON))
        method = request.args.get("method", HIJRI_METHOD)
        timezone = request.args.get("timezone", TIMEZONE)

        prediction = predict_end_of_month(lat, lon, method, timezone)
        return jsonify(prediction)
    except Exception as e:
        print("❌ Gagal prediksi akhir bulan:", e)
        return jsonify({"error": "Terjadi kesalahan saat memprediksi akhir bulan"}), 500

# Sajikan frontend static build
@app.route("/", defaults={"path": ""})
@app.route("/<path:path>")
def serve_frontend(path):
    dist_path = Path(app.static_folder)
    file_path = dist_path / path
    if file_path.exists():
        return send_from_directory(dist_path, path)
    else:
        return send_from_directory(dist_path, "index.html")

# Jalankan server
if __name__ == "__main__":
    port = int(os.getenv("PORT", 5000))
    app.run(debug=True, port=port)
