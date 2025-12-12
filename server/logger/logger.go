package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	if os.Getenv("GIN_MODE") == "release" {
		Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		return
	}

	Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
