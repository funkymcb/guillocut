package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const timeFormat = "02/Jan/2006:15:04:05 -0700"

// implement common log format
// host ident authuser date request status bytes
func log(logger *slog.Logger, logType string, ts time.Time, r *http.Request, rw *responseWriter) {
	switch strings.ToLower(logType) {
	case "info":
		logger.Info("HTTP",
			"ip_address", r.RemoteAddr,
			"user_identifier", "-",
			"user_authentication", "-",
			"timestamp", ts.Format(timeFormat),
			"request_method", r.Method,
			"request_url", r.URL.Path,
			"protocol", r.Proto,
			"status_code", rw.statusCode,
			"response_size", rw.responseSize,
			"user_agent", r.UserAgent(),
		)
	case "error":
		logger.Error("HTTP",
			"ip_address", r.RemoteAddr,
			"user_identifier", "-",
			"user_authentication", "-",
			"timestamp", ts.Format(timeFormat),
			"request_method", r.Method,
			"request_url", r.URL.Path,
			"protocol", r.Proto,
			"status_code", rw.statusCode,
			"response_size", rw.responseSize,
			"user_agent", r.UserAgent(),
		)
	case "warn":
		logger.Warn("HTTP",
			"ip_address", r.RemoteAddr,
			"user_identifier", "-",
			"user_authentication", "-",
			"timestamp", ts.Format(timeFormat),
			"request_method", r.Method,
			"request_url", r.URL.Path,
			"protocol", r.Proto,
			"status_code", rw.statusCode,
			"response_size", rw.responseSize,
			"user_agent", r.UserAgent(),
		)
	}
}
