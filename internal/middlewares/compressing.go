package middlewares

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressor struct {
	gz *gzip.Writer
}

func NewCompressor() (*compressor, error) {
	gz, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("init compressor error: %w", err)
	}

	return &compressor{gz: gz}, nil
}

// Compressing Сжимает ответ gzip
func (c *compressor) Compressing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		c.gz.Reset(w)
		defer func(gz *gzip.Writer) {
			_ = gz.Close()
		}(c.gz)

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: c.gz}, r)
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
