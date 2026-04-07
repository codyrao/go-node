package main

import (
	"os"
	"path/filepath"
	"strings"
)

// generateCallbackHeader writes callback bridge declarations used during cgo compilation.
func generateCallbackHeader(workDir string) error {
	// Keep header generation inside the source directory so local includes keep working.
	headerPath := filepath.Join(workDir, "callback.h")

	content := `/* Auto-generated callback header for go-node */
#ifndef GO_NODE_CALLBACK_H
#define GO_NODE_CALLBACK_H

#include <stdlib.h>
#include <stdint.h>
#include <string.h>

#ifdef _WIN32
#include <windows.h>
#endif

typedef void (*CallbackFunc)(const char*);
typedef void (*CallbackFuncWithName)(const char*, const char*);
typedef void (*CallbackFuncWithId)(int32_t, const char*);
typedef void (*CallbackControlFunc)(int32_t);

static void callCallback(void* ptr, const char* data) {
    if (ptr != NULL) {
        ((CallbackFunc)ptr)(data);
    }
}

static void callCallbackWithFuncName(void* ptr, const char* funcName, const char* data) {
    if (ptr != NULL) {
        ((CallbackFuncWithName)ptr)(funcName, data);
    }
}

static void callCallbackWithId(void* ptr, int32_t callbackId, const char* data) {
    if (ptr != NULL) {
        ((CallbackFuncWithId)ptr)(callbackId, data);
    }
}

static void keepCallback(void* ptr, int32_t callbackId) {
    if (ptr != NULL) {
        ((CallbackControlFunc)ptr)(callbackId);
    }
}

static void freeCallback(void* ptr, int32_t callbackId) {
    if (ptr != NULL) {
        ((CallbackControlFunc)ptr)(callbackId);
    }
}

#endif /* GO_NODE_CALLBACK_H */
`

	file, err := os.Create(headerPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Persist the generated callback declarations for this build run.
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

// generateWrapperCC creates the C++ Node addon wrapper used to call into exported Go DLL symbols.
func generateWrapperCC(workDir, moduleName string, functions []string, hexArray string) error {
	// Resolve one deterministic output file path for generated wrapper source.
	wrapperPath := filepath.Join(workDir, "wrapper.cc")

	funcDecls := ""
	funcImpls := ""
	nodeMethods := ""

	// Emit one wrapper function per export while skipping internal bridge symbols.
	for _, fn := range functions {
		funcName := fn
		jsName := toCamelCase(fn)

		if strings.Contains(fn, "Register") || strings.Contains(fn, "Free") {
			continue
		}

		funcDecls += "typedef const char* (*" + funcName + "Fn)(const char*, const char*);\n" +
			funcName + "Fn " + funcName + "Ptr = NULL;\n\n"

		funcImpls += "void " + funcName + "Wrapper(const FunctionCallbackInfo<Value>& args) {\n" +
			"    Isolate* isolate = args.GetIsolate();\n" +
			"    Local<Context> context = isolate->GetCurrentContext();\n" +
			"    \n" +
			"    if (hDLL == NULL) {\n" +
			"        isolate->ThrowException(Exception::Error(\n" +
			"            String::NewFromUtf8(isolate, \"Go library not loaded\").ToLocalChecked()));\n" +
			"        return;\n" +
			"    }\n" +
			"    \n" +
			"    if (" + funcName + "Ptr == NULL) {\n" +
			"        " + funcName + "Ptr = (" + funcName + "Fn)GetProcAddress(hDLL, \"" + funcName + "\");\n" +
			"        if (" + funcName + "Ptr == NULL) {\n" +
			"            isolate->ThrowException(Exception::Error(\n" +
			"                String::NewFromUtf8(isolate, \"Failed to find function " + funcName + "\").ToLocalChecked()));\n" +
			"            return;\n" +
			"        }\n" +
			"    }\n" +
			"    \n" +
			"    const char* arg1 = \"\";\n" +
			"    bool freeArg1 = false;\n" +
			"    std::string callbackToken;\n" +
			"    int32_t callbackId = -1;\n" +
			"    const char* arg2 = \"\";\n" +
			"    \n" +
			"    if (args.Length() > 0) {\n" +
			"        if (args[0]->IsString()) {\n" +
			"            String::Utf8Value utf8(isolate, args[0]);\n" +
			"            arg1 = strdup(*utf8);\n" +
			"            freeArg1 = true;\n" +
			"        } else if (args[0]->IsObject() && !args[0]->IsArray()) {\n" +
			"            Local<Object> obj = Local<Object>::Cast(args[0]);\n" +
			"            Local<String> jsonStr;\n" +
			"            if (v8::JSON::Stringify(context, obj).ToLocal(&jsonStr)) {\n" +
			"                String::Utf8Value jsonUtf8(isolate, jsonStr);\n" +
			"                arg1 = strdup(*jsonUtf8);\n" +
			"                freeArg1 = true;\n" +
			"            }\n" +
			"        }\n" +
			"    }\n" +
			"    \n" +
			"    if (args.Length() > 1 && args[1]->IsFunction()) {\n" +
			"        callbackId = RegisterCallback(isolate, args[1].As<Function>());\n" +
			"        callbackToken = std::to_string(callbackId);\n" +
			"        arg2 = callbackToken.c_str();\n" +
			"    }\n" +
			"    \n" +
			"    const char* result = " + funcName + "Ptr(arg1, arg2);\n" +
			"    if (callbackId >= 0) {\n" +
			"        ReleaseCallback(callbackId);\n" +
			"    }\n" +
			"    if (freeArg1) {\n" +
			"        free((void*)arg1);\n" +
			"    }\n" +
			"    \n" +
			"    if (result != NULL) {\n" +
			"        Local<Value> parsedResult = ParseJSONResult(isolate, result);\n" +
			"        args.GetReturnValue().Set(parsedResult);\n" +
			"    }\n" +
			"    if (freeCStringPtr != NULL && result != NULL) {\n" +
			"        freeCStringPtr(result);\n" +
			"    }\n" +
			"}\n\n"

		nodeMethods += "    NODE_SET_METHOD(exports, \"" + jsName + "\", " + funcName + "Wrapper);\n"
	}

	// Render C++ template content after all function-specific fragments are prepared.
	content := `#include <node.h>
#include <v8.h>
#include <uv.h>
#include <windows.h>
#include <map>
#include <string>
#include <stdlib.h>
#include <queue>
#include <mutex>
#include <fstream>

#pragma warning(disable: 4018)
#pragma warning(disable: 4996)

using namespace v8;

#define GO_NODE_DLL_RESOURCE_ID 1

// Extract embedded DLL from resources
std::string ExtractEmbeddedDLL() {
    char tempPath[MAX_PATH];
    DWORD result = GetTempPathA(MAX_PATH, tempPath);
    if (result == 0 || result > MAX_PATH) {
        return "";
    }

    char tempFilePath[MAX_PATH];
    if (GetTempFileNameA(tempPath, "gnd", 0, tempFilePath) == 0) {
        return "";
    }

    DeleteFileA(tempFilePath);

    std::string dllPath = std::string(tempFilePath) + ".dll";
    
    // First, try to find DLL in the same directory as the .node file
    char modulePath[MAX_PATH];
    HMODULE hModule = NULL;
    if (GetModuleHandleExA(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS | GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
                            (LPCSTR)&ExtractEmbeddedDLL, &hModule)) {
        GetModuleFileNameA(hModule, modulePath, MAX_PATH);
        
        // Get directory path
        char* lastSlash = strrchr(modulePath, '\\');
        if (lastSlash != NULL) {
            *lastSlash = '\0';
            std::string localDllPath = std::string(modulePath) + "\\" + "` + moduleName + `.dll";
            
            // Check if DLL exists in local directory
            std::ifstream testFile(localDllPath, std::ios::binary);
            if (testFile.good()) {
                testFile.close();
                return localDllPath;
            }
        }
    }
    
    // Try to extract from resources
    if (hModule != NULL) {
        HRSRC hRes = FindResourceA(hModule, MAKEINTRESOURCEA(GO_NODE_DLL_RESOURCE_ID), RT_RCDATA);
        if (hRes != NULL) {
            HGLOBAL hLoaded = LoadResource(hModule, hRes);
            if (hLoaded != NULL) {
                void* pData = LockResource(hLoaded);
                DWORD size = SizeofResource(hModule, hRes);
                if (pData != NULL && size > 0) {
                    std::ofstream outFile(dllPath, std::ios::binary);
                    if (outFile.is_open()) {
                        outFile.write(reinterpret_cast<const char*>(pData), size);
                        outFile.close();
                        return dllPath;
                    }
                }
            }
        }
    }
    
    return "";
}

// Clean up temporary DLL
void CleanupEmbeddedDLL(const std::string& dllPath) {
    if (!dllPath.empty()) {
        DeleteFileA(dllPath.c_str());
    }
}

// Parse JSON result and return object type
Local<Value> ParseJSONResult(Isolate* isolate, const char* jsonStr) {
    if (jsonStr == NULL || strlen(jsonStr) == 0) {
        return Null(isolate);
    }

    // Decode UTF-8 bytes into a V8 string before parsing JSON.
    Local<String> jsonString;
    if (!String::NewFromUtf8(isolate, jsonStr, NewStringType::kNormal, strlen(jsonStr)).ToLocal(&jsonString)) {
        return Null(isolate);
    }
    Local<Context> context = isolate->GetCurrentContext();

    // Parse using JSON::Parse to avoid script compilation overhead.
    Local<Value> parsed;
    if (!JSON::Parse(context, jsonString).ToLocal(&parsed)) {
        return Null(isolate);
    }
    
    // Only return object type
    if (parsed->IsObject()) {
        return parsed;
    }
    
    return Null(isolate);
}

` + funcDecls + `struct CallbackEntry {
    Global<Function> callback;
    bool closed;
    bool persistent;
    uint32_t pendingCount;
    uint32_t activeCalls;
};

// Global variables
HMODULE hDLL = NULL;
std::map<int32_t, CallbackEntry> callbackMap;
int32_t nextCallbackId = 0;
Isolate* gIsolate = NULL;
std::mutex callbackMutex;

// Async callback handling
struct CallbackData {
    int32_t callbackId;
    std::string jsonData;
};

std::queue<CallbackData> callbackQueue;
std::mutex queueMutex;
uv_async_t async_handle;

void AsyncCallback(uv_async_t* handle) {
    Isolate* isolate = gIsolate;
    if (isolate == NULL) return;
    
    HandleScope handleScope(isolate);
    Local<Context> context = isolate->GetCurrentContext();
    Context::Scope contextScope(context);
    
    while (true) {
        CallbackData data;
        {
            std::lock_guard<std::mutex> lock(queueMutex);
            if (callbackQueue.empty()) {
                break;
            }

            data = callbackQueue.front();
            callbackQueue.pop();
        }

        Local<Function> callback;
        {
            std::lock_guard<std::mutex> lock(callbackMutex);
            auto it = callbackMap.find(data.callbackId);
            if (it == callbackMap.end()) {
                continue;
            }

            callback = Local<Function>::New(isolate, it->second.callback);
        }

        if (callback.IsEmpty()) {
            continue;
        }
        
        Local<String> jsonString;
        Local<Value> parsedData;

        if (String::NewFromUtf8(isolate, data.jsonData.c_str(), NewStringType::kNormal, data.jsonData.length()).ToLocal(&jsonString) &&
            JSON::Parse(context, jsonString).ToLocal(&parsedData) &&
            parsedData->IsObject()) {
            Local<Value> argv[] = { parsedData };
            callback->Call(context, Null(isolate), 1, argv).ToLocalChecked();
        } else {
            Local<Object> emptyObj = Object::New(isolate);
            Local<Value> argv[] = { emptyObj };
            callback->Call(context, Null(isolate), 1, argv).ToLocalChecked();
        }

        {
            std::lock_guard<std::mutex> lock(callbackMutex);
            auto it = callbackMap.find(data.callbackId);
            if (it == callbackMap.end()) {
                continue;
            }

            if (it->second.pendingCount > 0) {
                it->second.pendingCount--;
            }

            if (it->second.closed && it->second.pendingCount == 0 && it->second.activeCalls == 0) {
                it->second.callback.Reset();
                callbackMap.erase(it);
            }
        }
    }
}

void InitializeAsyncHandle() {
    uv_loop_t* loop = uv_default_loop();
    uv_async_init(loop, &async_handle, AsyncCallback);
    async_handle.data = NULL;
}

// Callback registration and invocation
int32_t RegisterCallback(Isolate* isolate, Local<Function> callback) {
    std::lock_guard<std::mutex> lock(callbackMutex);

    int32_t id = nextCallbackId++;
    CallbackEntry entry;
    entry.callback = Global<Function>(isolate, callback);
    entry.closed = false;
    entry.persistent = false;
    entry.pendingCount = 0;
    entry.activeCalls = 1;
    callbackMap[id] = std::move(entry);

    return id;
}

typedef void (*CallNodeCallbackType)(int32_t, const char*);
typedef void (*CallbackControlType)(int32_t);
typedef void (*FreeCStringFn)(const char*);
FreeCStringFn freeCStringPtr = NULL;

void CallNodeCallback(int32_t callbackId, const char* jsonData);
void KeepCallback(int32_t callbackId);
void CloseCallback(int32_t callbackId);

bool BindGoRuntimeFunctions() {
    freeCStringPtr = (FreeCStringFn)GetProcAddress(hDLL, "FreeCString");

    typedef void (*RegisterGoCallbackType)(CallNodeCallbackType);
    typedef void (*RegisterGoCallbackExType)(CallNodeCallbackType, CallbackControlType, CallbackControlType);
    RegisterGoCallbackExType registerGoCallbackEx = (RegisterGoCallbackExType)GetProcAddress(hDLL, "RegisterGoCallbackEx");
    if (registerGoCallbackEx != NULL) {
        registerGoCallbackEx(CallNodeCallback, KeepCallback, CloseCallback);
        return true;
    }

    RegisterGoCallbackType registerGoCallback = (RegisterGoCallbackType)GetProcAddress(hDLL, "RegisterGoCallback");
    if (registerGoCallback != NULL) {
        registerGoCallback(CallNodeCallback);
    }

    return true;
}

void CallNodeCallback(int32_t callbackId, const char* jsonData) {
    {
        std::lock_guard<std::mutex> lock(callbackMutex);
        auto it = callbackMap.find(callbackId);
        if (it == callbackMap.end() || it->second.closed) {
            return;
        }

        it->second.pendingCount++;
    }

    CallbackData data;
    data.callbackId = callbackId;
    data.jsonData = jsonData;
    
    {
        std::lock_guard<std::mutex> lock(queueMutex);
        callbackQueue.push(data);
    }
    
    uv_async_send(&async_handle);
}

void KeepCallback(int32_t callbackId) {
    std::lock_guard<std::mutex> lock(callbackMutex);

    auto it = callbackMap.find(callbackId);
    if (it == callbackMap.end() || it->second.closed) {
        return;
    }

    it->second.persistent = true;
}

void ReleaseCallback(int32_t callbackId) {
    std::lock_guard<std::mutex> lock(callbackMutex);

    auto it = callbackMap.find(callbackId);
    if (it == callbackMap.end()) {
        return;
    }

    if (it->second.activeCalls > 0) {
        it->second.activeCalls--;
    }

    if (!it->second.persistent) {
        it->second.closed = true;
    }

    if (it->second.closed && it->second.pendingCount == 0 && it->second.activeCalls == 0) {
        it->second.callback.Reset();
        callbackMap.erase(it);
    }
}

void CloseCallback(int32_t callbackId) {
    std::lock_guard<std::mutex> lock(callbackMutex);

    auto it = callbackMap.find(callbackId);
    if (it == callbackMap.end()) {
        return;
    }

    it->second.closed = true;
    it->second.persistent = false;
    if (it->second.pendingCount == 0 && it->second.activeCalls == 0) {
        it->second.callback.Reset();
        callbackMap.erase(it);
    }
}

// Load DLL functions
void LoadGoLibrary(const FunctionCallbackInfo<Value>& args) {
    Isolate* isolate = args.GetIsolate();
    gIsolate = isolate;
    
    if (args.Length() < 1) {
        isolate->ThrowException(Exception::Error(
            String::NewFromUtf8(isolate, "DLL path required").ToLocalChecked()));
        return;
    }
    
    String::Utf8Value dllPath(isolate, args[0]);
    
    if (hDLL != NULL) {
        FreeLibrary(hDLL);
    }
    hDLL = LoadLibraryA(*dllPath);
    if (hDLL == NULL) {
        isolate->ThrowException(Exception::Error(
            String::NewFromUtf8(isolate, "Failed to load DLL").ToLocalChecked()));
        return;
    }

    BindGoRuntimeFunctions();
    
    args.GetReturnValue().Set(True(isolate));
}

void UnloadGoLibrary(const FunctionCallbackInfo<Value>& args) {
    Isolate* isolate = args.GetIsolate();
    
    std::lock_guard<std::mutex> lock(callbackMutex);

    for (auto& pair : callbackMap) {
        pair.second.callback.Reset();
    }
    callbackMap.clear();
    
    if (hDLL != NULL) {
        FreeLibrary(hDLL);
        hDLL = NULL;
    }
    
    args.GetReturnValue().Set(True(isolate));
}

// Exported callback functions for Go
extern "C" {
    __declspec(dllexport) void CallCallback(int32_t callbackId, const char* jsonData) {
        CallNodeCallback(callbackId, jsonData);
    }

    __declspec(dllexport) void FreeCallback(int32_t callbackId) {
        CloseCallback(callbackId);
    }
}

` + funcImpls + `// Forward declaration
void InitializeWrapper(const FunctionCallbackInfo<Value>& args);

void Initialize(Local<Object> exports) {
    Isolate* isolate = Isolate::GetCurrent();
    gIsolate = isolate;
    
    // Initialize async handle
    InitializeAsyncHandle();
    
    // Extract embedded DLL to temporary directory
    std::string dllPath = ExtractEmbeddedDLL();
    if (dllPath.empty()) {
        isolate->ThrowException(Exception::Error(
            String::NewFromUtf8(isolate, "Failed to extract embedded DLL").ToLocalChecked()));
        return;
    }
    
    hDLL = LoadLibraryA(dllPath.c_str());
    if (hDLL == NULL) {
        CleanupEmbeddedDLL(dllPath);
        isolate->ThrowException(Exception::Error(
            String::NewFromUtf8(isolate, ("Failed to load " + dllPath).c_str()).ToLocalChecked()));
        return;
    }
    
    BindGoRuntimeFunctions();
    
    NODE_SET_METHOD(exports, "init", InitializeWrapper);
    NODE_SET_METHOD(exports, "loadLibrary", LoadGoLibrary);
    NODE_SET_METHOD(exports, "unloadLibrary", UnloadGoLibrary);
` + nodeMethods + `}

void InitializeWrapper(const FunctionCallbackInfo<Value>& args) {
    Isolate* isolate = args.GetIsolate();
    gIsolate = isolate;
    
    if (args.Length() > 0) {
        String::Utf8Value dllPath(isolate, args[0]);
        
        if (hDLL != NULL) {
            FreeLibrary(hDLL);
        }
        hDLL = LoadLibraryA(*dllPath);
        if (hDLL == NULL) {
            isolate->ThrowException(Exception::Error(
                String::NewFromUtf8(isolate, "Failed to load Go library").ToLocalChecked()));
            return;
        }
        
        BindGoRuntimeFunctions();
    }
    
    args.GetReturnValue().Set(True(isolate));
}

NODE_MODULE(` + moduleName + `, Initialize)`

	file, err := os.Create(wrapperPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

// toCamelCase converts exported Go function names to JavaScript-style camelCase names.
func toCamelCase(s string) string {
	// Return early for empty identifiers.
	if len(s) == 0 {
		return s
	}

	// Split by underscore to support snake_case export names.
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return strings.ToLower(s[:1]) + s[1:]
	}

	// Lowercase the first part and capitalize the following parts to build camelCase.
	for i := 0; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			if i == 0 {
				parts[i] = strings.ToLower(parts[i][:1]) + parts[i][1:]
			} else {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
	}
	return strings.Join(parts, "")
}

// parseGoExports scans source code and returns names annotated with //export markers.
func parseGoExports(goFile string) ([]string, error) {
	// Read full source content so export scanning can run in one pass.
	content, err := os.ReadFile(goFile)
	if err != nil {
		return nil, err
	}

	var functions []string
	lines := strings.Split(string(content), "\n")

	// Filter out bridge-internal exports and keep user-facing symbols only.
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "//export") {
			fnName := strings.TrimSpace(strings.TrimPrefix(line, "//export"))
			if fnName != "" && fnName != "CallCallback" && fnName != "FreeCallback" && fnName != "RegisterGoCallback" && fnName != "incrementCounter" {
				functions = append(functions, fnName)
			}
		}
	}

	return functions, nil
}
