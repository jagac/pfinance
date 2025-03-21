package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"
)

// LoggingConfig is a configuration struct for the logging middleware.
// It holds a reference to a logger that will be used to log HTTP request information.
type LoggingConfig struct {
	Logger *slog.Logger // Logger to be used for logging HTTP request details
}

// Middleware logs incoming HTTP requests, their processing time, status code, and any errors.
// It wraps the provided handler and logs details before and after the request is processed.
func (l *LoggingConfig) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		l.Logger.Info("Incoming request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Time("start_time", start))

		recorder := &StatusRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start).Seconds()

		l.Logger.Info("Request processed",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Float64("duration", duration),
			slog.Int("status", recorder.statusCode))

		if recorder.statusCode >= 400 {
			l.Logger.Error("Request resulted in an error",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Int("status", recorder.statusCode))
		}
	})
}
