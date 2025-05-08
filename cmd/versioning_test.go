package cmd_test

import (
	"context"
	"log/slog"
	"testing"

	m "github.com/eu-erwin/protobuf-compiler/cmd"
)

func TestVersioning(t *testing.T) {
	tests := []struct {
		version  string
		expected m.VersionInfo
	}{
		{"0.0.0", m.VersionInfo{Type: m.Release, Major: 0, Minor: 0, Patch: 0}},
		{"0.0.1", m.VersionInfo{Type: m.Release, Major: 0, Minor: 0, Patch: 1}},
		{"0.1.0", m.VersionInfo{Type: m.Release, Major: 0, Minor: 1, Patch: 0}},
		{"1.0.0", m.VersionInfo{Type: m.Release, Major: 1, Minor: 0, Patch: 0}},
		{"1.0.0-alpha", m.VersionInfo{Type: m.Alpha, Major: 1, Minor: 0, Patch: 0}},
		{"1.0.0-beta", m.VersionInfo{Type: m.Beta, Major: 1, Minor: 0, Patch: 0}},
		{"1.0.0-release", m.VersionInfo{Type: m.Release, Major: 1, Minor: 0, Patch: 0}},
		{"1.1.0-Alpha", m.VersionInfo{Type: m.Alpha, Major: 1, Minor: 1, Patch: 0}},
		{"1.1.0-Beta", m.VersionInfo{Type: m.Beta, Major: 1, Minor: 1, Patch: 0}},
		{"1.1.0-Release", m.VersionInfo{Type: m.Release, Major: 1, Minor: 1, Patch: 0}},
		{"1.2.0-ALPHA", m.VersionInfo{Type: m.Alpha, Major: 1, Minor: 2, Patch: 0}},
		{"1.2.0-BETA", m.VersionInfo{Type: m.Beta, Major: 1, Minor: 2, Patch: 0}},
		{"1.2.0-RELEASE", m.VersionInfo{Type: m.Release, Major: 1, Minor: 2, Patch: 0}},
		{"v1.2.0", m.VersionInfo{Type: m.Release, Major: 1, Minor: 2, Patch: 0, Prefix: "v"}},
		{"ver1.2.0", m.VersionInfo{Type: m.Release, Major: 1, Minor: 2, Patch: 0, Prefix: "v"}},
		{"V1.2.0", m.VersionInfo{Type: m.Release, Major: 1, Minor: 2, Patch: 0, Prefix: "v"}},
	}

	for i := range tests {
		t.Run(tests[i].version, func(t *testing.T) {
			n := m.NewVersioning(
				context.Background(),
				slog.Default(),
			)
			result := n.Parse(tests[i].version)
			if result.Type != tests[i].expected.Type {
				t.Errorf("Type expected %q, got %q", tests[i].expected.Type, result.Type)
			}
			if result.Major != tests[i].expected.Major {
				t.Errorf("Major expected %q, got %q", tests[i].expected.Major, result.Major)
			}
			if result.Minor != tests[i].expected.Minor {
				t.Errorf("Minor expected %q, got %q", tests[i].expected.Minor, result.Minor)
			}
			if result.Patch != tests[i].expected.Patch {
				t.Errorf("Minor expected %q, got %q", tests[i].expected.Patch, result.Patch)
			}
			if result.Patch != tests[i].expected.Patch {
				t.Errorf("Minor expected %q, got %q", tests[i].expected.Patch, result.Patch)
			}
			if result.Prefix != tests[i].expected.Prefix {
				t.Errorf("Prefix expected %q, got %q", tests[i].expected.Prefix, result.Prefix)
			}
		})
	}
}
