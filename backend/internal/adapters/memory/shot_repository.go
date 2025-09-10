package memory

import (
	"context"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"nice-shot/backend/internal/domain"
	"nice-shot/backend/internal/ports"

	"github.com/google/uuid"
)

type ShotRepository struct {
	shots []domain.Shot
}

func NewShotRepositoryWithMocks(n int) ports.ShotRepository {
	repo := &ShotRepository{}
	repo.shots = generateMockShots(n)
	return repo
}

func (r *ShotRepository) Save(ctx context.Context, s domain.Shot) (domain.Shot, error) {
	// Upsert by ShotID
	for i := range r.shots {
		if r.shots[i].ShotID == s.ShotID {
			r.shots[i] = s
			return s, nil
		}
	}
	r.shots = append(r.shots, s)
	return s, nil
}

func (r *ShotRepository) FindByID(ctx context.Context, id string) (domain.Shot, error) {
	for _, s := range r.shots {
		if s.ShotID == id {
			return s, nil
		}
	}
	return domain.Shot{}, context.Canceled // simple sentinel; replace with custom error in real app
}

func (r *ShotRepository) List(ctx context.Context, limit int) ([]domain.Shot, error) {
	if limit <= 0 {
		limit = 100
	}
	n := len(r.shots)
	start := n - limit
	if start < 0 {
		start = 0
	}
	// newest first
	var recent []domain.Shot
	for i := n - 1; i >= start; i-- {
		recent = append(recent, r.shots[i])
		if len(recent) >= limit {
			break
		}
	}
	return recent, nil
}

func (r *ShotRepository) Delete(ctx context.Context, id string) error {
	for i := range r.shots {
		if r.shots[i].ShotID == id {
			r.shots = append(r.shots[:i], r.shots[i+1:]...)
			break
		}
	}
	return nil
}

// generateMockShots mirrors the original main.go generator
func generateMockShots(n int) []domain.Shot {
	var mocksArr []domain.Shot
	const seed int64 = 20240801
	prng := rand.New(rand.NewSource(seed))

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
		d := prng.Intn(daysRange)
		day := aug1.Add(time.Duration(d) * 24 * time.Hour)
		hour := 6 + prng.Intn(12)
		min := prng.Intn(60)
		brewAt := time.Date(day.Year(), day.Month(), day.Day(), hour, min, 0, 0, time.Local)

		grindTarget := 20 + prng.Intn(40)
		grindActual := clampInt(int(float64(grindTarget)+prng.NormFloat64()*3), 10, 80)
		doseTarget := 18.0 + prng.Float64()*4.0
		dose := doseTarget + (prng.Float64()-0.5)*2.0
		pressure := clampFloat(7.0+prng.NormFloat64()*1.2, 6.0, 11.0)
		brewSeconds := clampFloat(27.0+prng.NormFloat64()*4.0, 18.0, 40.0)

		shot := domain.Shot{
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

	sort.Slice(mocksArr, func(i, j int) bool { return mocksArr[i].BrewTime.Before(mocksArr[j].BrewTime) })
	return mocksArr
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
