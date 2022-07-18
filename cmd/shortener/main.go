package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/generator"
	"github.com/bgoldovsky/shortener/internal/app/hasher"
	urlsRepository "github.com/bgoldovsky/shortener/internal/app/repo/urls"
	authService "github.com/bgoldovsky/shortener/internal/app/services/auth"
	infraService "github.com/bgoldovsky/shortener/internal/app/services/infra"
	urlsService "github.com/bgoldovsky/shortener/internal/app/services/urls"
	"github.com/bgoldovsky/shortener/internal/config"
	"github.com/bgoldovsky/shortener/internal/handlers"
	"github.com/bgoldovsky/shortener/internal/middlewares"
)

func main() {
	// Config
	cfg, err := config.NewConfig()
	panicOnError(err)

	// Repositories
	urlsRepo, err := urlsRepository.Factory(cfg.FileStoragePath, cfg.DatabaseDSN)
	panicOnError(err)

	defer func(urlsRepo urlsRepository.Repo) {
		_ = urlsRepo.Close()
	}(urlsRepo)

	// Services
	gen := generator.NewGenerator()
	hash := hasher.NewHasher(cfg.Secret)
	urlSrv := urlsService.NewService(urlsRepo, gen, cfg.BaseURL)
	authSrv := authService.NewService(gen, hash)
	infraSrv := infraService.NewService(urlsRepo)

	// Router
	r := chi.NewRouter()

	// Middlewares
	compress, err := middlewares.NewCompressor()
	panicOnError(err)
	auth := middlewares.NewAuthenticator(authSrv)

	r.Use(middlewares.Logging)
	r.Use(middlewares.Recovering)
	r.Use(middlewares.Decompressing)
	r.Use(compress.Compressing)
	r.Use(auth.Auth)

	r.Post("/", handlers.New(urlSrv, auth, infraSrv).ShortenV1)
	r.Post("/api/shorten", handlers.New(urlSrv, auth, infraSrv).ShortenV2)
	r.Post("/api/shorten/batch", handlers.New(urlSrv, auth, infraSrv).ShortenBatch)
	r.Get("/{id}", handlers.New(urlSrv, auth, infraSrv).Expand)
	r.Get("/api/user/urls", handlers.New(urlSrv, auth, infraSrv).GetUrls)
	r.Get("/ping", handlers.New(urlSrv, auth, infraSrv).Ping)

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
