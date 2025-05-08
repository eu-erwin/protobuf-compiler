package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute(
	ctx context.Context,
	logger *slog.Logger,
) {
	versioningCmd := NewVersioningCmd(ctx, logger)
	namingCmd := NewNamingCmd(
		ctx,
		logger,
		WithLowercase(true),
		WithSnakeCase(true),
		WithKebabCase(true),
		WithCamelcase(true),
		WithCapitalize(true),
		WithUppercase(true),
		WithPascalCase(true),
	)
	rootCmd.AddCommand(
		versioningCmd,
		namingCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		logger.ErrorContext(ctx, "error executing root command", "error", err.Error())
		os.Exit(1)
	}
}
