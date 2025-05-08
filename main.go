package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/eu-erwin/protobuf-compiler/cmd"
)

var (
	appName    = "helper-cli"
	appCommit  = "none"
	appVersion = "0.0.0"
	appEnv     = "dev"
)

func main() {
	h := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	).WithAttrs([]slog.Attr{
		slog.String("app_name", appName),
		slog.String("app_commit", appCommit),
		slog.String("app_version", appVersion),
		slog.String("app_env", appEnv),
	})
	logger := slog.New(h)
	ctx := context.Background()

	cmd.Execute(ctx, logger)
}
