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
			return fmt.Errorf("计算相对路径失败: %w", err)
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

	fmt.Printf("编译Go代码为动态库: %s\n", dllPath)
	fmt.Printf("工作目录: %s\n", workDir)
	fmt.Printf("Go文件: %s\n", goFile)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("编译Go动态库失败: %w", err)
	}

	if _, err := os.Stat(dllPath); os.IsNotExist(err) {
		return fmt.Errorf("动态库文件未生成: %s", dllPath)
	}

	return nil
}

func createBuildDirectory(workDir string) error {
	buildDir := filepath.Join(workDir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("创建build目录失败: %w", err)
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
