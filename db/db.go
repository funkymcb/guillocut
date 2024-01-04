package db

import (
	"log/slog"
)

func Connect(slog *slog.Logger) error {
	slog.Info("connecting to mongo database...")

	// TODO mongo.Connect()

	return nil
}
