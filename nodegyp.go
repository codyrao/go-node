package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildWithNodeGyp(workDir, outputNodeDir string) error {
	releaseDir := filepath.Join(workDir, "build", "Release")

	if _, err := os.Stat(releaseDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Release directory does not exist: %s\n", releaseDir)
		return fmt.Errorf("Release directory does not exist: %s", releaseDir)
	}

	files, err := os.ReadDir(releaseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read Release directory: %v\n", err)
		return fmt.Errorf("read Release directory failed: %w", err)
	}

	var nodeFileName string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".node") {
			nodeFileName = f.Name()
			break
		}
	}

	if nodeFileName == "" {
		fmt.Fprintln(os.Stderr, "Error: .node file not found")
		return fmt.Errorf(".node file not found")
	}

	srcNode := filepath.Join(releaseDir, nodeFileName)
	dstNode := filepath.Join(outputNodeDir, nodeFileName)

	if err := copyFile(srcNode, dstNode); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to copy .node file: %v\n", err)
		return fmt.Errorf("copy .node file failed: %w", err)
	}

	fmt.Printf("Copied .node file to: %s\n", dstNode)
	return nil
}

func runNodeGyp(workDir string, isElectron bool, electronVersion string) error {
	nodeGypPath := "C:\\Users\\admin\\AppData\\Roaming\\npm\\node-gyp.cmd"

	var targetVersion, distURL string
	var runtime string

	if isElectron {
		if electronVersion != "" {
			targetVersion = electronVersion
		} else {
			targetVersion = "39.0.0"
		}
		distURL = "https://electronjs.org/headers"
		runtime = "electron"
	} else {
		targetVersion = "24.12.0"
		distURL = "https://nodejs.org/dist"
		runtime = "node"
	}

	args := []string{
		"configure",
		"--builddir=build",
		"--target=" + targetVersion,
		"--dist-url=" + distURL,
		"--runtime=" + runtime,
	}

	cmd := exec.Command(nodeGypPath, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running node-gyp configure for %s %s...\n", runtime, targetVersion)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: node-gyp configure failed: %v\n", err)
		return fmt.Errorf("node-gyp configure failed: %w", err)
	}

	args = []string{
		"build",
		"--target=" + targetVersion,
		"--dist-url=" + distURL,
		"--runtime=" + runtime,
	}

	cmd = exec.Command(nodeGypPath, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running node-gyp build...")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: node-gyp build failed: %v\n", err)
		return fmt.Errorf("node-gyp build failed: %w", err)
	}

	return nil
}

func copyDLLToOutput(srcDLL, dstPath string) error {
	dstDir := filepath.Dir(dstPath)
	if dstDir == "" {
		dstDir = "."
	}
	dstDLL := filepath.Join(dstDir, filepath.Base(srcDLL))

	if err := copyFile(srcDLL, dstDLL); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to copy DLL file: %v\n", err)
		return fmt.Errorf("copy DLL file failed: %w", err)
	}

	fmt.Printf("Copied DLL file to: %s\n", dstDLL)
	return nil
}
