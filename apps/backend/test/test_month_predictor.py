from app.core.month_predictor import predict_end_of_month


def test_month_predictor_only_on_day_29(astro):
    result = predict_end_of_month(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        **astro,
    )

    assert "today" in result
    assert "estimated_end_of_month" in result


def test_month_predictor_structure(astro):
    result = predict_end_of_month(
        lat=21.4225,
        lon=39.8262,
        method="global",
        timezone="Asia/Riyadh",
        **astro,
    )

    assert isinstance(result, dict)
