package main

import (
	"fmt"
	"net/http"

	"database-concurrency/config"
	"database-concurrency/internal/handler"
	"database-concurrency/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	cfg, err := config.NewConfig()
	if err != nil {
		log.WithError(err).Error("failed to load config")
		return
	}

	repo, err := repository.Init(cfg)
	if err != nil {
		log.WithError(err).Error("failed to init repo")
		return
	}

	// Handler
	hdl := handler.New(repo, log)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/healthz", healthz)

	apiV1 := e.Group("/api/v1")

	apiV1.Add(http.MethodGet, "/transactions/:hash", hdl.Transaction)

	if err := e.Start(fmt.Sprintf(":%s", cfg.SERVER_PORT)); err != nil {
		log.WithError(err).Error("failed to start server")
	}
}

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}
