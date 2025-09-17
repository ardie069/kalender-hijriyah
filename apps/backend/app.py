import os
from flask import Flask, jsonify, request
from dotenv import load_dotenv
from flask import send_from_directory
from hijri_calculator import get_hijri_date
from month_predictor import predict_end_of_month
from flask_cors import CORS

# Load .env file
load_dotenv()

# Inisialisasi aplikasi Flask
app = Flask(__name__)

# Mendapatkan nilai default dari variabel .env
DEFAULT_LAT = float(os.getenv("DEFAULT_LAT", 0))
DEFAULT_LON = float(os.getenv("DEFAULT_LON", 0))
HIJRI_METHOD = os.getenv("HIJRI_METHOD", "global")
TIMEZONE = os.getenv("TIMEZONE", "UTC")

# CORS (Cross-Origin Resource Sharing) untuk frontend
CORS(
    app,
    resources={
        r"/api/*": {"origins": ["http://127.0.0.1:5173", "http://127.0.0.1:5174"]}
    },
)


# Middleware untuk menghindari cache
@app.after_request
def add_no_cache_headers(response):
    response.cache_control.no_store = True
    response.cache_control.must_revalidate = True
    response.cache_control.no_cache = True
    response.headers["Pragma"] = "no-cache"
    response.headers["Expires"] = "0"
    return response


@app.route("/")
def index():
    return send_from_directory("static", "index.html")


@app.route("/<path:path>")
def static_proxy(path):
    return send_from_directory("static", path)


# Endpoint: Tanggal Hijriyah hari ini
@app.route("/api/hijri-date", methods=["GET"])
def hijri_date():
    lat = float(request.args.get("lat", DEFAULT_LAT))
    lon = float(request.args.get("lon", DEFAULT_LON))
    method = request.args.get("method", HIJRI_METHOD)
    timezone = request.args.get("timezone", TIMEZONE)

    if lat == 0 or lon == 0:
        return jsonify({"error": "Lokasi tidak valid"}), 400

    try:
        hijri_date = get_hijri_date(lat, lon, method, timezone)
        return jsonify({"hijriDate": hijri_date})
    except Exception as error:
        print("❌ Error saat mengambil tanggal Hijriyah:", error)
        import traceback

        traceback.print_exc()
        return jsonify({"error": "Terjadi kesalahan pada server"}), 500


# Endpoint: Prediksi akhir bulan Hijriyah
@app.route("/api/hijri-end-month", methods=["GET"])
def hijri_end_month():
    lat = float(request.args.get("lat", DEFAULT_LAT))
    lon = float(request.args.get("lon", DEFAULT_LON))
    method = request.args.get("method", HIJRI_METHOD)
    timezone = request.args.get("timezone", TIMEZONE)

    try:
        prediction = predict_end_of_month(lat, lon, method, timezone)
        return jsonify(
            {
                "today": prediction["today"],
                "estimatedEndOfMonth": prediction["estimated_end_of_month"],
                "message": prediction["message"],
            }
        )
    except Exception as error:
        print(f"❌ Gagal prediksi akhir bulan: {error}")
        return jsonify({"error": "Terjadi kesalahan saat memprediksi akhir bulan"}), 500


# Menjalankan aplikasi
if __name__ == "__main__":
    port = int(os.getenv("PORT", 5000))
    app.run(host="127.0.0.1", port=port, debug=True)
