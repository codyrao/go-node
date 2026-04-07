package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateWrapperCCUsesJSONParseCompatibilityHelper ensures generated wrapper code avoids direct JSON::Parse symbol dependencies.
func TestGenerateWrapperCCUsesJSONParseCompatibilityHelper(t *testing.T) {
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

	// Assert the generated wrapper uses a V8-version-safe helper instead of linking directly to JSON::Parse.
	if !strings.Contains(content, "Local<Value> ParseJSONString(") {
		t.Fatalf("expected generated wrapper to define ParseJSONString")
	}
	if !strings.Contains(content, "parseFunction->Call(context, jsonObject, 1, argv)") {
		t.Fatalf("expected generated wrapper to invoke JSON.parse through the global JSON object")
	}
	if strings.Contains(content, "JSON::Parse(") {
		t.Fatalf("expected generated wrapper to avoid direct JSON::Parse calls")
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

// TestGenerateWrapperCCResetsDLLStateOnReload verifies generated wrapper code clears cached DLL state before unload and reload operations.
func TestGenerateWrapperCCResetsDLLStateOnReload(t *testing.T) {
	// Create an isolated workspace so the generated wrapper can be inspected deterministically.
	workDir := t.TempDir()

	// Generate wrapper code for one export and load the file content for assertions.
	if err := generateWrapperCC(workDir, "hello", []string{"Hello1"}, ""); err != nil {
		t.Fatalf("generateWrapperCC returned error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(workDir, "wrapper.cc"))
	if err != nil {
		t.Fatalf("failed to read generated wrapper.cc: %v", err)
	}
	content := string(data)

	// Require an explicit DLL state reset helper and per-function pointer clearing to avoid stale addresses.
	if !strings.Contains(content, "void ResetDLLState()") {
		t.Fatalf("expected generated wrapper to define ResetDLLState")
	}
	if !strings.Contains(content, "Hello1Ptr = NULL;") {
		t.Fatalf("expected generated wrapper to clear cached export function pointers")
	}
	if !strings.Contains(content, "ResetDLLState();") {
		t.Fatalf("expected generated wrapper to call ResetDLLState during DLL lifecycle transitions")
	}
}

// TestGenerateWrapperCCBlocksUnsafeUnload verifies generated wrapper code rejects unload and reload while wrapper work is still active.
func TestGenerateWrapperCCBlocksUnsafeUnload(t *testing.T) {
	// Create an isolated workspace so the generated wrapper can be inspected deterministically.
	workDir := t.TempDir()

	// Generate wrapper code for one export and load the file content for assertions.
	if err := generateWrapperCC(workDir, "hello", []string{"Hello1"}, ""); err != nil {
		t.Fatalf("generateWrapperCC returned error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(workDir, "wrapper.cc"))
	if err != nil {
		t.Fatalf("failed to read generated wrapper.cc: %v", err)
	}
	content := string(data)

	// Require active-work tracking so wrapper cannot unload a DLL while Go code or callbacks still depend on it.
	if !strings.Contains(content, "bool HasActiveWrapperWork()") {
		t.Fatalf("expected generated wrapper to define HasActiveWrapperWork")
	}
	if !strings.Contains(content, "activeExportCalls") {
		t.Fatalf("expected generated wrapper to track in-flight exported function calls")
	}
	if !strings.Contains(content, "legacyCallbackSafetyLock") {
		t.Fatalf("expected generated wrapper to preserve unload safety for legacy callback registrations")
	}
	if !strings.Contains(content, "supportsManagedCallbacks") {
		t.Fatalf("expected generated wrapper to distinguish managed callback support from legacy callback support")
	}
	if !strings.Contains(content, "Cannot unload Go library while callbacks or exported calls are still active") {
		t.Fatalf("expected generated wrapper to reject unsafe DLL unload")
	}
	if !strings.Contains(content, "Cannot replace Go library while callbacks or exported calls are still active") {
		t.Fatalf("expected generated wrapper to reject unsafe DLL reload")
	}
}

// TestGenerateWrapperCCHandlesCallbackExceptionsSafely verifies generated wrapper code avoids fatal V8 unchecked callback invocation.
func TestGenerateWrapperCCHandlesCallbackExceptionsSafely(t *testing.T) {
	// Create an isolated workspace so the generated wrapper can be inspected deterministically.
	workDir := t.TempDir()

	// Generate wrapper code for one export and load the file content for assertions.
	if err := generateWrapperCC(workDir, "hello", []string{"Hello1"}, ""); err != nil {
		t.Fatalf("generateWrapperCC returned error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(workDir, "wrapper.cc"))
	if err != nil {
		t.Fatalf("failed to read generated wrapper.cc: %v", err)
	}
	content := string(data)

	// Require safe callback execution that preserves JS exceptions instead of terminating the host process.
	if !strings.Contains(content, "TryCatch tryCatch(isolate);") {
		t.Fatalf("expected generated wrapper to guard callback invocation with TryCatch")
	}
	if !strings.Contains(content, "MaybeLocal<Value> callbackResult = callback->Call(context, Null(isolate), 1, argv);") {
		t.Fatalf("expected generated wrapper to capture callback call results without ToLocalChecked")
	}
	if strings.Contains(content, "callback->Call(context, Null(isolate), 1, argv).ToLocalChecked();") {
		t.Fatalf("expected generated wrapper to avoid ToLocalChecked on callback invocation")
	}
}
