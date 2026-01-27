package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildGoSharedLibrary(cfg *Config, workDir string) error {
	workDir = filepath.Clean(workDir)

	dllPath := filepath.Join(workDir, "build", cfg.ModuleName+".dll")

	goFile := cfg.InputFile
	if filepath.IsAbs(goFile) {
		relPath, err := filepath.Rel(workDir, goFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to calculate relative path: %v\n", err)
			return fmt.Errorf("calculate relative path failed: %w", err)
		}
		goFile = relPath
	}

	args := []string{
		"build",
		"-buildmode=c-shared",
		"-o", dllPath,
	}

	args = append(args, goFile)

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"CC=gcc",
		"CXX=g++",
		"GOOS=windows",
		"GOARCH=amd64",
		"GO111MODULE=on",
		"GOPROXY=direct",
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

	return nil
}

func createBuildDirectory(workDir string) error {
	buildDir := filepath.Join(workDir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("Failed to create build directory: %w", err)
	}
	return nil
}

func cleanupBuild(workDir string) error {
	buildDir := filepath.Join(workDir, "build")
	if _, err := os.Stat(buildDir); err == nil {
		return os.RemoveAll(buildDir)
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	return err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getFileExt(path string) string {
	ext := filepath.Ext(path)
	return strings.ToLower(ext)
}
