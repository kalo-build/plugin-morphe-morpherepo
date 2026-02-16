package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile"
	"github.com/kalo-build/plugin-morphe-morpherepo/pkg/compile/cfg"
)

type PluginConfig struct {
	InputPath     string   `json:"inputPath"`
	OutputPath    string   `json:"outputPath"`
	ExcludeModels []string `json:"excludeModels,omitempty"`
	Verbose       bool     `json:"verbose,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-morpherepo <config>")
		os.Exit(3)
	}

	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(os.Args[1]), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(4)
	}

	if pluginConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required")
		os.Exit(12)
	}
	if pluginConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(13)
	}

	inputAbs, err := filepath.Abs(pluginConfig.InputPath)
	if err == nil {
		pluginConfig.InputPath = inputAbs
	}
	outputAbs, err := filepath.Abs(pluginConfig.OutputPath)
	if err == nil {
		pluginConfig.OutputPath = outputAbs
	}

	compileConfig := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      filepath.Join(pluginConfig.InputPath, "enums"),
			RegistryModelsDirPath:     filepath.Join(pluginConfig.InputPath, "models"),
			RegistryStructuresDirPath: filepath.Join(pluginConfig.InputPath, "structures"),
			RegistryEntitiesDirPath:   filepath.Join(pluginConfig.InputPath, "entities"),
		},
		ExcludeModels: pluginConfig.ExcludeModels,
		OutputDirPath: pluginConfig.OutputPath,
	}

	if compileErr := compile.MorpheToMorpheRepo(compileConfig); compileErr != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", compileErr)
		os.Exit(1)
	}

	os.Exit(0)
}
