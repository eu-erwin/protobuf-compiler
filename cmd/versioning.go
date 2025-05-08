package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type VersionType int

const (
	Alpha VersionType = iota
	Beta
	Release
)

type VersionInfo struct {
	Name   string
	Type   VersionType
	Prefix string
	Major  int
	Minor  int
	Patch  int
}

func NewVersioning(
	ctx context.Context,
	logger *slog.Logger,
) *Versioning {
	logger.InfoContext(ctx, "initializing versioning")
	return &Versioning{
		logger: logger,
		Parsers: []VersionDataParser{
			ParseVersionPrefix,
			ParseVersionNumber,
			ParseVersionType,
		},
	}
}

func NewVersioningCmd(
	ctx context.Context,
	logger *slog.Logger,
) *cobra.Command {
	logger.InfoContext(ctx, "initializing versioning cmd")
	versioning := NewVersioning(ctx, logger)
	c := &cobra.Command{
		Use:   "versioning",
		Short: "Helper cmd for versioning",
		Long:  ``,
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			info := versioning.Parse("1.0.0")
			fmt.Printf("%v", info.Major)
			fmt.Printf("%v", args)
		},
	}
	return c
}

type VersioningOpt func(*Versioning)

type VersionDataParser func(info *VersionInfo, version string)

type Versioning struct {
	logger  *slog.Logger
	Parsers []VersionDataParser
}

func (v *Versioning) Parse(version string) VersionInfo {
	data := VersionInfo{}
	for _, parser := range v.Parsers {
		parser(&data, version)
	}
	return data
}

func ParseVersionPrefix(data *VersionInfo, version string) {
	if version == "" {
		return
	}

	if version[0] == 'v' || version[0] == 'V' {
		data.Prefix = "v"
	}

	if strings.HasPrefix(version, "ver") {
		data.Prefix = "v"
	}

	if strings.HasPrefix(version, "Ver") {
		data.Prefix = "v"
	}
}

func ParseVersionNumber(data *VersionInfo, version string) {
	if version == "" {
		return
	}

	parts := strings.Split(SanitizeVersion(version), ".")
	for i := range parts {
		if i == 0 {
			data.Major = AdaptVersionNumber(parts[i])
		} else if i == 1 {
			data.Minor = AdaptVersionNumber(parts[i])
		} else if i == 2 {
			data.Patch = AdaptVersionNumber(parts[i])
		}
	}
}

func AdaptVersionNumber(number string) int {
	result, err := strconv.Atoi(number)
	if err != nil {
		return 0
	}
	return result
}

func ParseVersionType(data *VersionInfo, version string) {
	data.Type = Release
	if version == "" {
		return
	}

	separators := []string{"-", "+", "_"}
	for i := range separators {
		if strings.Contains(version, separators[i]) {
			parts := strings.Split(version, separators[i])
			if len(parts) > 1 {
				data.Type = AdaptVersionType(parts[1])
			}
		}
	}
}

func AdaptVersionType(versionType string) VersionType {
	switch strings.ToLower(versionType) {
	case "alpha":
		return Alpha
	case "beta":
		return Beta
	case "release":
		return Release
	default:
		return Release
	}
}

func SanitizeVersion(version string) string {
	if version == "" {
		return ""
	}

	if strings.Contains(version, "-") {
		parts := strings.Split(version, "-")
		version = parts[0]
	}

	if strings.Contains(version, "+") {
		parts := strings.Split(version, "+")
		version = parts[0]
	}

	if strings.HasPrefix(version, "ver") {
		version = version[3:]
	}

	if strings.HasPrefix(version, "Ver") {
		version = version[3:]
	}

	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	if strings.HasPrefix(version, "V") {
		version = version[1:]
	}

	return version
}
