$path = 'C:\Users\admin\Documents\goproject\xlasgame\game-asst-sdk\main.go'
$content = Get-Content -Raw $path

$content = $content -replace 'var gCallNodeCallback uintptr', @'
var (
	gCallNodeCallback uintptr
	gKeepNodeCallback uintptr
	gFreeNodeCallback uintptr
)
'@

$registerOld = @'
//export RegisterGoCallback
func RegisterGoCallback(fn uintptr) {
	gCallNodeCallback = fn
}
'@

$registerNew = @'
//export RegisterGoCallback
func RegisterGoCallback(fn uintptr) {
	gCallNodeCallback = fn
}

//export RegisterGoCallbackEx
func RegisterGoCallbackEx(callFn uintptr, keepFn uintptr, freeFn uintptr) {
	gCallNodeCallback = callFn
	gKeepNodeCallback = keepFn
	gFreeNodeCallback = freeFn
}
'@

$content = $content.Replace($registerOld, $registerNew)

$sendOld = @'
func (cb sdkCallback) send(data map[string]interface{}) {
	if cb.ptr == 0 {
		return
	}

	jsonData, _ := json.Marshal(data)
	cJSON := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(cJSON))

	if cb.id >= 0 {
		C.callCallbackWithId(unsafe.Pointer(cb.ptr), C.int32_t(cb.id), cJSON)
		return
	}

	C.callCallback(unsafe.Pointer(cb.ptr), cJSON)
}
'@

$sendNew = @'
func (cb sdkCallback) send(data map[string]interface{}) {
	if cb.ptr == 0 {
		return
	}

	jsonData, _ := json.Marshal(data)
	cJSON := C.CString(string(jsonData))
	defer C.free(unsafe.Pointer(cJSON))

	if cb.id >= 0 {
		C.callCallbackWithId(unsafe.Pointer(cb.ptr), C.int32_t(cb.id), cJSON)
		return
	}

	C.callCallback(unsafe.Pointer(cb.ptr), cJSON)
}

func (cb sdkCallback) keep() {
	if cb.id < 0 || gKeepNodeCallback == 0 {
		return
	}

	C.keepCallback(unsafe.Pointer(gKeepNodeCallback), C.int32_t(cb.id))
}

func (cb sdkCallback) free() {
	if cb.id < 0 || gFreeNodeCallback == 0 {
		return
	}

	C.freeCallback(unsafe.Pointer(gFreeNodeCallback), C.int32_t(cb.id))
}
'@

$content = $content.Replace($sendOld, $sendNew)

$scanGamesPattern = '(?s)//export ScanGames\r?\nfunc ScanGames\(req \*C\.char, callbackType \*C\.char\) \{.*?\r?\n\}\r?\n\r?\n//export Init'
$scanGamesReplacement = @'
//export ScanGames
func ScanGames(req *C.char, callbackType *C.char) {
	callback := newSDKCallback(callbackType)
	if !callback.available() {
		log.Printf("ScanGames callback not registered (callbackType=%s)\n", callback.token)
		return
	}

	reqStr := C.GoString(req)
	callback.keep()

	go func(cb sdkCallback, rawReq string) {
		defer cb.free()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("ScanGames panic: %v\n", r)
				cb.send(map[string]interface{}{
					"success": false,
					"error":   fmt.Sprintf("ScanGames panic: %v", r),
				})
			}
		}()

		log.Printf("ScanGames req: %v\n", rawReq)

		var reqData ScanGameReq
		if err := json.Unmarshal([]byte(rawReq), &reqData); err != nil {
			log.Printf("ScanGames parse request error: %v\n", err)
			cb.send(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		log.Printf("ScanGames reqData: %v\n", reqData)
		config.SetASBaseURL(reqData.BaseURL)

		gamesCh, err := scanner.ScanGames()
		if err != nil {
			log.Printf("ScanGames get game info error: %v\n", err)
			cb.send(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		for game := range gamesCh {
			firstScanFinished := game.FirstScanFinished
			callbackData := map[string]interface{}{}
			if firstScanFinished {
				callbackData = map[string]interface{}{
					"success": true,
					"error":   "ok",
					"data": map[string]interface{}{
						"firstScanFinished": firstScanFinished,
					},
				}
			} else {
				callbackData = map[string]interface{}{
					"success": true,
					"error":   "ok",
					"data": map[string]interface{}{
						"firstScanFinished": firstScanFinished,
						"game":              game,
					},
				}
			}

			cb.send(callbackData)
		}
	}(callback, reqStr)
}

//export Init
'@

$content = [regex]::Replace($content, $scanGamesPattern, $scanGamesReplacement)

Set-Content -Path $path -Value $content -NoNewline
