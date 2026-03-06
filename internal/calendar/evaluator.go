package calendar

import "time"

type VisibilityResult struct {
	Method     string  `json:"method"`
	IsNewMonth bool    `json:"is_new_month"`
	Altitude   float64 `json:"altitude"`
	Elongation float64 `json:"elongation"`
	LagTimeMin float64 `json:"lag_time_minutes"`
}

func (l *Logic) EvaluateLocalHisab(method string, t time.Time, alt, elong float64, sunset, moonset, ijtimaTime time.Time) VisibilityResult {
	res := VisibilityResult{Method: method, Altitude: alt, Elongation: elong}
	isMoonAfterSun, lag := CheckLagTime(sunset, moonset)
	res.LagTimeMin = lag

	switch method {
	case "MABIMS":
		if isMoonAfterSun {
			res.IsNewMonth = IsMABIMS(alt, elong)
		}
	case "TURKEY_2016":
		if isMoonAfterSun {
			res.IsNewMonth = IsTurkey2016(alt, elong)
		}
	case "WUJUDUL_HILAL":
		res.IsNewMonth = IsWujudulHilal(alt, sunset.After(ijtimaTime))
	case "UMM_AL_QURA":
		mLat, mLon := 21.4225, 39.8262
		sMakkah, _ := l.Astro.GetSunset(t, mLat, mLon)
		mMakkah, _ := l.Astro.GetMoonset(t, mLat, mLon)
		res.IsNewMonth = IsUmmAlQura(ijtimaTime, sMakkah, mMakkah)
	}
	return res
}
