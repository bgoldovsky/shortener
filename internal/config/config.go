package config

import "os"

// Port Возвращает порт HTTP сервера
func Port() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8080"
	}

	return ":" + p
}

// ShortcutHost Возвращает хост для генерации сокращенного URL
func ShortcutHost() string {
	h := os.Getenv("SHORTCUT_HOST")
	if h == "" {
		h = "http://localhost:8080"
	}

	return h
}
