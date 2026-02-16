package cfg

import (
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
)

// CompileConfig holds all configuration for the morphe-to-morpherepo compilation pipeline.
type CompileConfig struct {
	rcfg.MorpheLoadRegistryConfig

	// ExcludeModels lists model names to skip during generation.
	ExcludeModels []string

	// OutputDirPath is the directory for generated .repo files.
	OutputDirPath string
}
