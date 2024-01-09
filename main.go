package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/funkymcb/guillocut/config"
	"github.com/funkymcb/guillocut/handlers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Get()
	if err != nil {
		slog.Error("error processing config", "message", err.Error())
		os.Exit(1)
	}

	/* TODO we only need to connect to db on login
	db.Logger = logger
	c := context.Background()
	if err := db.Connect(); err != nil {
		slog.Error("could not connect to mongo database", "message", err)
		os.Exit(1)
	} */

	handlers.Logger = logger
	http.Handle("/", handlers.NewHomeHandler())
	http.Handle("/login", handlers.NewLoginHandler())

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	logger.Info("Listening...",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
	)
	logger.Error(http.ListenAndServe(addr, nil).Error())
}
