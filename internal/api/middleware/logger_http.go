package middleware

import (
	"net/http"
	"time"

	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/sirupsen/logrus"
)

// Logger returns HTTP middleware for request logging
func Logger(cfg config.LoggingConfig) func(http.Handler) http.Handler {
	log := logrus.New()

	switch cfg.Level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			latency := time.Since(start)
			fields := logrus.Fields{
				"status":  rw.statusCode,
				"method":  r.Method,
				"path":    r.URL.Path,
				"ip":      r.RemoteAddr,
				"latency": latency,
			}

			entry := log.WithFields(fields)
			if rw.statusCode >= 500 {
				entry.Error("Server error")
			} else if rw.statusCode >= 400 {
				entry.Warn("Client error")
			} else {
				entry.Info("Request completed")
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
