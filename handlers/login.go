package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/funkymcb/guillocut/components"
)

type LoginService interface {
	Auth(ctx context.Context, r *http.Request)
}

type LoginHandler struct {
	Log          *slog.Logger
	LoginService LoginService
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

	switch r.Method {
	case http.MethodGet:
		templ.Handler(components.Login()).ServeHTTP(rw, r)
		log(lh.Log, "info", ts, r, rw)

	case http.MethodPost:
		// TODO handle POST (user auth) here

	default:
		errorHandler(lh.Log, ts, rw, r, http.StatusMethodNotAllowed)
		return
	}
}
