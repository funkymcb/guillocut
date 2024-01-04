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

const (
	ListeningHostENV     = "LISTENING_HOST"
	ListeningPortENV     = "LISTENING_PORT"
	ListeningHostDefault = "0.0.0.0"
	ListeningPortDefault = "8080"
)

func main() {
	slog := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := config.Read(slog); err != nil {
		slog.Error("error processing config", "message", err.Error())
		os.Exit(1)
	}

	// TODO fix empty config fmt.Println(config.Cfg)

	if err := db.Connect(slog); err != nil {
		slog.Error("could not connect to mongo database", "message", err)
		os.Exit(1)
	}

	a := getServerAdress()
	h := handlers.New(slog)

	server := &http.Server{
		Addr:         a,
		Handler:      h,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	// TODO slog this! fmt.Printf("Listening on %v\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func getServerAdress() string {
	host := os.Getenv(ListeningHostENV)
	if host == "" {
		host = ListeningHostDefault
	}
	port := os.Getenv(ListeningPortENV)
	if port == "" {
		port = ListeningPortDefault
	}

	return fmt.Sprintf("%s:%s", host, port)
}
