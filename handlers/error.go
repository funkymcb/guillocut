package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/funkymcb/guillocut/components"
)

func errorHandler(l *slog.Logger, t time.Time, w *responseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	log(l, "error", t, r, w)

	if status == http.StatusNotFound {
		templ.Handler(components.NotFound()).ServeHTTP(w, r)
	}
}
