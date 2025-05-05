package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	appName    = "naming-cli"
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

func NewNaming(opts ...NamingOpt) *Naming {
	n := &Naming{}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

type NamingOpt func(*Naming)

func WithLowercase(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.lowercase = enabled
	}
}

func WithUppercase(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.uppercase = enabled
	}
}

func WithCapitalize(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.capitalize = enabled
	}
}

func WithCamelcase(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.camelCase = enabled
	}
}

func WithSnakeCase(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.snakeCase = enabled
	}
}

func WithKebabCase(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.kebabCase = enabled
	}
}

func WithPascalCase(enabled bool) NamingOpt {
	return func(n *Naming) {
		n.pascalCase = enabled
	}
}

type Naming struct {
	lowercase  bool
	uppercase  bool
	capitalize bool
	camelCase  bool
	snakeCase  bool
	kebabCase  bool
	pascalCase bool
}

func (n *Naming) Execute(name string) (string, error) {
	if n.lowercase && n.uppercase {
		return "", fmt.Errorf("cannot use both lowercase and uppercase")
	}

	if n.snakeCase && n.kebabCase {
		return "", fmt.Errorf("cannot use both snakecase and kebabcase")
	}

	if n.camelCase && n.pascalCase {
		return "", fmt.Errorf("cannot use both camelcase and pascalcase")
	}

	if n.lowercase {
		name = strings.ToLower(name)
	}

	if n.uppercase {
		name = strings.ToUpper(name)
	}

	if n.capitalize {
		name = cases.Title(language.English).String(name)
	}

	if n.camelCase {
		words := strings.Fields(name)
		words[0] = strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			words[i] = cases.Title(language.English).String(words[i])
		}
		name = strings.Join(words, "")
	}

	if n.pascalCase {
		words := strings.Fields(name)
		for i := 0; i < len(words); i++ {
			words[i] = cases.Title(language.English).String(words[i])
		}
		name = strings.Join(words, "")
	}

	if n.snakeCase {
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ReplaceAll(name, "-", "_")
	}

	if n.kebabCase {
		name = strings.ReplaceAll(name, " ", "-")
		name = strings.ReplaceAll(name, "_", "-")
	}

	return name, nil
}
