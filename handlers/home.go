package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/funkymcb/guillocut/components"
)

var Logger *slog.Logger

type HomeHandler struct {
	Log *slog.Logger
}

func NewHomeHandler() HomeHandler {
	return HomeHandler{
		Log: Logger,
	}
}

func (hh HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := newResponseWriter(w)
	ts := time.Now().UTC()

	if r.URL.Path != "/" {
		errorHandler(hh.Log, ts, rw, r, http.StatusNotFound)
		return
	}

	templ.Handler(components.Home()).ServeHTTP(rw, r)

	log(hh.Log, "info", ts, r, rw)
}
