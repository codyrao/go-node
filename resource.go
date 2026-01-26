package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generateResourceRC(workDir, moduleName string, dllPath string) error {
	resourcePath := filepath.Join(workDir, moduleName+".rc")

	content := `#include <winuser.h>

` + moduleName + `_DLL RCDATA "` + strings.ReplaceAll(dllPath, "\\", "\\\\") + `"
`

	file, err := os.Create(resourcePath)
	if err != nil {
		return fmt.Errorf("创建resource.rc文件失败: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("写入resource.rc内容失败: %w", err)
	}

	fmt.Printf("生成resource.rc: %s\n", resourcePath)
	return nil
}
