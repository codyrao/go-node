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
		return fmt.Errorf("Release目录不存在: %s", releaseDir)
	}
	
	files, err := os.ReadDir(releaseDir)
	if err != nil {
		return fmt.Errorf("读取Release目录失败: %w", err)
	}
	
	var nodeFileName string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".node") {
			nodeFileName = f.Name()
			break
		}
	}
	
	if nodeFileName == "" {
		return fmt.Errorf("未找到.node文件")
	}
	
	srcNode := filepath.Join(releaseDir, nodeFileName)
	dstNode := filepath.Join(outputNodeDir, nodeFileName)
	
	if err := copyFile(srcNode, dstNode); err != nil {
		return fmt.Errorf("复制.node文件失败: %w", err)
	}
	
	fmt.Printf("复制.node文件到: %s\n", dstNode)
	return nil
}

func runNodeGyp(workDir string) error {
	nodeGypPath := "C:\\Users\\admin\\AppData\\Roaming\\npm\\node-gyp.cmd"

	args := []string{
		"configure",
		"--builddir=build",
	}

	cmd := exec.Command(nodeGypPath, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("运行node-gyp configure...\n")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("node-gyp configure失败: %w", err)
	}

	args = []string{
		"build",
	}

	cmd = exec.Command(nodeGypPath, args...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("运行node-gyp build...\n")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("node-gyp build失败: %w", err)
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
		return fmt.Errorf("复制DLL文件失败: %w", err)
	}
	
	fmt.Printf("复制DLL文件到: %s\n", dstDLL)
	return nil
}
