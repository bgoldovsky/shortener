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
	// Repositories
	repo := urlsRepo.NewRepo()

	// Services
	gen := generator.NewGenerator()
	service := urlsSrv.NewService(repo, gen, config.ShortcutHost())

	// Router
	r := chi.NewRouter()
	r.Use(middlewares.Logging)
	r.Use(middlewares.Recovering)
	r.Post("/", handlers.New(service).Shorten)
	r.Get("/{id}", handlers.New(service).Expand)

	// Start service
	port := config.Port()
	logrus.WithField("port", port).Info("server starts")
	logrus.Fatal(http.ListenAndServe(port, r))
}
