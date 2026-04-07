package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateWrapperCCUsesJSONParse ensures generated wrapper code uses JSON::Parse instead of script compilation.
func TestGenerateWrapperCCUsesJSONParse(t *testing.T) {
	// Create an isolated workspace so generation side effects stay local to this test.
	workDir := t.TempDir()

	// Generate wrapper code for one exported function and then load the generated source.
	if err := generateWrapperCC(workDir, "hello", []string{"Hello1"}, ""); err != nil {
		t.Fatalf("generateWrapperCC returned error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(workDir, "wrapper.cc"))
	if err != nil {
		t.Fatalf("failed to read generated wrapper.cc: %v", err)
	}
	content := string(data)

	// Assert the optimized JSON path is present and legacy Script::Compile parsing is absent.
	if !strings.Contains(content, "JSON::Parse(context, jsonString)") {
		t.Fatalf("expected generated wrapper to contain JSON::Parse")
	}
	if strings.Contains(content, "Script::Compile(") {
		t.Fatalf("expected generated wrapper to avoid Script::Compile")
	}
}

// TestGenerateWrapperCCIncludesResultFreeing verifies generated wrapper code releases Go-owned C strings through DLL exports.
func TestGenerateWrapperCCIncludesResultFreeing(t *testing.T) {
	// Create a temp output directory for deterministic file inspection.
	workDir := t.TempDir()

	// Generate wrapper code and inspect free-function declarations and call sites.
	if err := generateWrapperCC(workDir, "hello", []string{"Hello1"}, ""); err != nil {
		t.Fatalf("generateWrapperCC returned error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(workDir, "wrapper.cc"))
	if err != nil {
		t.Fatalf("failed to read generated wrapper.cc: %v", err)
	}
	content := string(data)

	// Ensure both the function-pointer declaration and runtime free call exist.
	if !strings.Contains(content, "typedef void (*FreeCStringFn)(const char*);") {
		t.Fatalf("expected generated wrapper to declare FreeCStringFn")
	}
	if !strings.Contains(content, "freeCStringPtr = (FreeCStringFn)GetProcAddress(hDLL, \"FreeCString\")") {
		t.Fatalf("expected generated wrapper to resolve FreeCString symbol")
	}
	if !strings.Contains(content, "if (freeCStringPtr != NULL && result != NULL)") {
		t.Fatalf("expected generated wrapper to free Go result with FreeCString")
	}
}

// TestParseGoExportsSkipsInternalGlue verifies parseGoExports ignores internal bridge functions.
func TestParseGoExportsSkipsInternalGlue(t *testing.T) {
	// Build a synthetic Go source with both user and internal export markers.
	workDir := t.TempDir()
	goFile := filepath.Join(workDir, "sample.go")
	content := `package main

import "C"

//export Hello1
func Hello1(params *C.char, callbackType *C.char) *C.char { return nil }

//export RegisterGoCallback
func RegisterGoCallback(fn uintptr) {}

//export FreeCallback
func FreeCallback(id int32) {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write sample Go file: %v", err)
	}

	// Parse exports and verify only user-facing exports survive filtering.
	functions, err := parseGoExports(goFile)
	if err != nil {
		t.Fatalf("parseGoExports returned error: %v", err)
	}
	if len(functions) != 1 || functions[0] != "Hello1" {
		t.Fatalf("unexpected exported function list: %#v", functions)
	}
}
