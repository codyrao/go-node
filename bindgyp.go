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
        "wrapper.cc",
        "{{.ModuleName}}.dll"
      ],
      "msvs_settings": {
        "VCLinkerTool": {
          "AdditionalLibraryDirectories": [
            "$(ConfigurationPath)"
          ],
          "AdditionalOptions": [ "/OPT:REF", "/OPT:ICF" ]
        },
        "VCResourceCompilerTool": {
          "ResourceFileName": "$(IntDir)%(Filename).res"
        },
        "VCCLCompilerTool": {
          "Optimization": 2,
          "FavorSizeOrSpeed": 2,
          "AdditionalOptions": [ "/wd4018", "/wd4996" ]
        }
      },
      "conditions": [
        ["OS=='win'", {
          "sources": [
            "wrapper.cc",
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
		return fmt.Errorf("Failed to parse binding.gyp template: %w", err)
	}

	data := BindingGypData{
		ModuleName: moduleName,
	}

	file, err := os.Create(bindingPath)
	if err != nil {
		return fmt.Errorf("Failed to create binding.gyp file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("Failed to write binding.gyp content: %w", err)
	}

	fmt.Printf("Generated binding.gyp: %s\n", bindingPath)
	return nil
}
