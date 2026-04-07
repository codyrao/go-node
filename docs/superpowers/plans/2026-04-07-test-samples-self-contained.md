# Test Samples Self-Contained Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the `test/` sample directory compile cleanly as a normal Go package while preserving existing CLI sample file paths.

**Architecture:** Move shared callback runtime code into one package-level sample support file, keep individual sample exports in their existing files, and make `callback.h` lifecycle management respect checked-in headers. The CLI continues to accept `-input=test/hello.go` and `-input=test/test_callback_fix.go`.

**Tech Stack:** Go, cgo, sample Node addon fixtures

---

### Task 1: Lock callback header lifecycle with tests

**Files:**
- Modify: `C:\Users\admin\Documents\goproject\xlasgame\go-node\builder_test.go`
- Test: `C:\Users\admin\Documents\goproject\xlasgame\go-node\builder_test.go`

- [ ] **Step 1: Write the failing test**

Add a regression test proving cleanup preserves a pre-existing `callback.h` while still deleting generated helper files and build output.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test . -run "TestCleanupBuild" -count=1`
Expected: FAIL because `cleanupBuild` currently deletes all `callback.h` files.

### Task 2: Make builder cleanup ownership-aware

**Files:**
- Modify: `C:\Users\admin\Documents\goproject\xlasgame\go-node\builder.go`
- Test: `C:\Users\admin\Documents\goproject\xlasgame\go-node\builder_test.go`

- [ ] **Step 1: Write the minimal implementation**

Track whether `generateCallbackHeader` created the header for the current build and only delete it during cleanup when this run created it.

- [ ] **Step 2: Run focused tests to verify they pass**

Run: `go test . -run "TestCleanupBuild|TestGenerateRuntimeHelperFile" -count=1`
Expected: PASS

### Task 3: Make the sample package self-contained

**Files:**
- Create: `C:\Users\admin\Documents\goproject\xlasgame\go-node\test\callback_support.go`
- Create: `C:\Users\admin\Documents\goproject\xlasgame\go-node\test\callback.h`
- Modify: `C:\Users\admin\Documents\goproject\xlasgame\go-node\test\hello.go`
- Modify: `C:\Users\admin\Documents\goproject\xlasgame\go-node\test\test_callback_fix.go`

- [ ] **Step 1: Write the minimal implementation**

Move shared runtime declarations into `callback_support.go`, leave only sample exports in the two sample files, and add function/key-flow comments required by the repository rules.

- [ ] **Step 2: Run package verification**

Run: `go test ./... -count=1`
Expected: PASS

### Task 4: Re-run sample integration verification

**Files:**
- No code changes

- [ ] **Step 1: Build the sample addons**

Run:
- `go build -o go-node.exe .`
- `.\go-node.exe --% -input=test/hello.go -name=hello -output=test/output`
- `.\go-node.exe --% -input=test/test_callback_fix.go -name=test_callback_fix -output=test/output`

Expected: both sample builds succeed.

- [ ] **Step 2: Run Node verification**

Run:
- `node test/test.js`
- `node test/test_callback_fix.js`

Expected: both scripts pass.
