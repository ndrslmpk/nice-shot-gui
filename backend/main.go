package main

import (
	"net/http"
	"time"

	httpadapter "nice-shot/backend/internal/adapters/http"
	"nice-shot/backend/internal/adapters/memory"
	"nice-shot/backend/internal/application"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	repo := memory.NewShotRepositoryWithMocks(250)
	svc := application.NewShotService(repo, time.Now)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173", "*"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions},
	}))

	api := e.Group("/api")
	h := httpadapter.NewShotHandlers(svc)
	h.Register(api)

	e.Logger.Fatal(e.Start(":8080"))
}
