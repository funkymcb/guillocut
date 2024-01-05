package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/funkymcb/guillocut/config"
	"github.com/funkymcb/guillocut/db"
	"github.com/funkymcb/guillocut/handlers"
)

func main() {
	slog := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Get()
	if err != nil {
		slog.Error("error processing config", "message", err.Error())
		os.Exit(1)
	}

	if err := db.Connect(slog); err != nil {
		slog.Error("could not connect to mongo database", "message", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	hdlr := handlers.New(slog)

	server := &http.Server{
		Addr:         addr,
		Handler:      hdlr,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	slog.Info("Listening...",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
	)
	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
