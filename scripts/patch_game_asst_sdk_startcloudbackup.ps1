$path = 'C:\Users\admin\Documents\goproject\xlasgame\game-asst-sdk\main.go'
$content = Get-Content -Raw $path

$old = @'
//export StartCloudBackup
func StartCloudBackup(req *C.char, callbackType *C.char) {
	callback := newSDKCallback(callbackType)
	if !callback.available() {
		log.Printf("StartCloudBackup callback not registered (callbackType=%s)\n", callback.token)
		return
	}

	reqStr := C.GoString(req)
	log.Printf("StartCloudBackup req: %v\n", reqStr)

	var reqData StartCloudBackupReq
	if err := json.Unmarshal([]byte(reqStr), &reqData); err != nil {
		log.Printf("StartCloudBackup parameter error: %v\n", err)
		callback.send(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	log.Printf("StartCloudBackup req: %v\n", reqData)
	config.SetASBaseURL(reqData.BaseURL)
	config.SetJWTToken(reqData.JWTToken)

	opts := &save.BackupOptions{
		GameID:         reqData.GameID,
		GameExecutable: reqData.GameExecutable,
		GameName:       reqData.GameName,
		InstallDir:     reqData.InstallDir,
		JWTToken:       reqData.JWTToken,
		ServerURL:      reqData.BaseURL,
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("StartCloudBackup panic: %v\n", r)
			callback.send(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("StartCloudBackup panic: %v", r),
			})
		}
	}()

	err := save.StartCloudBackup(nil, opts)
	if err != nil {
		log.Printf("StartCloudBackup failed: %v\n", err)
		callback.send(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	callback.send(map[string]interface{}{
		"success": true,
		"error":   "ok",
	})
}
'@

$new = @'
//export StartCloudBackup
func StartCloudBackup(req *C.char, callbackType *C.char) {
	callback := newSDKCallback(callbackType)
	if !callback.available() {
		log.Printf("StartCloudBackup callback not registered (callbackType=%s)\n", callback.token)
		return
	}

	reqStr := C.GoString(req)
	callback.keep()

	go func(cb sdkCallback, rawReq string) {
		defer cb.free()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("StartCloudBackup panic: %v\n", r)
				cb.send(map[string]interface{}{
					"success": false,
					"error":   fmt.Sprintf("StartCloudBackup panic: %v", r),
				})
			}
		}()

		log.Printf("StartCloudBackup req: %v\n", rawReq)

		var reqData StartCloudBackupReq
		if err := json.Unmarshal([]byte(rawReq), &reqData); err != nil {
			log.Printf("StartCloudBackup parameter error: %v\n", err)
			cb.send(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		log.Printf("StartCloudBackup req: %v\n", reqData)
		config.SetASBaseURL(reqData.BaseURL)
		config.SetJWTToken(reqData.JWTToken)

		opts := &save.BackupOptions{
			GameID:         reqData.GameID,
			GameExecutable: reqData.GameExecutable,
			GameName:       reqData.GameName,
			InstallDir:     reqData.InstallDir,
			JWTToken:       reqData.JWTToken,
			ServerURL:      reqData.BaseURL,
		}

		err := save.StartCloudBackup(nil, opts)
		if err != nil {
			log.Printf("StartCloudBackup failed: %v\n", err)
			cb.send(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		cb.send(map[string]interface{}{
			"success": true,
			"error":   "ok",
		})
	}(callback, reqStr)
}
'@

if (-not $content.Contains($old)) {
	throw 'StartCloudBackup block not found'
}

$content = $content.Replace($old, $new)
Set-Content -Path $path -Value $content -NoNewline
