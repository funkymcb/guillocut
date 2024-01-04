package handlers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/funkymcb/guillocut/components"
)

type GuillocutAPI interface {
	Get(ctx context.Context, sessionID string) (err error)
}

type DefaultHandler struct {
	Log *slog.Logger
}

func New(log *slog.Logger) *DefaultHandler {
	return &DefaultHandler{
		Log: log,
	}
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Get(w, r)
}

func (h *DefaultHandler) Get(w http.ResponseWriter, r *http.Request) {
	// call any backend service from here (../services/...)
	h.View(w, r)
	h.log("info", "incoming request", r)
}

func (h *DefaultHandler) View(w http.ResponseWriter, r *http.Request) {
	if err := components.Login().Render(r.Context(), w); err != nil {
		h.log("error", err.Error(), r)
	}
}

func (h *DefaultHandler) log(logType, msg string, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.Log.Error(err.Error())
	}

	switch strings.ToLower(logType) {
	case "info":
		h.Log.Info(
			msg,
			"method", r.Method,
			"uri", fmt.Sprintf("%s%s", r.Host, r.URL.String()),
			"body", b,
			"user-agent", r.UserAgent(),
		)
	case "error":
		h.Log.Error(
			msg,
			"method", r.Method,
			"uri", fmt.Sprintf("%s%s", r.Host, r.URL.String()),
			"body", b,
			"user-agent", r.UserAgent(),
		)
	}
}
