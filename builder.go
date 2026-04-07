package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const generatedRuntimeHelperName = "go_node_generated_runtime.go"
const generatedCallbackHeaderMarkerName = ".go-node-generated-callback-header"
const sampleCallbackSupportFileName = "callback_support.go"

// buildGoSharedLibrary compiles user Go code into a DLL and injects runtime helper exports used by the Node wrapper.
func buildGoSharedLibrary(cfg *Config, workDir string) error {
	// Normalize the work directory first so path operations stay deterministic.
	workDir = filepath.Clean(workDir)

	dllPath := filepath.Join(workDir, "build", cfg.ModuleName+".dll")

	// Resolve input file as a workDir-relative path because the command runs with cmd.Dir=workDir.
	goFile := cfg.InputFile
	if filepath.IsAbs(goFile) {
		relPath, err := filepath.Rel(workDir, goFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to calculate relative path: %v\n", err)
			return fmt.Errorf("calculate relative path failed: %w", err)
		}
		goFile = relPath
	}

	// Generate helper exports (for example FreeCString) so wrapper memory ownership stays in the Go DLL.
	helperPath, err := generateRuntimeHelperFile(workDir, cfg.PackageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate runtime helper file: %v\n", err)
		return fmt.Errorf("generate runtime helper file failed: %w", err)
	}
	helperFile := filepath.Base(helperPath)

	args := []string{
		"build",
		"-buildmode=c-shared",
		"-trimpath",
		"-buildvcs=false",
		"-ldflags", "-s -w -buildid=",
		"-o", dllPath,
	}

	// Resolve the final Go source list so sample support files can participate without changing the CLI surface.
	buildFiles, err := resolveBuildGoFiles(workDir, goFile, helperFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to resolve build Go files: %v\n", err)
		return fmt.Errorf("resolve build Go files failed: %w", err)
	}
	args = append(args, buildFiles...)

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"CC=gcc",
		"CXX=g++",
		"GOOS=windows",
		"GOARCH=amd64",
		"GO111MODULE=on",
		"GOPROXY=direct",
		"CGO_CFLAGS=-I. -Os -ffunction-sections -fdata-sections",
		"CGO_LDFLAGS=-s -w -Wl,--gc-sections",
	)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Compiling Go code to shared library: %s\n", dllPath)
	fmt.Printf("Working directory: %s\n", workDir)
	fmt.Printf("Go file: %s\n", goFile)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to compile Go shared library: %v\n", err)
		return fmt.Errorf("compile Go shared library failed: %w", err)
	}

	if _, err := os.Stat(dllPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Shared library file not generated: %s\n", dllPath)
		return fmt.Errorf("shared library file not generated: %s", dllPath)
	}

	// 使用UPX压缩DLL
	fmt.Printf("Compressing DLL with UPX...\n")
	upxCmd := exec.Command("upx", "--best", "--lzma", dllPath)
	upxCmd.Stdout = os.Stdout
	upxCmd.Stderr = os.Stderr
	if err := upxCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: UPX compression failed (may not be installed): %v\n", err)
		fmt.Fprintf(os.Stderr, "Continuing without UPX compression...\n")
	}

	return nil
}

// resolveBuildGoFiles returns the ordered Go files that should participate in a shared-library build.
func resolveBuildGoFiles(workDir, goFile, helperFile string) ([]string, error) {
	// Start with the selected input file because it defines the user-facing entrypoint for this build.
	buildFiles := []string{goFile}

	// Include the optional sample support file when it exists so checked-in examples compile as one package.
	supportPath := filepath.Join(workDir, sampleCallbackSupportFileName)
	if _, err := os.Stat(supportPath); err == nil {
		if sampleCallbackSupportFileName != goFile && sampleCallbackSupportFileName != helperFile {
			buildFiles = append(buildFiles, sampleCallbackSupportFileName)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Add the generated runtime helper last so its exported cleanup helpers are always compiled in.
	if helperFile != goFile && helperFile != sampleCallbackSupportFileName {
		buildFiles = append(buildFiles, helperFile)
	}

	return buildFiles, nil
}

// generateRuntimeHelperFile writes a temporary Go source file that exports helpers required by generated wrapper code.
func generateRuntimeHelperFile(workDir, packageName string) (string, error) {
	// Keep package naming aligned with user input while falling back to main for safety.
	if strings.TrimSpace(packageName) == "" {
		packageName = "main"
	}

	// Keep helper minimal to avoid inflating binary size while providing explicit free semantics.
	content := `package ` + packageName + `

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// FreeCString releases a C string allocated inside the Go shared library.
//
//export FreeCString
func FreeCString(ptr *C.char) {
	// Guard against nil pointers to keep release operation idempotent.
	if ptr == nil {
		return
	}

	// Free inside the same DLL allocator boundary to avoid cross-runtime deallocation issues.
	C.free(unsafe.Pointer(ptr))
}
`

	helperPath := filepath.Join(workDir, generatedRuntimeHelperName)
	if err := os.WriteFile(helperPath, []byte(content), 0644); err != nil {
		return "", err
	}
	return helperPath, nil
}

// createBuildDirectory ensures the build output directory exists before compilation starts.
func createBuildDirectory(workDir string) error {
	// Create the build directory eagerly so later steps can write outputs without repeated checks.
	buildDir := filepath.Join(workDir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("Failed to create build directory: %w", err)
	}
	return nil
}

// cleanupStaleBuildArtifacts removes transient build outputs left behind by earlier runs while preserving checked-in support files.
func cleanupStaleBuildArtifacts(workDir string) error {
	// Reuse the build-directory cleanup so stale DLLs and intermediate files never leak into the next build.
	buildDir := filepath.Join(workDir, "build")
	if _, err := os.Stat(buildDir); err == nil {
		if err := os.RemoveAll(buildDir); err != nil {
			return err
		}
	}

	// Remove stale callback ownership markers without deleting a checked-in callback header that now owns the path.
	callbackHeaderMarker := filepath.Join(workDir, generatedCallbackHeaderMarkerName)
	if err := os.Remove(callbackHeaderMarker); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Remove generated runtime helpers so the next build always recreates them from the current source.
	runtimeHelperPath := filepath.Join(workDir, generatedRuntimeHelperName)
	if _, err := os.Stat(runtimeHelperPath); err == nil {
		if err := os.Remove(runtimeHelperPath); err != nil {
			return err
		}
	}

	return nil
}

// cleanupBuild removes generated build artifacts that should not persist across runs.
func cleanupBuild(workDir string) error {
	// Remove the build directory so stale DLLs and intermediate files do not contaminate the next run.
	buildDir := filepath.Join(workDir, "build")
	if _, err := os.Stat(buildDir); err == nil {
		if err := os.RemoveAll(buildDir); err != nil {
			return err
		}
	}

	// Remove generated callback header from source directory only when a build marker proves ownership.
	callbackHeader := filepath.Join(workDir, "callback.h")
	callbackHeaderMarker := filepath.Join(workDir, generatedCallbackHeaderMarkerName)
	if _, err := os.Stat(callbackHeaderMarker); err == nil {
		if _, err := os.Stat(callbackHeader); err == nil {
			if err := os.Remove(callbackHeader); err != nil {
				return err
			}
		}
		if err := os.Remove(callbackHeaderMarker); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	// Remove generated runtime helper file from source directory.
	runtimeHelperPath := filepath.Join(workDir, generatedRuntimeHelperName)
	if _, err := os.Stat(runtimeHelperPath); err == nil {
		if err := os.Remove(runtimeHelperPath); err != nil {
			return err
		}
	}

	return nil
}

// copyFile copies one file to another location without loading the full file into memory.
func copyFile(src, dst string) error {
	// Open source first so early failures do not create destination files.
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file and stream bytes directly from source.
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	return err
}

// fileExists reports whether the path currently exists.
func fileExists(path string) bool {
	// Rely on os.Stat to keep existence checks cheap and platform-independent.
	_, err := os.Stat(path)
	return err == nil
}

// getFileExt returns the lowercase file extension of a path.
func getFileExt(path string) string {
	// Normalize extension casing to simplify extension comparisons across platforms.
	ext := filepath.Ext(path)
	return strings.ToLower(ext)
}
