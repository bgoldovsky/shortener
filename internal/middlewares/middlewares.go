package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// Recovering Восстанавливает работу сервера после паники
func Recovering(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if a := recover(); a != nil {
				logrus.WithFields(logrus.Fields{
					"method":    r.Method,
					"path":      r.URL.Path,
					"stack":     string(debug.Stack()),
					"recovered": a,
				}).Error("panic recovered")

				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal service error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Logging Логирует детали запроса
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithField("method", r.Method).
			WithField("path", r.URL.Path).
			Infoln("request")
		next.ServeHTTP(w, r)
	})
}
