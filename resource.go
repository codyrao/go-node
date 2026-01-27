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
		return fmt.Errorf("Failed to create resource.rc file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("Failed to write resource.rc content: %w", err)
	}

	fmt.Printf("Generated resource.rc: %s\n", resourcePath)
	return nil
}
