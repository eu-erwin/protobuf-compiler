package naming

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
	opts ...Opt,
) *Naming {
	logger.InfoContext(ctx, "initializing naming")

	n := &Naming{Logger: logger}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

type Opt func(*Naming)

func WithLowercase(enabled bool) Opt {
	return func(n *Naming) {
		n.Lowercase = enabled
	}
}

func WithUppercase(enabled bool) Opt {
	return func(n *Naming) {
		n.Uppercase = enabled
	}
}

func WithCapitalize(enabled bool) Opt {
	return func(n *Naming) {
		n.Capitalize = enabled
	}
}

func WithCamelcase(enabled bool) Opt {
	return func(n *Naming) {
		n.CamelCase = enabled
	}
}

func WithSnakeCase(enabled bool) Opt {
	return func(n *Naming) {
		n.SnakeCase = enabled
	}
}

func WithKebabCase(enabled bool) Opt {
	return func(n *Naming) {
		n.KebabCase = enabled
	}
}

func WithPascalCase(enabled bool) Opt {
	return func(n *Naming) {
		n.PascalCase = enabled
	}
}

type Naming struct {
	Logger     *slog.Logger
	Lowercase  bool
	Uppercase  bool
	Capitalize bool
	CamelCase  bool
	SnakeCase  bool
	KebabCase  bool
	PascalCase bool
}

func (n *Naming) Execute(name string) (string, error) {
	if n.Lowercase && n.Uppercase {
		return "", fmt.Errorf("cannot use both Lowercase and Uppercase")
	}

	if n.SnakeCase && n.KebabCase {
		return "", fmt.Errorf("cannot use both snakecase and kebabcase")
	}

	if n.CamelCase && n.PascalCase {
		return "", fmt.Errorf("cannot use both camelcase and pascalcase")
	}

	if n.Lowercase {
		name = strings.ToLower(name)
	}

	if n.Uppercase {
		name = strings.ToUpper(name)
	}

	if n.Capitalize {
		name = cases.Title(language.English).String(name)
	}

	if n.CamelCase {
		words := strings.Fields(name)
		words[0] = strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			words[i] = cases.Title(language.English).String(words[i])
		}
		name = strings.Join(words, "")
	}

	if n.PascalCase {
		words := strings.Fields(name)
		for i := 0; i < len(words); i++ {
			words[i] = cases.Title(language.English).String(words[i])
		}
		name = strings.Join(words, "")
	}

	if n.SnakeCase {
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ReplaceAll(name, "-", "_")
	}

	if n.KebabCase {
		name = strings.ReplaceAll(name, " ", "-")
		name = strings.ReplaceAll(name, "_", "-")
	}

	return name, nil
}
