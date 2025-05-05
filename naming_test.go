package main_test

import (
	"context"
	"log/slog"
	"testing"

	m "github.com/eu-erwin/protobuf-compiler"
)

func TestNaming(t *testing.T) {
	tests := []struct {
		name       string
		expected   string
		lowercase  bool
		capitalize bool
		camelcase  bool
		snakeCase  bool
		kebabCase  bool
		uppercase  bool
		pascalCase bool
	}{
		{name: "Hello World", expected: "Hello World"},
		{name: "Hello World", expected: "hello world", lowercase: true},
		{name: "Hello World", expected: "HELLO WORLD", uppercase: true},
		{name: "Hello World", expected: "helloWorld", camelcase: true},
		{name: "Hello World", expected: "HelloWorld", pascalCase: true},
		{name: "Hello World", expected: "Hello World", capitalize: true},
		{name: "Hello World", expected: "Hello_World", snakeCase: true},
		{name: "Hello World", expected: "Hello-World", kebabCase: true},
		{name: "Hello World", expected: "hello_world", lowercase: true, snakeCase: true},
		{name: "Hello World", expected: "hello-world", lowercase: true, kebabCase: true},
		{name: "Hello World", expected: "HELLO_WORLD", uppercase: true, snakeCase: true},
		{name: "Hello World", expected: "HELLO-WORLD", uppercase: true, kebabCase: true},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			n := m.NewNaming(
				context.Background(),
				slog.Default(),
				m.WithLowercase(tests[i].lowercase),
				m.WithSnakeCase(tests[i].snakeCase),
				m.WithKebabCase(tests[i].kebabCase),
				m.WithCamelcase(tests[i].camelcase),
				m.WithCapitalize(tests[i].capitalize),
				m.WithUppercase(tests[i].uppercase),
				m.WithPascalCase(tests[i].pascalCase),
			)
			result, err := n.Execute(tests[i].name)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tests[i].expected {
				t.Errorf("expected %q, got %q", tests[i].expected, result)
			}
		})
	}
}

func TestErrorNaming(t *testing.T) {
	tests := []struct {
		lowercase  bool
		capitalize bool
		camelcase  bool
		snakeCase  bool
		kebabCase  bool
		uppercase  bool
		pascalCase bool
	}{
		{lowercase: true, uppercase: true},
		{camelcase: true, pascalCase: true},
		{kebabCase: true, snakeCase: true},
	}

	for i := range tests {
		t.Run("Test"+string(rune(i)), func(t *testing.T) {
			n := m.NewNaming(
				context.Background(),
				slog.Default(),
				m.WithLowercase(tests[i].lowercase),
				m.WithSnakeCase(tests[i].snakeCase),
				m.WithKebabCase(tests[i].kebabCase),
				m.WithCamelcase(tests[i].camelcase),
				m.WithCapitalize(tests[i].capitalize),
				m.WithUppercase(tests[i].uppercase),
				m.WithPascalCase(tests[i].pascalCase),
			)
			_, err := n.Execute("test")
			if err == nil {
				t.Errorf("expected error, get nil")
			}
		})
	}
}
