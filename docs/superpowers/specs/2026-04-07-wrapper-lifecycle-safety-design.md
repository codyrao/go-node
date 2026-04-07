# Wrapper Lifecycle Safety Design

**Context**

The generated Node wrapper currently has three P0 failure modes:
stale cached function pointers survive DLL unload/reload, DLL unload can happen while Go code is still executing asynchronously, and JS callback exceptions can crash the host process through `ToLocalChecked()`.

**Design**

Use a minimal containment fix inside generated `wrapper.cc`.

1. Centralize DLL runtime state reset.
Reset `hDLL`, `freeCStringPtr`, and every cached exported function pointer whenever the wrapper unloads a DLL or starts a fresh load path. This removes stale pointers before any later call site can reuse freed code addresses.

2. Block unsafe unload and reload.
Track in-flight exported function calls and queued callback deliveries in generated wrapper state. If a DLL unload or replacement load is requested while work is still active, reject the operation with a JS exception instead of calling `FreeLibrary` immediately.

3. Make async callback delivery exception-safe.
Wrap callback invocation in `TryCatch` and stop using `ToLocalChecked()` on the callback result. JS exceptions must stay inside V8 exception handling and must not terminate the Node/Electron host process.

**Testing**

Regression coverage will live in `wrapper_test.go` by asserting the generated `wrapper.cc` contains the new lifecycle guards and safe callback invocation paths.
