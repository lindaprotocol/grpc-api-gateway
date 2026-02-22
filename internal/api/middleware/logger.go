package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lindaprotocol/grpc-api-gateway/internal/config"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(cfg config.LoggingConfig) gin.HandlerFunc {
	log := logrus.New()
	
	// Set log level
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

	// Set log format
	switch cfg.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		fields := logrus.Fields{
			"status":     statusCode,
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
		}

		// Add user info if authenticated
		if user, exists := c.Get("auth_user"); exists {
			fields["user_id"] = user
		}

		// Add API key if present
		if apiKey := c.GetHeader("LINDA-PRO-API-KEY"); apiKey != "" {
			fields["api_key"] = apiKey[:8] + "..." // Mask for privacy
		}

		entry := log.WithFields(fields)

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else if statusCode >= 500 {
			entry.Error("Server error")
		} else if statusCode >= 400 {
			entry.Warn("Client error")
		} else {
			entry.Info("Request completed")
		}
	}
}