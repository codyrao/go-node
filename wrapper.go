package main

import (
	"os"
	"path/filepath"
	"strings"
)

func generateCallbackHeader(workDir string) error {
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

static void callCallback(void* ptr, const char* data) {
    if (ptr != NULL) {
        ((CallbackFunc)ptr)(data);
    }
}

#endif /* GO_NODE_CALLBACK_H */
`

	file, err := os.Create(headerPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

func generateWrapperCC(workDir, moduleName string, functions []string, hexArray string) error {
	wrapperPath := filepath.Join(workDir, "wrapper.cc")

	funcDecls := ""
	funcImpls := ""
	nodeMethods := ""

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
			"    const char* arg2 = \"\";\n" +
			"    \n" +
			"    if (args.Length() > 0 && args[0]->IsObject() && !args[0]->IsArray()) {\n" +
			"        Local<Object> obj = Local<Object>::Cast(args[0]);\n" +
			"        Local<Context> context = isolate->GetCurrentContext();\n" +
			"        Local<String> jsonStr = v8::JSON::Stringify(context, obj).ToLocalChecked();\n" +
			"        String::Utf8Value jsonUtf8(isolate, jsonStr);\n" +
			"        arg1 = strdup(*jsonUtf8);\n" +
			"    }\n" +
			"    \n" +
			"    if (args.Length() > 1 && args[1]->IsFunction()) {\n" +
			"        RegisterCallback(isolate, args[1].As<Function>());\n" +
			"        arg2 = \"callback\";\n" +
			"    }\n" +
			"    \n" +
			"    const char* result = " + funcName + "Ptr(arg1, arg2);\n" +
			"    \n" +
			"    if (result != NULL) {\n" +
			"        Local<Value> parsedResult = ParseJSONResult(isolate, result);\n" +
			"        args.GetReturnValue().Set(parsedResult);\n" +
			"    }\n" +
			"}\n\n"

		nodeMethods += "    NODE_SET_METHOD(exports, \"" + jsName + "\", " + funcName + "Wrapper);\n"
	}

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

using namespace v8;

// Extract embedded DLL from resources
std::string ExtractEmbeddedDLL() {
    char tempPath[MAX_PATH];
    DWORD result = GetTempPathA(MAX_PATH, tempPath);
    if (result == 0 || result > MAX_PATH) {
        return "";
    }
    
    std::string dllPath = std::string(tempPath) + "` + moduleName + `.dll";
    
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
        HRSRC hRes = FindResourceA(hModule, "` + moduleName + `_DLL", RT_RCDATA);
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

// Callback function type and global variable
typedef void (*CallNodeCallbackType)(const char*);
CallNodeCallbackType gCallNodeCallback = NULL;

// Parse JSON result and return object type
Local<Value> ParseJSONResult(Isolate* isolate, const char* jsonStr) {
    if (jsonStr == NULL || strlen(jsonStr) == 0) {
        return Null(isolate);
    }
    
    Local<String> jsonString = String::NewFromUtf8(isolate, jsonStr).ToLocalChecked();
    Local<Context> context = isolate->GetCurrentContext();
    
    // Wrap JSON string in parentheses to make it a valid JavaScript expression
    std::string wrappedJson = "(" + std::string(jsonStr) + ")";
    Local<String> wrappedJsonStr = String::NewFromUtf8(isolate, wrappedJson.c_str()).ToLocalChecked();
    
    TryCatch tryCatch(isolate);
    Local<Script> script;
    if (!Script::Compile(context, wrappedJsonStr).ToLocal(&script)) {
        return Null(isolate);
    }
    
    Local<Value> parsed;
    if (!script->Run(context).ToLocal(&parsed)) {
        return Null(isolate);
    }
    
    // Only return object type
    if (parsed->IsObject()) {
        return parsed;
    }
    
    return Null(isolate);
}

` + funcDecls + `// Global variables
HMODULE hDLL = NULL;
std::map<int32_t, Global<Function>> callbackMap;
int32_t nextCallbackId = 0;
Isolate* gIsolate = NULL;

// Async callback handling
struct CallbackData {
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
    
    std::lock_guard<std::mutex> lock(queueMutex);
    
    while (!callbackQueue.empty()) {
        CallbackData data = callbackQueue.front();
        callbackQueue.pop();
        
        if (callbackMap.size() > 0) {
            auto it = callbackMap.begin();
            Local<Function> callback = Local<Function>::New(isolate, it->second);
            
            Local<String> jsonStr = String::NewFromUtf8(isolate, data.jsonData.c_str()).ToLocalChecked();
            
            // Wrap JSON string in parentheses to make it a valid JavaScript expression
            std::string wrappedJson = "(" + data.jsonData + ")";
            Local<String> wrappedJsonStr = String::NewFromUtf8(isolate, wrappedJson.c_str()).ToLocalChecked();
            
            TryCatch tryCatch(isolate);
            Local<Script> script;
            Local<Value> parsedData;
            
            if (Script::Compile(context, wrappedJsonStr).ToLocal(&script) && script->Run(context).ToLocal(&parsedData) && parsedData->IsObject()) {
                Local<Value> argv[] = { parsedData };
                callback->Call(context, Null(isolate), 1, argv).ToLocalChecked();
            } else {
                Local<Object> emptyObj = Object::New(isolate);
                Local<Value> argv[] = { emptyObj };
                callback->Call(context, Null(isolate), 1, argv).ToLocalChecked();
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
    int32_t id = nextCallbackId++;
    callbackMap[id] = Global<Function>(isolate, callback);
    return id;
}

void CallNodeCallback(const char* jsonData) {
    CallbackData data;
    data.jsonData = jsonData;
    
    {
        std::lock_guard<std::mutex> lock(queueMutex);
        callbackQueue.push(data);
    }
    
    uv_async_send(&async_handle);
}

void ClearCallback(int32_t callbackId) {
    auto it = callbackMap.find(callbackId);
    if (it != callbackMap.end()) {
        it->second.Reset();
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
    
    args.GetReturnValue().Set(True(isolate));
}

void UnloadGoLibrary(const FunctionCallbackInfo<Value>& args) {
    Isolate* isolate = args.GetIsolate();
    
    for (auto& pair : callbackMap) {
        pair.second.Reset();
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
    __declspec(dllexport) void CallCallback(const char* jsonData) {
        CallNodeCallback(jsonData);
    }
    
    __declspec(dllexport) void FreeCallback(int32_t callbackId) {
        ClearCallback(callbackId);
    }
    
    __declspec(dllexport) void RegisterGoCallback(CallNodeCallbackType fn) {
        gCallNodeCallback = fn;
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
    
    // Register callback function in Go DLL
    typedef void (*RegisterGoCallbackType)(CallNodeCallbackType);
    RegisterGoCallbackType registerGoCallback = (RegisterGoCallbackType)GetProcAddress(hDLL, "RegisterGoCallback");
    if (registerGoCallback != NULL) {
        registerGoCallback(CallNodeCallback);
        printf("Called RegisterGoCallback\\n");
    } else {
        printf("RegisterGoCallback not found in DLL\\n");
    }
    
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
        
        // Register callback function in Go DLL
        typedef void (*RegisterGoCallbackType)(CallNodeCallbackType);
        RegisterGoCallbackType registerGoCallback = (RegisterGoCallbackType)GetProcAddress(hDLL, "RegisterGoCallback");
        if (registerGoCallback != NULL) {
            registerGoCallback(CallNodeCallback);
            printf("Called RegisterGoCallback\\n");
        } else {
            printf("RegisterGoCallback not found in DLL\\n");
        }
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

func toCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return strings.ToLower(s[:1]) + s[1:]
	}
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

func parseGoExports(goFile string) ([]string, error) {
	content, err := os.ReadFile(goFile)
	if err != nil {
		return nil, err
	}

	var functions []string
	lines := strings.Split(string(content), "\n")

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
