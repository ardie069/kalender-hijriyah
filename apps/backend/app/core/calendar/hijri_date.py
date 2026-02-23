def increment_hijri_day(date):
    if date["day"] < 30:
        return {**date, "day": date["day"] + 1}
    return start_new_month(date)


def decrement_hijri_day(date):
    d, m, y = date["day"], date["month"], date["year"]
    if d > 1:
        return {**date, "day": d - 1}
    pm, py = (m - 1, y) if m > 1 else (12, y - 1)
    return {"year": py, "month": pm, "day": 30}


def start_new_month(date):
    m, y = date["month"], date["year"]
    nm, ny = (m + 1, y) if m < 12 else (1, y + 1)
    return {"year": ny, "month": nm, "day": 1}
