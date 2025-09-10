package domain

import "time"

// Shot models one espresso shot record in the core domain.
type Shot struct {
	ShotID          string
	BrewTime        time.Time
	MachineID       string
	UserID          string
	SoftwareBundle  string
	CoffeeType      string
	RecipeID        string
	GrindSizeActual int
	GrindSizeTarget int
	DoseGrams       float64
	DoseTargetGrams float64
	BrewTimeSeconds float64
	PeakPressureBar float64
	LastStatus      string
}
