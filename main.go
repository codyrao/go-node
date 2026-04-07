package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const toolName = "go-node"
const toolVersion = "0.1.0"

type Config struct {
	InputFile   string
	OutputDir   string
	ModuleName  string
	PackageName string
	SourceDir   string
	NoCleanup   bool
	BuildDLL    bool
	ElectronVer string
	ShowVersion bool
}

// main parses command-line flags and executes the requested CLI action.
func main() {
	// Execute the CLI with the process-standard streams so normal shell usage stays unchanged.
	if err := execute(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// execute parses CLI arguments, handles version-only output, and runs the main build workflow.
func execute(args []string, stdout, stderr io.Writer) error {
	// Parse arguments first so version requests can short-circuit before build validation runs.
	config, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}

	// Print the pinned tool version immediately when the user requests version output.
	if config.ShowVersion {
		_, err := fmt.Fprintln(stdout, toolVersion)
		return err
	}

	// Delegate normal compilation requests to the existing build workflow.
	return run(config)
}

// parseFlags parses the provided command-line arguments into a CLI configuration.
func parseFlags(args []string, stderr io.Writer) (*Config, error) {
	// Use a dedicated FlagSet so tests can parse arguments without mutating global process flags.
	flagSet := flag.NewFlagSet(toolName, flag.ContinueOnError)
	flagSet.SetOutput(stderr)
	cfg := &Config{}
	flagSet.StringVar(&cfg.InputFile, "input", "", "Input Go source file path")
	flagSet.StringVar(&cfg.OutputDir, "output", "./output", "Output directory, default ./output")
	flagSet.StringVar(&cfg.ModuleName, "name", "", "Module name")
	flagSet.StringVar(&cfg.PackageName, "package", "main", "Go package name, default main (only functions in main package will be compiled)")
	flagSet.StringVar(&cfg.SourceDir, "source", "", "Go source file directory")
	flagSet.BoolVar(&cfg.NoCleanup, "no-cleanup", false, "Do not cleanup temporary files after compilation")
	flagSet.BoolVar(&cfg.BuildDLL, "dll", false, "Build DLL file instead of Node.js native module (for node-ffi)")
	flagSet.StringVar(&cfg.ElectronVer, "ev", "", "Electron version (e.g., 28.0.0). If not specified, uses node-gyp's default Node.js version")
	flagSet.BoolVar(&cfg.ShowVersion, "version", false, "Print go-node version and exit")

	flagSet.Usage = func() {
		// Render a stable help banner and examples so the CLI consistently presents the go-node name.
		fmt.Fprint(stderr, usageBanner())
		flagSet.PrintDefaults()
		fmt.Fprint(stderr, usageExamples())
	}

	// Parse the caller-provided arguments against the dedicated FlagSet.
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

// usageBanner returns the static CLI help header shown before the generated flag table.
func usageBanner() string {
	// Keep banner text centralized so tests can lock the public tool name to go-node.
	return `go-node - Compile Go code to Node.js native module

Usage:
  go-node [options]

Options:
`
}

// usageExamples returns the static example block appended after the flag table in CLI help output.
func usageExamples() string {
	// Keep examples centralized so command-name changes stay synchronized across help text.
	return `
Examples:
  go-node -input=hello.go -name=hello
  go-node -input=hello.go -name=hello -output=./dist
  go-node -input=hello.go -name=hello -dll
`
}

// run validates build configuration and executes the go-node compilation workflow.
func run(cfg *Config) error {
	// Allow version-only calls to bypass all build validation even when run is invoked directly in tests.
	if cfg.ShowVersion {
		fmt.Println(toolVersion)
		return nil
	}

	// Validate required user input before deriving paths or writing any build artifacts.
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

	if err := cleanupStaleBuildArtifacts(workDir); err != nil {
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
	isElectron := cfg.ElectronVer != ""
	if err := runNodeGyp(tmpDir, isElectron, cfg.ElectronVer); err != nil {
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

// generateHexArray converts raw DLL bytes into a C-style hex array declaration string.
func generateHexArray(data []byte) string {
	// Append one byte at a time so generated declarations remain deterministic and compact.
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

// init retains the existing strings package reference used by older builds.
func init() {
	// Keep a no-op strings call so the compiler preserves the import in this generated-tool entrypoint.
	strings.Contains("", "")
}
