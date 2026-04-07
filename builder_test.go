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
	if err := os.WriteFile(filepath.Join(workDir, generatedCallbackHeaderMarkerName), []byte("generated"), 0644); err != nil {
		t.Fatalf("failed to write callback header marker: %v", err)
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

// TestCleanupBuildPreservesCheckedInCallbackHeader ensures cleanup does not delete a pre-existing callback header owned by the repository.
func TestCleanupBuildPreservesCheckedInCallbackHeader(t *testing.T) {
	// Create a workspace that simulates a checked-in callback header with no generated marker file.
	workDir := t.TempDir()
	headerPath := filepath.Join(workDir, "callback.h")
	if err := os.WriteFile(headerPath, []byte("checked-in"), 0644); err != nil {
		t.Fatalf("failed to write callback header: %v", err)
	}

	// Run cleanup and verify the pre-existing callback header remains on disk.
	if err := cleanupBuild(workDir); err != nil {
		t.Fatalf("cleanupBuild returned error: %v", err)
	}
	data, err := os.ReadFile(headerPath)
	if err != nil {
		t.Fatalf("expected callback.h to remain readable, got err=%v", err)
	}
	if string(data) != "checked-in" {
		t.Fatalf("unexpected callback.h content after cleanup: %q", string(data))
	}
}

// TestCleanupBuildRemovesGeneratedCallbackHeader ensures cleanup removes callback headers created by the build pipeline.
func TestCleanupBuildRemovesGeneratedCallbackHeader(t *testing.T) {
	// Create a workspace that simulates a generated callback header paired with its ownership marker.
	workDir := t.TempDir()
	headerPath := filepath.Join(workDir, "callback.h")
	markerPath := filepath.Join(workDir, generatedCallbackHeaderMarkerName)
	if err := os.WriteFile(headerPath, []byte("generated"), 0644); err != nil {
		t.Fatalf("failed to write callback header: %v", err)
	}
	if err := os.WriteFile(markerPath, []byte("generated"), 0644); err != nil {
		t.Fatalf("failed to write callback header marker: %v", err)
	}

	// Run cleanup and verify both generated artifacts are removed together.
	if err := cleanupBuild(workDir); err != nil {
		t.Fatalf("cleanupBuild returned error: %v", err)
	}
	if _, err := os.Stat(headerPath); !os.IsNotExist(err) {
		t.Fatalf("expected generated callback.h to be removed, got err=%v", err)
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("expected generated callback header marker to be removed, got err=%v", err)
	}
}

// TestResolveBuildGoFilesIncludesCallbackSupport verifies sample support files are compiled alongside the selected sample input.
func TestResolveBuildGoFilesIncludesCallbackSupport(t *testing.T) {
	// Create a workspace with one input file, one callback support file, and one generated helper file.
	workDir := t.TempDir()
	inputFile := "hello.go"
	supportFile := sampleCallbackSupportFileName
	helperFile := generatedRuntimeHelperName
	for _, fileName := range []string{inputFile, supportFile, helperFile} {
		if err := os.WriteFile(filepath.Join(workDir, fileName), []byte("package main"), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", fileName, err)
		}
	}

	// Resolve build inputs and assert the callback support file is preserved between the main input and helper file.
	files, err := resolveBuildGoFiles(workDir, inputFile, helperFile)
	if err != nil {
		t.Fatalf("resolveBuildGoFiles returned error: %v", err)
	}
	expected := []string{inputFile, supportFile, helperFile}
	if strings.Join(files, ",") != strings.Join(expected, ",") {
		t.Fatalf("unexpected build files: got %v want %v", files, expected)
	}
}

// TestCleanupStaleBuildArtifactsPreservesCheckedInCallbackHeader ensures pre-build cleanup removes stale ownership markers without deleting repository-owned headers.
func TestCleanupStaleBuildArtifactsPreservesCheckedInCallbackHeader(t *testing.T) {
	// Create a workspace with a checked-in callback header and a stale generated marker from an earlier build.
	workDir := t.TempDir()
	headerPath := filepath.Join(workDir, "callback.h")
	markerPath := filepath.Join(workDir, generatedCallbackHeaderMarkerName)
	if err := os.WriteFile(headerPath, []byte("checked-in"), 0644); err != nil {
		t.Fatalf("failed to write callback header: %v", err)
	}
	if err := os.WriteFile(markerPath, []byte("generated"), 0644); err != nil {
		t.Fatalf("failed to write callback header marker: %v", err)
	}

	// Run pre-build cleanup and verify only the stale marker is removed.
	if err := cleanupStaleBuildArtifacts(workDir); err != nil {
		t.Fatalf("cleanupStaleBuildArtifacts returned error: %v", err)
	}
	data, err := os.ReadFile(headerPath)
	if err != nil {
		t.Fatalf("expected callback.h to remain readable, got err=%v", err)
	}
	if string(data) != "checked-in" {
		t.Fatalf("unexpected callback.h content after stale cleanup: %q", string(data))
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("expected stale callback marker to be removed, got err=%v", err)
	}
}
