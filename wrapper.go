package main

import (
	"os"
	"path/filepath"
	"strings"
)

func generateWrapperCC(workDir, moduleName string, functions []string) error {
	wrapperPath := filepath.Join(workDir, "wrapper.cc")

	funcDecls := ""
	funcImpls := ""
	nodeMethods := ""

	for _, fn := range functions {
		funcName := fn
		jsName := toCamelCase(fn)

		if strings.HasPrefix(fn, "Hello1") || strings.HasPrefix(fn, "Hello2") || strings.HasPrefix(fn, "Hello3") {
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
				"    if (args.Length() > 0 && args[0]->IsString()) {\n" +
				"        String::Utf8Value arg1Str(isolate, args[0]);\n" +
				"        arg1 = strdup(*arg1Str);\n" +
				"    }\n" +
				"    \n" +
				"    if (args.Length() > 1 && args[1]->IsString()) {\n" +
				"        String::Utf8Value arg2Str(isolate, args[1]);\n" +
				"        arg2 = strdup(*arg2Str);\n" +
				"    }\n" +
				"    \n" +
				"    const char* result = " + funcName + "Ptr(arg1, arg2);\n" +
				"    \n" +
				"    if (result != NULL) {\n" +
				"        args.GetReturnValue().Set(String::NewFromUtf8(isolate, result).ToLocalChecked());\n" +
				"    }\n" +
				"}\n\n"
		} else if strings.HasPrefix(fn, "Hello4") || strings.HasPrefix(fn, "Hello5") || strings.HasPrefix(fn, "Hello6") {
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
				"    const char* callbackType = \"default\";\n" +
				"    \n" +
				"    if (args.Length() > 0 && args[0]->IsString()) {\n" +
				"        String::Utf8Value arg1Str(isolate, args[0]);\n" +
				"        arg1 = strdup(*arg1Str);\n" +
				"    }\n" +
				"    \n" +
				"    if (args.Length() > 1 && args[1]->IsFunction()) {\n" +
				"        callbackType = \"callback\";\n" +
				"        RegisterCallback(isolate, args[1].As<Function>());\n" +
				"    }\n" +
				"    \n" +
				"    const char* result = " + funcName + "Ptr(arg1, callbackType);\n" +
				"    \n" +
				"    if (result != NULL) {\n" +
				"        args.GetReturnValue().Set(String::NewFromUtf8(isolate, result).ToLocalChecked());\n" +
				"    }\n" +
				"}\n\n"
		}

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

using namespace v8;

// Callback function type and global variable
typedef void (*CallNodeCallbackType)(const char*, const char*);
CallNodeCallbackType gCallNodeCallback = NULL;

` + funcDecls + `// Global variables
HMODULE hDLL = NULL;
std::map<int32_t, Global<Function>> callbackMap;
int32_t nextCallbackId = 0;
Isolate* gIsolate = NULL;

// Async callback handling
struct CallbackData {
    std::string callbackType;
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
            
            Local<Value> argv[] = {
                String::NewFromUtf8(isolate, data.callbackType.c_str()).ToLocalChecked(),
                String::NewFromUtf8(isolate, data.jsonData.c_str()).ToLocalChecked()
            };
            
            callback->Call(context, Null(isolate), 2, argv).ToLocalChecked();
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

void CallNodeCallback(const char* callbackType, const char* jsonData) {
    CallbackData data;
    data.callbackType = callbackType;
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
    __declspec(dllexport) void CallCallback(const char* callbackType, const char* jsonData) {
        CallNodeCallback(callbackType, jsonData);
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
    
    // Get path of current module (.node file)
    char modulePath[MAX_PATH];
    HMODULE hModule = NULL;
    if (GetModuleHandleExA(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS | GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
                             (LPCSTR)&Initialize, &hModule)) {
        GetModuleFileNameA(hModule, modulePath, MAX_PATH);
        
        // Get directory path
        std::string modulePathStr(modulePath);
        size_t lastSlash = modulePathStr.find_last_of("\\\\/");
        if (lastSlash != std::string::npos) {
            std::string moduleDir = modulePathStr.substr(0, lastSlash);
            std::string dllPath = moduleDir + "\\\\" + "` + moduleName + `.dll";
            hDLL = LoadLibraryA(dllPath.c_str());
            if (hDLL == NULL) {
                isolate->ThrowException(Exception::Error(
                    String::NewFromUtf8(isolate, ("Failed to load " + dllPath).c_str()).ToLocalChecked()));
                return;
            }
        } else {
            isolate->ThrowException(Exception::Error(
                String::NewFromUtf8(isolate, "Failed to get module directory").ToLocalChecked()));
            return;
        }
    } else {
        isolate->ThrowException(Exception::Error(
            String::NewFromUtf8(isolate, "Failed to get module handle").ToLocalChecked()));
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
		return s
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
