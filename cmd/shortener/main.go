package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/generator"
	urlsRepo "github.com/bgoldovsky/shortener/internal/app/repo/urls"
	urlsSrv "github.com/bgoldovsky/shortener/internal/app/services/urls"
	"github.com/bgoldovsky/shortener/internal/config"
	"github.com/bgoldovsky/shortener/internal/handlers"
	"github.com/bgoldovsky/shortener/internal/middlewares"
)

func main() {
	// Config
	cfg, err := config.New()
	panicOnError(err)

	// Repositories
	repo := urlsRepo.Factory(cfg.FileStoragePath)

	// Services
	gen := generator.NewGenerator()
	service := urlsSrv.NewService(repo, gen, cfg.BaseURL)

	// Router
	r := chi.NewRouter()

	// Middlewares
	compress, err := middlewares.NewCompressor()
	panicOnError(err)

	r.Use(middlewares.Logging)
	r.Use(middlewares.Recovering)
	r.Use(middlewares.Decompressing)
	r.Use(compress.Compressing)

	r.Post("/", handlers.New(service).ShortenV1)
	r.Post("/api/shorten", handlers.New(service).ShortenV2)
	r.Get("/{id}", handlers.New(service).Expand)

	// Start service
	address := cfg.ServerAddress
	logrus.WithField("address", address).Info("server starts")
	logrus.Fatal(http.ListenAndServe(address, r))
}

func panicOnError(err error) {
	if err != nil {
		logrus.WithError(err).Error("fatal error")
		panic(err)
	}
}
