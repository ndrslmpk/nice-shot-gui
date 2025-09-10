package httpadapter

import (
	"net/http"
	"strconv"
	"time"

	"nice-shot/backend/internal/application"
	"nice-shot/backend/internal/domain"

	"github.com/labstack/echo/v4"
)

type ShotHandlers struct {
	service *application.ShotService
}

func NewShotHandlers(service *application.ShotService) *ShotHandlers {
	return &ShotHandlers{service: service}
}

// Register sets up routes on the given echo.Group or echo.Echo
func (h *ShotHandlers) Register(g *echo.Group) {
	g.GET("/health", h.Health)
	g.GET("/shots", h.List)
}

func (h *ShotHandlers) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *ShotHandlers) List(c echo.Context) error {
	limit := 100
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > 1000 {
				n = 1000
			}
			limit = n
		}
	}
	shots, err := h.service.List(c.Request().Context(), limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// map domain -> DTO to preserve original JSON field names
	dtos := make([]shotDTO, 0, len(shots))
	for _, s := range shots {
		dtos = append(dtos, toDTO(s))
	}
	return c.JSON(http.StatusOK, dtos)
}

// shotDTO mirrors the original API JSON schema (snake_case)
type shotDTO struct {
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

func toDTO(s domain.Shot) shotDTO {
	return shotDTO{
		ShotID:          s.ShotID,
		BrewTime:        s.BrewTime,
		MachineID:       s.MachineID,
		UserID:          s.UserID,
		SoftwareBundle:  s.SoftwareBundle,
		CoffeeType:      s.CoffeeType,
		RecipeID:        s.RecipeID,
		GrindSizeActual: s.GrindSizeActual,
		GrindSizeTarget: s.GrindSizeTarget,
		DoseGrams:       s.DoseGrams,
		DoseTargetGrams: s.DoseTargetGrams,
		BrewTimeSeconds: s.BrewTimeSeconds,
		PeakPressureBar: s.PeakPressureBar,
		LastStatus:      s.LastStatus,
	}
}
