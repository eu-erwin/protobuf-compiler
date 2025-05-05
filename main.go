package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
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

	lowercase := flag.Bool("l", false, "lower case")
	uppercase := flag.Bool("u", false, "upper case")
	capitalize := flag.Bool("t", false, "capitalize / title case")
	camelcase := flag.Bool("c", false, "camel case")
	snakeCase := flag.Bool("s", false, "snake case")
	kebabCase := flag.Bool("k", false, "kebab case")
	pascalCase := flag.Bool("p", false, "pascal case")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		logger.ErrorContext(ctx, "no name provided")
		return
	}

	logger.DebugContext(ctx, "provided args", "args", slog.String("args", strings.Join(args, " ")))
	naming := NewNaming(
		ctx,
		logger,
		WithLowercase(*lowercase),
		WithUppercase(*uppercase),
		WithCapitalize(*capitalize),
		WithCamelcase(*camelcase),
		WithSnakeCase(*snakeCase),
		WithKebabCase(*kebabCase),
		WithPascalCase(*pascalCase),
	)
	name, err := naming.Execute(strings.Join(args[0:], " "))
	if err != nil {
		logger.ErrorContext(ctx, "error executing naming", "error", err)
		return
	}
	fmt.Print(name)
}
