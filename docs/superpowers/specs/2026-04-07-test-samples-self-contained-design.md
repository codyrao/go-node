# Test Samples Self-Contained Design

**Context**

The `test/` directory is currently a collection of standalone sample source files that depend on generated side effects.
That makes `go test ./...` fail for two separate reasons:

1. `test/hello.go` and `test/test_callback_fix.go` both declare the same callback runtime symbols such as `RegisterGoCallback` and `main`.
2. The package depends on a generated `callback.h`, so a clean checkout cannot compile the package directly.

The user wants a long-term self-contained layout while keeping existing CLI compatibility with commands such as `-input=test/hello.go`.

**Design**

Use a package-based sample layout inside the existing `test/` directory while preserving file paths.

1. Add one shared sample runtime file in `test/`.
It will own the cgo preamble, callback bridge state, the single exported `RegisterGoCallback`, reusable helpers, and the single empty `main`.

2. Keep `test/hello.go` and `test/test_callback_fix.go` as sample entry files.
They will keep their exported sample functions, but remove duplicated callback runtime declarations so the package can compile cleanly as one package.

3. Check in `test/callback.h` as a stable sample dependency.
The `test` package must be buildable from a clean checkout without first running the CLI.

4. Stop treating every `callback.h` as disposable generated output.
The builder should only create `callback.h` when it is missing and should only delete the header when this build run created it. Existing checked-in headers must be preserved.

**Testing**

Add regression coverage in `builder_test.go` for preserving pre-existing `callback.h`.
Then verify with `go test ./...`, plus existing sample build and Node verification commands.
