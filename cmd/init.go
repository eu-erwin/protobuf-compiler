package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/eu-erwin/protobuf-compiler/cmd/templating"
)

func NewInitCmd(
	ctx context.Context,
	logger *slog.Logger,
) *cobra.Command {
	logger.InfoContext(ctx, "initializing init cmd")
	p := &project{
		Logger:     logger,
		templating: templating.NewTemplating(ctx, logger),
	}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Helper cmd for initializing a new project",
		Long:  ``,
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
		Run:   p.Init,
	}
	cmd.Flags().StringVarP(&p.Git, "git", "g", "", "Git provider: gitlab, github")
	cmd.Flags().StringVarP(&p.Name, "name", "n", "", "Project name")
	cmd.Flags().StringVarP(&p.Organization, "organization", "o", "", "Git organization (e.g., your git repo organization)")
	return cmd
}

type project struct {
	Logger       *slog.Logger
	templating   *templating.Templating
	Git          string
	Name         string
	Organization string
}

func (c *project) Init(cmd *cobra.Command, args []string) {
	c.Logger.InfoContext(cmd.Context(), "initializing project")
	modifiers := templating.ModifierFactory(c.Logger, c.Name, c.Organization)

	switch c.Git {
	case "gitlab":
		c.templating.Execute(".gitlab-ci.yml", modifiers...)
	case "github":
		c.templating.Execute(".github/workflows/ci.yml", modifiers...)
	}
	c.templating.Execute("README.md", modifiers...)

	c.templating.Execute("message.proto", modifiers...)
	c.templating.Execute("services.proto", modifiers...)
}
