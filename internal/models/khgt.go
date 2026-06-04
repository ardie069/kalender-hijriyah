package models

import "time"

type KHGTResult struct {
	Date               time.Time         `json:"date"`
	IsGlobalValid      bool              `json:"is_global_valid"`
	IsAmericaException bool              `json:"is_america_exception"`
	BestVisibility     *HilalPrediction  `json:"best_visibility"`
	ValidLocations     []HilalPrediction `json:"valid_locations"`
}
