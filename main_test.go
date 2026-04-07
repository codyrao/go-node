package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestParseFlagsSupportsVersionFlags verifies both -version and --version map to the same CLI behavior.
func TestParseFlagsSupportsVersionFlags(t *testing.T) {
	// Iterate over the supported version flag spellings so the CLI remains compatible with common shell usage.
	for _, args := range [][]string{{"-version"}, {"--version"}} {
		var stderr bytes.Buffer

		// Parse each argument form independently and assert the version mode is enabled without extra build inputs.
		cfg, err := parseFlags(args, &stderr)
		if err != nil {
			t.Fatalf("parseFlags(%v) returned error: %v", args, err)
		}
		if !cfg.ShowVersion {
			t.Fatalf("expected ShowVersion for args %v", args)
		}
		if cfg.InputFile != "" || cfg.ModuleName != "" {
			t.Fatalf("expected version args %v to avoid build requirements, got input=%q name=%q", args, cfg.InputFile, cfg.ModuleName)
		}
	}
}

// TestExecutePrintsVersion verifies version requests print the pinned tool version and skip build validation.
func TestExecutePrintsVersion(t *testing.T) {
	// Capture stdout and stderr so the test can assert the exact version-only CLI output.
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Execute the version path and verify it exits successfully without demanding build flags.
	if err := execute([]string{"--version"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute(--version) returned error: %v", err)
	}
	if stdout.String() != toolVersion+"\n" {
		t.Fatalf("unexpected version output: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for version output, got %q", stderr.String())
	}
}

// TestUsageTextUsesGoNode verifies CLI help text no longer exposes the legacy go2node name.
func TestUsageTextUsesGoNode(t *testing.T) {
	// Build the static help fragments so public CLI naming stays consistent in usage and examples.
	banner := usageBanner()
	examples := usageExamples()

	// Assert the new public name is present everywhere and the legacy branding is gone.
	if !strings.Contains(banner, "go-node") || !strings.Contains(examples, "go-node") {
		t.Fatalf("expected usage text to mention go-node, got banner=%q examples=%q", banner, examples)
	}
	if strings.Contains(strings.ToLower(banner+examples), "go2node") {
		t.Fatalf("expected usage text to exclude go2node, got %q", banner+examples)
	}
}
