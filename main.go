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
	flag.StringVar(&cfg.InputFile, "input", "", "输入的Go源文件路径")
	flag.StringVar(&cfg.OutputDir, "output", "./output", "输出目录，默认./output")
	flag.StringVar(&cfg.ModuleName, "name", "", "模块名称")
	flag.StringVar(&cfg.PackageName, "package", "main", "Go包名，默认main（仅main包下的函数会被编译）")
	flag.StringVar(&cfg.SourceDir, "source", "", "Go源文件所在目录")
	flag.BoolVar(&cfg.NoCleanup, "no-cleanup", false, "编译后不清理临时文件")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Go2Node - 将Go代码编译为Node.js原生模块

用法:
  go2node [选项]

选项:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
示例:
  go2node -input=hello.go -name=hello
  go2node -input=hello.go -name=hello -output=./dist
`)
	}

	flag.Parse()

	return cfg
}

func run(cfg *Config) error {
	if cfg.InputFile == "" {
		return fmt.Errorf("请指定输入的Go源文件（-input）")
	}

	if cfg.ModuleName == "" {
		return fmt.Errorf("请指定模块名称（-name）")
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
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Errorf("生成随机ID失败: %w", err)
	}
	randomID := hex.EncodeToString(randomBytes)
	tmpDir := filepath.Join(cfg.SourceDir, "go-node-tmp_"+randomID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}

	workDir, _ = filepath.Abs(workDir)
	tmpDirAbs, _ := filepath.Abs(tmpDir)
	fmt.Printf("工作目录: %s\n", workDir)
	fmt.Printf("临时目录: %s\n\n", tmpDirAbs)

	if err := cleanupBuild(workDir); err != nil {
		return fmt.Errorf("清理build目录失败: %w", err)
	}

	fmt.Println("步骤1/4: 编译Go代码为动态库...")
	if err := createBuildDirectory(workDir); err != nil {
		return fmt.Errorf("创建build目录失败: %w", err)
	}
	if err := buildGoSharedLibrary(cfg, workDir); err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("步骤2/4: 解析Go导出函数...")
	functions, err := parseGoExports(absInputFile)
	if err != nil {
		return fmt.Errorf("解析Go导出函数失败: %w", err)
	}
	if len(functions) == 0 {
		return fmt.Errorf("未找到任何导出函数（使用 //export 标记）")
	}
	fmt.Printf("找到导出函数: %v\n\n", functions)

	fmt.Println("步骤3/4: 生成绑定文件...")
	dllPath := filepath.Join(workDir, "build", cfg.ModuleName+".dll")
	if err := generateResourceRC(tmpDir, cfg.ModuleName, dllPath); err != nil {
		return err
	}
	if err := generateBindingGyp(tmpDir, cfg.ModuleName); err != nil {
		return err
	}
	if err := generateWrapperCC(tmpDir, cfg.ModuleName, functions); err != nil {
		return err
	}
	fmt.Println()

	fmt.Println("步骤4/4: 使用node-gyp编译...")
	if err := runNodeGyp(tmpDir); err != nil {
		return err
	}

	if err := buildWithNodeGyp(tmpDir, outputNodeDir); err != nil {
		return err
	}
	fmt.Println()

	if !cfg.NoCleanup {
		fmt.Println("清理临时文件...")
		if err := cleanupBuild(workDir); err != nil {
			fmt.Printf("警告: 清理build目录失败: %v\n", err)
		}
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Printf("警告: 清理临时目录失败: %v\n", err)
		} else {
			fmt.Printf("已删除临时目录: %s\n", tmpDir)
		}
	}

	fmt.Println("========================")
	fmt.Printf("编译完成!\n")
	fmt.Printf("输出文件:\n")
	fmt.Printf("  - %s\n", filepath.Join(outputNodeDir, cfg.ModuleName+".node"))
	fmt.Println("========================")

	return nil
}

func init() {
	strings.Contains("", "")
}
