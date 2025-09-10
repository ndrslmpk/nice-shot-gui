package main

import (
	"math"
	"math/rand"
	"net/http"
	"nice-shot/backend/models"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	// shots is our in-memory dataset
	shots []models.Shot
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

	e.Logger.Fatal(e.Start(":8080"))
}

// generateMockShots creates N records with dates from Aug 1 to today
func generateMockShots(n int) []models.Shot {
	var mocksArr []models.Shot
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
	daysRange := int(now.Sub(aug1).Hours()/24) + 1

	for i := 0; i < n; i++ {
		// Prepare MockData
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

		shot := models.Shot{
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
		mocksArr = append(mocksArr, shot)
	}

	// Sort by time asc to ease later slicing; handlers can reverse when needed
	sort.Slice(mocksArr, func(i, j int) bool { return mocksArr[i].BrewTime.Before(mocksArr[j].BrewTime) })
	return mocksArr
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
	var recent []models.Shot
	for i := n - 1; i >= start; i-- {
		recent = append(recent, shots[i])
		if len(recent) >= limit {
			break
		}
	}
	return c.JSON(http.StatusOK, recent)
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
