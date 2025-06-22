package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"

	"github.com/eu-erwin/protobuf-compiler/cmd/naming"
)

func NewNamingCmd(
	ctx context.Context,
	logger *slog.Logger,
	opts ...naming.Opt,
) *cobra.Command {
	logger.InfoContext(ctx, "initializing naming")

	n := naming.NewNaming(ctx, logger, opts...)
	cmd := &cobra.Command{
		Use:   "naming",
		Short: "Helper cmd for naming",
		Long:  ``,
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			r, _ := n.Execute(strings.Join(args, " "))
			fmt.Printf("%v", r)
		},
	}
	cmd.Flags().BoolVarP(&n.PascalCase, "p", "p", false, "Output in Pascal case")
	cmd.Flags().BoolVarP(&n.Uppercase, "u", "u", false, "Output in Upper case")
	cmd.Flags().BoolVarP(&n.Lowercase, "l", "l", false, "Output in Lower case")
	cmd.Flags().BoolVarP(&n.Capitalize, "t", "t", false, "Output in Capitalize")
	cmd.Flags().BoolVarP(&n.SnakeCase, "s", "s", false, "Output in Snake case")
	cmd.Flags().BoolVarP(&n.KebabCase, "k", "k", false, "Output in Kebab case")
	cmd.Flags().BoolVarP(&n.CamelCase, "c", "c", false, "Output in Camel case")
	return cmd
}
