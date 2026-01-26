package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const bindingGypTemplate = `{
  "targets": [
    {
      "target_name": "{{.ModuleName}}",
      "cflags!": [ "-fno-exceptions" ],
      "cflags_cc!": [ "-fno-exceptions" ],
      "sources": [ 
        "wrapper.cc"
      ],
      "msvs_settings": {
        "VCLinkerTool": {
          "AdditionalLibraryDirectories": [
            "$(ConfigurationPath)"
          ]
        },
        "VCResourceCompilerTool": {
          "ResourceFileName": "$(IntDir)%(Filename).res"
        }
      },
      "conditions": [
        ["OS=='win'", {
          "sources": [
            "{{.ModuleName}}.rc"
          ]
        }],
        ["OS=='linux'", {
          "cflags!": [ "-fno-exceptions" ],
          "cflags_cc!": [ "-fno-exceptions" ]
        }],
        ["OS=='mac'", {
          "ldflags": [ "-l{{.ModuleName}}" ],
          "xcode_settings": {
            "CLANG_CXX_LIBRARY": "libc++",
            "CLANG_CXX_FLAGS": "-std=c++14 -fno-exceptions"
          }
        }]
      ]
    }
  ]
}
`

type BindingGypData struct {
	ModuleName string
}

func generateBindingGyp(workDir, moduleName string) error {
	bindingPath := filepath.Join(workDir, "binding.gyp")

	tmpl, err := template.New("binding.gyp").Parse(bindingGypTemplate)
	if err != nil {
		return fmt.Errorf("解析binding.gyp模板失败: %w", err)
	}

	data := BindingGypData{
		ModuleName: moduleName,
	}

	file, err := os.Create(bindingPath)
	if err != nil {
		return fmt.Errorf("创建binding.gyp文件失败: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("写入binding.gyp内容失败: %w", err)
	}

	fmt.Printf("生成binding.gyp: %s\n", bindingPath)
	return nil
}
