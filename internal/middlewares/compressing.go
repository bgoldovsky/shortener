package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var gz *gzip.Writer

func init() {
	var err error
	gz, err = gzip.NewWriterLevel(nil, gzip.BestSpeed)
	if err != nil {
		panic(err)
	}
}

// Compressing Сжимает ответ gzip
func Compressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz.Reset(w)
		defer func(gz *gzip.Writer) {
			_ = gz.Close()
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
