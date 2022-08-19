package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/bgoldovsky/shortener/internal/app/generator"
	"github.com/bgoldovsky/shortener/internal/app/hasher"
	"github.com/bgoldovsky/shortener/internal/app/models"
	urlsRepository "github.com/bgoldovsky/shortener/internal/app/repositories/urls"
	authService "github.com/bgoldovsky/shortener/internal/app/services/auth"
	cleanerService "github.com/bgoldovsky/shortener/internal/app/services/cleaner"
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

	// Channels
	deleteCh := make(chan models.UserCollection)
	doneCh := make(chan struct{})
	defer func() { doneCh <- struct{}{} }()

	// Repositories
	urlsRepo, err := urlsRepository.Factory(cfg.FileStoragePath, cfg.DatabaseDSN)
	panicOnError(err)
	defer func(urlsRepo urlsRepository.Repository) {
		_ = urlsRepo.Close()
	}(urlsRepo)

	// Services
	gen := generator.NewGenerator()
	hash := hasher.NewHasher(cfg.Secret)
	urlsSrv := urlsService.NewService(urlsRepo, gen, cfg.BaseURL)
	authSrv := authService.NewService(gen, hash)
	infraSrv := infraService.NewService(urlsRepo)
	cleanerSrv := cleanerService.NewService(urlsRepo, deleteCh, doneCh)
	cleanerSrv.Run()

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

	r.Post("/", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).ShortenV1)
	r.Post("/api/shorten", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).ShortenV2)
	r.Post("/api/shorten/batch", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).ShortenBatch)
	r.Get("/{id}", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).Expand)
	r.Get("/api/user/urls", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).GetUrls)
	r.Delete("/api/user/urls", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).DeleteUrls)
	r.Get("/ping", handlers.New(urlsSrv, auth, infraSrv, cleanerSrv).Ping)

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
