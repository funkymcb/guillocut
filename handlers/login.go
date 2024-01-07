package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/funkymcb/guillocut/components"
)

type LoginHandler struct {
	Log *slog.Logger
}

func NewLoginHandler() LoginHandler {
	return LoginHandler{
		Log: Logger,
	}
}

func (lh LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := newResponseWriter(w)
	ts := time.Now().UTC()

	if r.URL.Path != "/login" {
		errorHandler(lh.Log, ts, rw, r, http.StatusNotFound)
		return
	}

	templ.Handler(components.Login()).ServeHTTP(rw, r)

	log(lh.Log, "info", ts, r, rw)
}
