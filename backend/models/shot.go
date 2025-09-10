package models

import "time"

// Shot models one espresso shot record
type Shot struct {
	ShotID          string    `json:"shot_id"`
	BrewTime        time.Time `json:"brew_time"`
	MachineID       string    `json:"machine_id"`
	UserID          string    `json:"user_id"`
	SoftwareBundle  string    `json:"software_bundle"`
	CoffeeType      string    `json:"coffee_type"`
	RecipeID        string    `json:"recipe_id"`
	GrindSizeActual int       `json:"grind_size_actual"`
	GrindSizeTarget int       `json:"grind_size_target"`
	DoseGrams       float64   `json:"dose_grams"`
	DoseTargetGrams float64   `json:"dose_target_grams"`
	BrewTimeSeconds float64   `json:"brew_time_seconds"`
	PeakPressureBar float64   `json:"peak_pressure_bar"`
	LastStatus      string    `json:"last_status"`
}
