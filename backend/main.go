package main

import (
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

var (
	// shots is our in-memory dataset
	shots []Shot
)

func main() {
	// Generate deterministic yet varied mock data
	shots = generateMockShots(250)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173", "*"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions},
	}))

	api := e.Group("/api")
	api.GET("/health", func(c echo.Context) error { return c.JSON(http.StatusOK, map[string]string{"status": "ok"}) })
	api.GET("/shots", handleGetShots)
	api.GET("/stats/overview", handleGetOverviewStats)
	api.GET("/stats/daily", handleGetDailyStats)

	e.Logger.Fatal(e.Start(":8080"))
}

// generateMockShots creates N records with dates from Aug 1 to today
func generateMockShots(n int) []Shot {
	const seed int64 = 20240801
	prng := rand.New(rand.NewSource(seed))

	// Date range: Aug 1 of current or previous year if today < Aug
	now := time.Now()
	year := now.Year()
	aug1 := time.Date(year, time.August, 1, 8, 0, 0, 0, time.Local)
	if now.Before(aug1) {
		aug1 = time.Date(year-1, time.August, 1, 8, 0, 0, 0, time.Local)
	}

	machines := []string{"nxlc-100", "nxlc-200", "nxlc-300"}
	users := []string{"barista.alex", "barista.sam", "barista.taylor", "barista.jordan"}
	bundles := []string{"stable-1.4.2", "stable-1.5.0", "edge-1.6.0"}
	coffeeTypes := []string{"espresso", "ristretto", "lungo"}
	statuses := []string{"ok", "ok", "ok", "ok", "warning", "error"}

	var out []Shot
	daysRange := int(now.Sub(aug1).Hours()/24) + 1
	for i := 0; i < n; i++ {
		d := prng.Intn(daysRange)
		day := aug1.Add(time.Duration(d) * 24 * time.Hour)
		// Random time of day between 6:00 and 18:00
		hour := 6 + prng.Intn(12)
		min := prng.Intn(60)
		brewAt := time.Date(day.Year(), day.Month(), day.Day(), hour, min, 0, 0, time.Local)

		grindTarget := 20 + prng.Intn(40) // 20..59
		grindActual := clampInt(int(float64(grindTarget)+prng.NormFloat64()*3), 10, 80)
		doseTarget := 18.0 + prng.Float64()*4.0 // 18..22
		dose := doseTarget + (prng.Float64()-0.5)*2.0
		pressure := clampFloat(7.0+prng.NormFloat64()*1.2, 6.0, 11.0)
		brewSeconds := clampFloat(27.0+prng.NormFloat64()*4.0, 18.0, 40.0)

		shot := Shot{
			ShotID:          uuid.NewString(),
			BrewTime:        brewAt,
			MachineID:       pick(machines, prng),
			UserID:          pick(users, prng),
			SoftwareBundle:  pick(bundles, prng),
			CoffeeType:      pick(coffeeTypes, prng),
			RecipeID:        "rx-" + strconv.Itoa(100+prng.Intn(50)),
			GrindSizeActual: grindActual,
			GrindSizeTarget: grindTarget,
			DoseGrams:       round2(dose),
			DoseTargetGrams: round2(doseTarget),
			BrewTimeSeconds: round2(brewSeconds),
			PeakPressureBar: round2(pressure),
			LastStatus:      pick(statuses, prng),
		}
		out = append(out, shot)
	}

	// Sort by time asc to ease later slicing; handlers can reverse when needed
	sort.Slice(out, func(i, j int) bool { return out[i].BrewTime.Before(out[j].BrewTime) })
	return out
}

func handleGetShots(c echo.Context) error {
	// default limit 100, return most recent first
	limit := 100
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > 1000 {
				n = 1000
			}
			limit = n
		}
	}
	n := len(shots)
	start := n - limit
	if start < 0 {
		start = 0
	}
	// reverse order (newest first)
	var recent []Shot
	for i := n - 1; i >= start; i-- {
		recent = append(recent, shots[i])
		if len(recent) >= limit {
			break
		}
	}
	return c.JSON(http.StatusOK, recent)
}

type OverviewStats struct {
	TotalShots         int     `json:"total_shots"`
	AvgBrewTimeSeconds float64 `json:"avg_brew_time_seconds"`
	MinBrewTimeSeconds float64 `json:"min_brew_time_seconds"`
	MaxBrewTimeSeconds float64 `json:"max_brew_time_seconds"`
	AvgPeakPressureBar float64 `json:"avg_peak_pressure_bar"`
	SuccessRatePercent float64 `json:"success_rate_percent"`
}

func handleGetOverviewStats(c echo.Context) error {
	if len(shots) == 0 {
		return c.JSON(http.StatusOK, OverviewStats{})
	}
	var sumBrew, minBrew, maxBrew, sumPressure float64
	minBrew = math.MaxFloat64
	var okCount int
	for _, s := range shots {
		sumBrew += s.BrewTimeSeconds
		if s.BrewTimeSeconds < minBrew {
			minBrew = s.BrewTimeSeconds
		}
		if s.BrewTimeSeconds > maxBrew {
			maxBrew = s.BrewTimeSeconds
		}
		sumPressure += s.PeakPressureBar
		if s.LastStatus == "ok" {
			okCount++
		}
	}
	n := float64(len(shots))
	stats := OverviewStats{
		TotalShots:         len(shots),
		AvgBrewTimeSeconds: round2(sumBrew / n),
		MinBrewTimeSeconds: round2(minBrew),
		MaxBrewTimeSeconds: round2(maxBrew),
		AvgPeakPressureBar: round2(sumPressure / n),
		SuccessRatePercent: round2(float64(okCount) / n * 100.0),
	}
	return c.JSON(http.StatusOK, stats)
}

type DailyStat struct {
	Date               string  `json:"date"`
	Count              int     `json:"count"`
	AvgBrewTimeSeconds float64 `json:"avg_brew_time_seconds"`
	AvgPeakPressureBar float64 `json:"avg_peak_pressure_bar"`
}

func handleGetDailyStats(c echo.Context) error {
	// Aggregate per YYYY-MM-DD
	type agg struct {
		sumBrew, sumPressure float64
		count                int
	}
	m := map[string]*agg{}
	for _, s := range shots {
		key := s.BrewTime.Format("2006-01-02")
		if _, ok := m[key]; !ok {
			m[key] = &agg{}
		}
		a := m[key]
		a.sumBrew += s.BrewTimeSeconds
		a.sumPressure += s.PeakPressureBar
		a.count++
	}
	var days []string
	for k := range m {
		days = append(days, k)
	}
	sort.Strings(days)
	var res []DailyStat
	for _, d := range days {
		a := m[d]
		res = append(res, DailyStat{
			Date:               d,
			Count:              a.count,
			AvgBrewTimeSeconds: round2(a.sumBrew / float64(a.count)),
			AvgPeakPressureBar: round2(a.sumPressure / float64(a.count)),
		})
	}
	return c.JSON(http.StatusOK, res)
}

func pick[T any](list []T, r *rand.Rand) T { return list[r.Intn(len(list))] }
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
func round2(v float64) float64 { return math.Round(v*100) / 100 }
