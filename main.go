package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	InputFile   string
	OutputDir   string
	ModuleName  string
	PackageName string
	SourceDir   string
	NoCleanup   bool
	BuildDLL    bool
}

func main() {
	config := parseFlags()

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.InputFile, "input", "", "Input Go source file path")
	flag.StringVar(&cfg.OutputDir, "output", "./output", "Output directory, default ./output")
	flag.StringVar(&cfg.ModuleName, "name", "", "Module name")
	flag.StringVar(&cfg.PackageName, "package", "main", "Go package name, default main (only functions in main package will be compiled)")
	flag.StringVar(&cfg.SourceDir, "source", "", "Go source file directory")
	flag.BoolVar(&cfg.NoCleanup, "no-cleanup", false, "Do not cleanup temporary files after compilation")
	flag.BoolVar(&cfg.BuildDLL, "dll", false, "Build DLL file instead of Node.js native module (for node-ffi)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Go2Node - Compile Go code to Node.js native module

Usage:
  go2node [options]

Options:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  go2node -input=hello.go -name=hello
  go2node -input=hello.go -name=hello -output=./dist
  go2node -input=hello.go -name=hello -dll
`)
	}

	flag.Parse()

	return cfg
}

func run(cfg *Config) error {
	if cfg.InputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Please specify input Go source file (-input)")
		return fmt.Errorf("input file not specified")
	}

	if cfg.ModuleName == "" {
		fmt.Fprintln(os.Stderr, "Error: Please specify module name (-name)")
		return fmt.Errorf("module name not specified")
	}

	if cfg.SourceDir == "" {
		cfg.SourceDir = filepath.Dir(cfg.InputFile)
		if cfg.SourceDir == "." {
			cfg.SourceDir, _ = os.Getwd()
		}
	}

	absInputFile, _ := filepath.Abs(cfg.InputFile)
	workDir := cfg.SourceDir

	if filepath.IsAbs(workDir) {
		workDir = filepath.Clean(workDir)
	} else {
		workDirAbs, _ := filepath.Abs(workDir)
		workDir = workDirAbs
	}

	cfg.SourceDir = workDir

	baseName := filepath.Base(cfg.InputFile)
	if cfg.SourceDir == filepath.Dir(absInputFile) {
		cfg.InputFile = baseName
	}

	outputNodeDir := cfg.OutputDir
	if err := os.MkdirAll(outputNodeDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create output directory: %v\n", err)
		return fmt.Errorf("create output directory failed: %w", err)
	}

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate random ID: %v\n", err)
		return fmt.Errorf("generate random ID failed: %w", err)
	}
	randomID := hex.EncodeToString(randomBytes)
	tmpDir := filepath.Join(cfg.SourceDir, "go-node-tmp_"+randomID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create temp directory: %v\n", err)
		return fmt.Errorf("create temp directory failed: %w", err)
	}

	defer func() {
		if !cfg.NoCleanup {
			fmt.Println("Cleaning up temporary files...")
			if err := cleanupBuild(workDir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to cleanup build directory: %v\n", err)
			}
			if err := os.RemoveAll(tmpDir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to cleanup temp directory: %v\n", err)
			} else {
				fmt.Printf("Deleted temp directory: %s\n", tmpDir)
			}
		}
	}()

	workDir, _ = filepath.Abs(workDir)
	tmpDirAbs, _ := filepath.Abs(tmpDir)
	fmt.Printf("Working directory: %s\n", workDir)
	fmt.Printf("Temporary directory: %s\n\n", tmpDirAbs)

	if err := cleanupBuild(workDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to cleanup build directory: %v\n", err)
		return fmt.Errorf("cleanup build directory failed: %w", err)
	}

	fmt.Println("Step 1/4: Generating callback header and compiling Go code...")
	if err := createBuildDirectory(workDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create build directory: %v\n", err)
		return fmt.Errorf("create build directory failed: %w", err)
	}
	// Generate callback.h in source directory for Go compilation
	if err := generateCallbackHeader(workDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate callback header: %v\n", err)
		return fmt.Errorf("generate callback header failed: %w", err)
	}
	if err := buildGoSharedLibrary(cfg, workDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	fmt.Println()

	fmt.Println("Step 2/4: Parsing Go export functions...")
	functions, err := parseGoExports(absInputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse Go export functions: %v\n", err)
		return fmt.Errorf("parse Go export functions failed: %w", err)
	}
	if len(functions) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No export functions found (use //export to mark functions)")
		return fmt.Errorf("no export functions found")
	}
	fmt.Printf("Found export functions: %v\n\n", functions)

	if cfg.BuildDLL {
		// 直接输出 DLL 文件（用于 node-ffi）
		fmt.Println("Step 3/3: Copying DLL file...")
		dllPath := filepath.Join(workDir, "build", cfg.ModuleName+".dll")
		outputDLLPath := filepath.Join(outputNodeDir, cfg.ModuleName+".dll")

		if err := copyFile(dllPath, outputDLLPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to copy DLL to output directory: %v\n", err)
			return err
		}
		fmt.Printf("Copied DLL file to: %s\n", outputDLLPath)
		fmt.Println()

		fmt.Println("========================")
		fmt.Println("Compilation completed!")
		fmt.Println("Output files:")
		fmt.Printf("  - %s (for node-ffi)\n", outputDLLPath)
		fmt.Println("========================")

		return nil
	}

	fmt.Println("Step 3/4: Generating binding files...")
	dllPath := filepath.Join(workDir, "build", cfg.ModuleName+".dll")

	// Copy DLL to temp directory for resource embedding
	tmpDLLPath := filepath.Join(tmpDir, cfg.ModuleName+".dll")
	if err := copyFile(dllPath, tmpDLLPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to copy DLL to temp directory: %v\n", err)
		return err
	}

	if err := generateResourceRC(tmpDir, cfg.ModuleName, tmpDLLPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if err := generateWrapperCC(tmpDir, cfg.ModuleName, functions, ""); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if err := generateBindingGyp(tmpDir, cfg.ModuleName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Println()

	fmt.Println("Step 4/4: Compiling with node-gyp...")
	if err := runNodeGyp(tmpDir, true); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	if err := buildWithNodeGyp(tmpDir, outputNodeDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	fmt.Println()

	fmt.Println("========================")
	fmt.Println("Compilation completed!")
	fmt.Println("Output files:")
	fmt.Printf("  - %s\n", filepath.Join(outputNodeDir, cfg.ModuleName+".node"))
	fmt.Println("========================")

	return nil
}

func generateHexArray(data []byte) string {
	var hex strings.Builder
	hex.WriteString("unsigned char dllData[] = {")
	for i, b := range data {
		if i > 0 {
			hex.WriteString(",")
		}
		hex.WriteString(fmt.Sprintf("0x%02X", b))
	}
	hex.WriteString("};")
	return hex.String()
}

func init() {
	strings.Contains("", "")
}
