package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NewNaming(
	ctx context.Context,
	logger *slog.Logger,
	opts ...NamingOpt,
) *Naming {
	logger.InfoContext(ctx, "initializing naming")

	n := &Naming{logger: logger}
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
	logger     *slog.Logger
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
