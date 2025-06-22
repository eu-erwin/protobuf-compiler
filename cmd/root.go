package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/eu-erwin/protobuf-compiler/cmd/naming"
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
	initCmd := NewInitCmd(ctx, logger)
	templatingCmd := NewTemplatingCmd(ctx, logger)
	versioningCmd := NewVersioningCmd(ctx, logger)
	namingCmd := NewNamingCmd(
		ctx,
		logger,
		naming.WithLowercase(true),
		naming.WithSnakeCase(true),
		naming.WithKebabCase(true),
		naming.WithCamelcase(true),
		naming.WithCapitalize(true),
		naming.WithUppercase(true),
		naming.WithPascalCase(true),
	)
	rootCmd.AddCommand(
		initCmd,
		templatingCmd,
		versioningCmd,
		namingCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		logger.ErrorContext(ctx, "error executing root command", "error", err.Error())
		os.Exit(1)
	}
}
