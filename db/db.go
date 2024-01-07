package db

import "log/slog"

var Logger *slog.Logger

func Connect() error {
	Logger.Info("connecting to mongo database.", "status", "pending")

	// TODO mongo.Connect()

	Logger.Info("connecting to mongo database.", "status", "success")

	return nil
}
