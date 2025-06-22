package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/eu-erwin/protobuf-compiler/cmd/templating"
)

func NewTemplatingCmd(
	ctx context.Context,
	logger *slog.Logger,
) *cobra.Command {
	logger.InfoContext(ctx, "initializing templating cmd")
	t := templating.NewTemplating(ctx, logger)
	cmd := &cobra.Command{
		Use:   "templating",
		Short: "Helper cmd for templating",
		Long:  ``,
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.Flags().StringVarP(&t.Name, "name", "n", "", "Name (Generic Name)")
	cmd.Flags().StringVarP(&t.Organization, "organization", "o", "", "Organization (your git repo organization)")
	return cmd
}
