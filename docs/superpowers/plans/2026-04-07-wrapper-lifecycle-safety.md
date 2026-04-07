# Wrapper Lifecycle Safety Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Eliminate the generated wrapper's stale function-pointer, unsafe DLL unload, and fatal callback exception crash paths.

**Architecture:** Keep the fix local to the generated `wrapper.cc` template in `wrapper.go`. Add explicit wrapper runtime state management, reject unload/reload while work is active, and replace fatal callback invocation helpers with safe V8 exception handling.

**Tech Stack:** Go, generated C++ Node addon wrapper, V8, libuv

---

### Task 1: Lock the new behavior with regression tests

**Files:**
- Modify: `C:\Users\admin\Documents\goproject\xlasgame\go-node\wrapper_test.go`
- Test: `C:\Users\admin\Documents\goproject\xlasgame\go-node\wrapper_test.go`

- [ ] **Step 1: Write the failing tests**

Add wrapper-generation assertions that require:
- cached exported function pointers to be reset during unload/reload paths
- unload/reload to reject when active calls or queued callbacks still exist
- async callback invocation to use `TryCatch` and avoid `ToLocalChecked()`

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test . -run "TestGenerateWrapperCC" -count=1`
Expected: FAIL because the generated wrapper does not yet include the new guard code.

### Task 2: Implement minimal wrapper lifecycle protection

**Files:**
- Modify: `C:\Users\admin\Documents\goproject\xlasgame\go-node\wrapper.go`
- Test: `C:\Users\admin\Documents\goproject\xlasgame\go-node\wrapper_test.go`

- [ ] **Step 1: Write the minimal implementation**

Update the generated wrapper template to:
- define helper functions for resetting DLL state and checking whether unload is safe
- increment/decrement active export-call counters around Go DLL entrypoints
- reject unload/reload while the wrapper still has active work
- invoke JS callbacks through `TryCatch` without `ToLocalChecked()`

- [ ] **Step 2: Run focused tests to verify they pass**

Run: `go test . -run "TestGenerateWrapperCC" -count=1`
Expected: PASS

### Task 3: Run repository verification

**Files:**
- No code changes

- [ ] **Step 1: Run broader verification**

Run: `go test ./... -count=1`
Expected: Root package passes. If the historical `test` package still fails due to duplicate exports, report that separately as a pre-existing verification gap unless fixed as part of this task.
