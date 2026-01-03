from datetime import datetime
import pytz  # type: ignore
from core.hijri_explain import explain_hijri_decision


def test_explain_global_structure(astro):
    now = datetime(2026, 1, 3, 1, 0, tzinfo=pytz.utc)

    result = explain_hijri_decision(
        lat=21.4225,
        lon=39.8262,
        method="global",
        timezone="Asia/Riyadh",
        now_local=now,
        **astro,
    )

    assert result["method"] == "global"
    assert "baseline_date" in result
    assert "decision" in result
    assert "final_hijri_date" in result


def test_explain_hisab_decision_keys(astro):
    now = datetime(2026, 1, 18, 16, 0, tzinfo=pytz.utc)

    result = explain_hijri_decision(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=now,
        **astro,
    )

    assert result["method"] == "hisab"
    assert "decision" in result

    if result["decision"] != "no_evaluation_needed":
        assert "conjunction" in result


def test_explain_rukyat_visibility_block(astro):
    now = datetime(2026, 1, 19, 15, 0, tzinfo=pytz.utc)

    result = explain_hijri_decision(
        lat=-6.2,
        lon=106.8,
        method="rukyat",
        timezone="Asia/Jakarta",
        now_local=now,
        **astro,
    )

    assert result["method"] == "rukyat"
    assert "decision" in result

    if result["decision"] != "no_evaluation_needed":
        assert "visibility" in result
        assert "is_visible" in result["visibility"]


def test_explain_no_evaluation_needed(astro):
    now = datetime(2026, 1, 4, 10, 0, tzinfo=pytz.utc)

    result = explain_hijri_decision(
        lat=-6.2,
        lon=106.8,
        method="hisab",
        timezone="Asia/Jakarta",
        now_local=now,
        **astro,
    )

    assert result["decision"] in (
        "no_evaluation_needed",
        "global_standard",
    )
