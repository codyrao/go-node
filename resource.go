package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func generateResourceRC(workDir, moduleName string, dllPath string) error {
	resourcePath := filepath.Join(workDir, moduleName+".rc")

	// Use only filename, not full path
	dllFileName := filepath.Base(dllPath)

	// Use string name instead of numeric ID
	content := `#include <winuser.h>

` + moduleName + `_DLL RCDATA "` + dllFileName + `"
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
	fmt.Printf("Resource path: %s\n", dllFileName)
	return nil
}
