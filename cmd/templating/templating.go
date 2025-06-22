package templating

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/eu-erwin/protobuf-compiler/cmd/naming"
)

type TemplateModifier func(template string) string

type Templating struct {
	Logger       *slog.Logger
	Name         string
	Organization string
	TargetPath   string
}

func NewTemplating(
	_ context.Context,
	logger *slog.Logger,
) *Templating {
	logger.Info("initializing templating")
	tpl := &Templating{
		Logger:     logger,
		TargetPath: "/code",
	}
	return tpl
}

func (t *Templating) Execute(fileName string, modifiers ...TemplateModifier) string {
	filePath := fmt.Sprintf("/var/protobuf/template/%s.temp", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Logger.Error("template file does not exist", "file", filePath)
		return ""
	}

	template := ""
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Logger.Error("failed to read template file", "file", filePath, "error", err)
		return ""
	}

	template = string(content)
	for i := range modifiers {
		template = modifiers[i](template)
	}

	target := fmt.Sprintf("%s/%s", t.TargetPath, template)
	t.Logger.Info("writing template file", "file", target)
	if err := os.WriteFile(target, []byte(template), 0644); err != nil {
		t.Logger.Error("failed to write file", "file", target, "error", err)
		return ""
	}
	t.Logger.Info("template file written successfully", "file", target)
	return template
}

func ModifierFactory(logger *slog.Logger, name, organization string) []TemplateModifier {
	return []TemplateModifier{
		NameModifierFn(logger, name),
		PackageModifierFn(logger, name),
		FilenameModifierFn(logger, name),
		NamespaceModifierFn(logger, name),
		TitleModifierFn(logger, name),
		OrganizationModifierFn(logger, organization),
	}
}

func NameModifierFn(logger *slog.Logger, name string) TemplateModifier {
	exe := naming.Naming{Logger: logger, Lowercase: true, SnakeCase: true}
	value, err := exe.Execute(name)
	if err != nil {
		logger.Error("failed to convert name in name modifier", "error", err)
		os.Exit(1)
	}
	return func(content string) string {
		return modifier(content, "__name__", value)
	}
}

func PackageModifierFn(logger *slog.Logger, packageName string) TemplateModifier {
	exe := naming.Naming{Logger: logger, Lowercase: true, SnakeCase: true}
	packageName, err := exe.Execute(packageName)
	if err != nil {
		logger.Error("failed to convert package name in package modifier", "error", err)
		os.Exit(1)
	}
	return func(content string) string {
		return modifier(content, "__package__", packageName)
	}
}

func FilenameModifierFn(logger *slog.Logger, filename string) TemplateModifier {
	exe := naming.Naming{Logger: logger, Lowercase: true, SnakeCase: true}
	filename, err := exe.Execute(filename)
	if err != nil {
		logger.Error("failed to convert package name in namespace modifier", "error", err)
		os.Exit(1)
	}
	return func(content string) string {
		return modifier(content, "__filename__", filename)
	}
}

func NamespaceModifierFn(logger *slog.Logger, namespace string) TemplateModifier {
	exe := naming.Naming{Logger: logger, PascalCase: true}
	namespace, err := exe.Execute(namespace)
	if err != nil {
		logger.Error("failed to convert package name in namespace modifier", "error", err)
		os.Exit(1)
	}
	return func(content string) string {
		return modifier(content, "__namespace__", namespace)
	}
}

func TitleModifierFn(logger *slog.Logger, title string) TemplateModifier {
	exe := naming.Naming{Logger: logger, PascalCase: true}
	title, err := exe.Execute(title)
	if err != nil {
		logger.Error("failed to convert package name in namespace modifier", "error", err)
		os.Exit(1)
	}
	return func(content string) string {
		return modifier(content, "__capitalize__", title)
	}
}

func OrganizationModifierFn(logger *slog.Logger, organization string) TemplateModifier {
	exe := naming.Naming{Logger: logger, KebabCase: true}
	organization, err := exe.Execute(organization)
	if err != nil {
		logger.Error("failed to convert organization in namespace modifier", "error", err)
		os.Exit(1)
	}
	return func(content string) string {
		return modifier(content, "__organization__", organization)
	}
}

func modifier(content, replace, name string) string {
	return strings.ReplaceAll(content, replace, name)
}
