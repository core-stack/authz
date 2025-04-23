package authz

import (
	"log/slog"
	"os"
)

func initLogs() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
}
