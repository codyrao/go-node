/* Auto-generated callback header for go-node */
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
