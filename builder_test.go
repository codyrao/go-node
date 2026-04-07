package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateRuntimeHelperFileWritesFreeCString verifies runtime helper generation includes the FreeCString export.
func TestGenerateRuntimeHelperFileWritesFreeCString(t *testing.T) {
	// Prepare a temporary workspace for helper generation.
	workDir := t.TempDir()

	// Generate helper file content and read it back for assertions.
	helperPath, err := generateRuntimeHelperFile(workDir, "main")
	if err != nil {
		t.Fatalf("generateRuntimeHelperFile returned error: %v", err)
	}
	data, err := os.ReadFile(helperPath)
	if err != nil {
		t.Fatalf("failed to read generated runtime helper file: %v", err)
	}
	content := string(data)

	// Ensure expected exported functions exist so wrapper can free result memory safely.
	if !strings.Contains(content, "//export FreeCString") {
		t.Fatalf("expected helper file to export FreeCString")
	}
	if !strings.Contains(content, "func FreeCString(ptr *C.char)") {
		t.Fatalf("expected helper file to define FreeCString function")
	}
}

// TestCleanupBuildRemovesGeneratedRuntimeHelper ensures cleanup removes helper artifacts created for build-time linking.
func TestCleanupBuildRemovesGeneratedRuntimeHelper(t *testing.T) {
	// Create a fake build workspace with generated files that should be cleaned.
	workDir := t.TempDir()
	buildDir := filepath.Join(workDir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		t.Fatalf("failed to create build directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "callback.h"), []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to write callback header: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workDir, generatedRuntimeHelperName), []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to write runtime helper: %v", err)
	}

	// Run cleanup and verify generated artifacts are removed from disk.
	if err := cleanupBuild(workDir); err != nil {
		t.Fatalf("cleanupBuild returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workDir, "callback.h")); !os.IsNotExist(err) {
		t.Fatalf("expected callback.h to be removed, got err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(workDir, generatedRuntimeHelperName)); !os.IsNotExist(err) {
		t.Fatalf("expected generated runtime helper to be removed, got err=%v", err)
	}
	if _, err := os.Stat(buildDir); !os.IsNotExist(err) {
		t.Fatalf("expected build directory to be removed, got err=%v", err)
	}
}
