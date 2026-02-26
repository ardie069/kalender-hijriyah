HIJRI_MONTHS = [
    "Muharram",
    "Shafar",
    "Rabiul awwal",
    "Rabiul akhir",
    "Jumadil awwal",
    "Jumadil akhir",
    "Rajab",
    "Sha'ban",
    "Ramadan",
    "Shawwal",
    "Dhul Qa'dah",
    "Dhul Hijjah",
]

# ====== Location =======

AL_HARAM_LOCATION = {
    "lat": 21.4225,
    "lon": 39.8262,
    "timezone": "Asia/Riyadh",
}

NZ_LOCATION = {
    "lat": -41.2865,
    "lon": 174.7762,
    "timezone": "Pacific/Auckland",
}

REGIONAL_RUKYAT_CONFIG = {
    "indonesia": {
        "criteria": "MABIMS",
        "sites": [
            {"name": "Aceh", "lat": 5.89, "lon": 95.32, "timezone": "Asia/Jakarta"},
            {
                "name": "Jakarta",
                "lat": -6.17,
                "lon": 106.82,
                "timezone": "Asia/Jakarta",
            },
            {
                "name": "Makassar",
                "lat": -5.13,
                "lon": 119.42,
                "timezone": "Asia/Makassar",
            },
            {
                "name": "Jayapura",
                "lat": -2.59,
                "lon": 140.68,
                "timezone": "Asia/Jayapura",
            },
        ],
    },
    "turkey": {
        "criteria": "TURKEY_2016_TOPOCENTRIC",
        "sites": [
            {
                "name": "Ankara",
                "lat": 39.93,
                "lon": 32.85,
                "timezone": "Europe/Istanbul",
            }
        ],
    },
}
