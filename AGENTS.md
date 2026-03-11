# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go CLI that builds Node.js native modules from Go source. Core implementation files live at the repository root: `main.go` handles flags and orchestration, `builder.go` compiles the shared library, `wrapper.go` and `callback.h` provide the bridge layer, and `nodegyp.go` / `bindgyp.go` generate and build the Node addon. Documentation is in `README.md` and `README_CN.md`. Use `test/` for sample Go inputs, JS verification scripts, and generated artifacts under `test/output/`.

## Build, Test, and Development Commands
Use Go 1.23.x and Node.js 16+ on Windows.

```bash
go build -o go-node.exe .
```
Builds the CLI from the current source tree.

```bash
go run . -input=test/hello.go -name=hello -output=test/output
```
Compiles the sample Go file into `test/output/hello.node` and `hello.dll`.

```bash
node test/test.js
```
Runs the main Node-side verification script against the generated addon.

```bash
node test/test_callback_fix.js
npx electron@v39.0.0 electron-test.js
```
Checks callback behavior and Electron compatibility.

## Coding Style & Naming Conventions
Format all Go files with `gofmt -w .` before submitting. Follow Go defaults: tabs for indentation, exported identifiers in `PascalCase`, internal helpers in `camelCase`, and concise error messages wrapped with `%w` where appropriate. Keep generated filenames and module names lowercase, for example `hello.node`. JavaScript test scripts use simple CommonJS files named `test_*.js` or `*-test.js`.

## Testing Guidelines
There is no `go test` suite yet; validation is example-driven. Add or update fixtures in `test/` when changing bridge behavior, callbacks, DLL loading, or Electron support. Prefer small, reproducible samples such as `test/hello.go`. A change is not ready until the relevant `go run . ...` build succeeds and the matching Node or Electron script runs cleanly.

## Commit & Pull Request Guidelines
Recent history uses short, focused subjects such as `[fea] update readme.md` and `Update .gitattributes to exclude node_modules/test`. Keep commits scoped to one change, start with a compact prefix like `fea`, `fix`, or `docs`, and use either English or Chinese consistently within the message. PRs should include the purpose, the exact commands used for verification, linked issues if any, and screenshots only when UI or Electron behavior changes.
