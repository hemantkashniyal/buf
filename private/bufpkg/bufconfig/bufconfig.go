// Copyright 2020-2021 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package bufconfig contains the configuration functionality.
package bufconfig

import (
	"context"

	"github.com/bufbuild/buf/private/bufpkg/bufcheck/bufbreaking/bufbreakingconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/buflint/buflintconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/storage"
)

const (
	// ExternalConfigV1FilePath is the default configuration file path.
	ExternalConfigV1FilePath = "buf.yaml"

	// ExternalConfigV1Beta1FilePath is the v1beta1 file path.
	ExternalConfigV1Beta1FilePath = "buf.yaml"

	// V1Version is the v1 version.
	V1Version = "v1"

	// V1Beta1Version is the v1beta1 version.
	V1Beta1Version = "v1beta1"

	// backupExternalConfigV1FilePath is another acceptable configuration file path for v1.
	//
	// Originally we thought we were going to move to buf.mod, and had this around for
	// a while, but then reverted back to buf.yaml. We still need to support buf.mod as
	// we released with it, however.
	backupExternalConfigV1FilePath = "buf.mod"
)

var (
	// All versions are all the versions in order.
	AllVersions = []string{
		V1Beta1Version,
		V1Version,
	}

	// AllConfigFilePaths are all acceptable config file paths without overrides.
	//
	// These are in the order we should check.
	AllConfigFilePaths = []string{
		ExternalConfigV1FilePath,
		backupExternalConfigV1FilePath,
	}
)

// Config is the user config.
type Config struct {
	Version        string
	ModuleIdentity bufmoduleref.ModuleIdentity
	Build          *bufmoduleconfig.Config
	Breaking       *bufbreakingconfig.Config
	Lint           *buflintconfig.Config
}

// GetConfigForBucket gets the Config for the YAML data at ConfigFilePath.
//
// If the data is of length 0, returns the default config.
func GetConfigForBucket(ctx context.Context, readBucket storage.ReadBucket) (*Config, error) {
	return getConfigForBucket(ctx, readBucket)
}

// GetConfigForData gets the Config for the given JSON or YAML data.
//
// If the data is of length 0, returns the default config.
func GetConfigForData(ctx context.Context, data []byte) (*Config, error) {
	return getConfigForData(ctx, data)
}

// WriteConfig writes an initial configuration file into the bucket.
func WriteConfig(
	ctx context.Context,
	writeBucket storage.WriteBucket,
	options ...WriteConfigOption,
) error {
	return writeConfig(
		ctx,
		writeBucket,
		options...,
	)
}

// WriteConfigOption is an option for WriteConfig.
type WriteConfigOption func(*writeConfigOptions)

// WriteConfigWithDocumentationComments returns a new WriteConfigOption that documents the resulting configuration file.
func WriteConfigWithDocumentationComments() WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.documentationComments = true
	}
}

// WriteConfigWithUncomment returns a new WriteConfigOption that uncomments the resulting configuration file
// options that are commented out by default.
//
// If this option is used, WriteConfigWithDocumentationComments must also be used.
func WriteConfigWithUncomment() WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.uncomment = true
	}
}

// WriteConfigWithModuleIdentity returns a new WriteConfigOption that sets the name of the
// module to the given ModuleIdentity.
//
// The default is to not set the name.
//
// Only use in tests.
func WriteConfigWithModuleIdentity(moduleIdentity bufmoduleref.ModuleIdentity) WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.moduleIdentity = moduleIdentity
	}
}

// WriteConfigWithDependencyModuleReferences returns a new WriteConfigOption that sets the
// dependencies of the module.
//
// The default is to not have any dependencies.
// If this option is used, WriteConfigWithModuleIdentity must also be used.
//
// Only use in tests.
func WriteConfigWithDependencyModuleReferences(dependencyModuleReferences ...bufmoduleref.ModuleReference) WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.dependencyModuleReferences = dependencyModuleReferences
	}
}

// WriteConfigWithBreakingConfig returns a new WriteConfigOption that sets the breaking change
// config for the module.
//
// If this option is used then the version used must be consistent with the rest of the config.
func WriteConfigWithBreakingConfig(breakingConfig *bufbreakingconfig.Config) WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.breakingConfig = breakingConfig
	}
}

// WriteConfigWithLintConfig returns a new WriteConfigOption that sets the lint config for the module.
//
// If this option is used then the version used must be consistent with the rest of the config.
func WriteConfigWithLintConfig(lintConfig *buflintconfig.Config) WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.lintConfig = lintConfig
	}
}

// WriteConfigWithVersion returns a new WriteConfigOption that sets the version of the config
// being written.
//
// If this is not set, the default is v1beta1.
func WriteConfigWithVersion(version string) WriteConfigOption {
	return func(writeConfigOptions *writeConfigOptions) {
		writeConfigOptions.version = version
	}
}

// ReadConfigOS reads the configuration from the OS or an override, if any.
//
// ONLY USE IN CLI TOOLS.
func ReadConfigOS(
	ctx context.Context,
	readBucket storage.ReadBucket,
	options ...ReadConfigOSOption,
) (*Config, error) {
	return readConfigOS(
		ctx,
		readBucket,
		options...,
	)
}

// ReadConfigOSOption is an option for ReadConfig.
type ReadConfigOSOption func(*readConfigOSOptions)

// ReadConfigOSWithOverride sets the override.
//
// If override is set, this will first check if the override ends in .json or .yaml, if so,
// this reads the file at this path and uses it. Otherwise, this assumes this is configuration
// data in either JSON or YAML format, and unmarshals it.
//
// If no override is set, this reads ExternalConfigFilePath in the bucket.
func ReadConfigOSWithOverride(override string) ReadConfigOSOption {
	return func(readConfigOSOptions *readConfigOSOptions) {
		readConfigOSOptions.override = override
	}
}

// ExistingConfigFilePath checks if a configuration file exists, and if so, returns the path
// within the ReadBucket of this configuration file.
//
// Returns empty string and no error if no configuration file exists.
func ExistingConfigFilePath(ctx context.Context, readBucket storage.ReadBucket) (string, error) {
	for _, configFilePath := range AllConfigFilePaths {
		exists, err := storage.Exists(ctx, readBucket, configFilePath)
		if err != nil {
			return "", err
		}
		if exists {
			return configFilePath, nil
		}
	}
	return "", nil
}

// ExternalConfigV1Beta1 represents the on-disk representation of the Config
// at version v1beta1.
type ExternalConfigV1Beta1 struct {
	Version  string                                  `json:"version,omitempty" yaml:"version,omitempty"`
	Name     string                                  `json:"name,omitempty" yaml:"name,omitempty"`
	Deps     []string                                `json:"deps,omitempty" yaml:"deps,omitempty"`
	Build    bufmoduleconfig.ExternalConfigV1Beta1   `json:"build,omitempty" yaml:"build,omitempty"`
	Breaking bufbreakingconfig.ExternalConfigV1Beta1 `json:"breaking,omitempty" yaml:"breaking,omitempty"`
	Lint     buflintconfig.ExternalConfigV1Beta1     `json:"lint,omitempty" yaml:"lint,omitempty"`
}

// ExternalConfigV1 represents the on-disk representation of the Config
// at version v1.
type ExternalConfigV1 struct {
	Version  string                             `json:"version,omitempty" yaml:"version,omitempty"`
	Name     string                             `json:"name,omitempty" yaml:"name,omitempty"`
	Deps     []string                           `json:"deps,omitempty" yaml:"deps,omitempty"`
	Build    bufmoduleconfig.ExternalConfigV1   `json:"build,omitempty" yaml:"build,omitempty"`
	Breaking bufbreakingconfig.ExternalConfigV1 `json:"breaking,omitempty" yaml:"breaking,omitempty"`
	Lint     buflintconfig.ExternalConfigV1     `json:"lint,omitempty" yaml:"lint,omitempty"`
}

// ExternalConfigVersion defines the subset of all config
// file versions that is used to determine the configuration version.
type ExternalConfigVersion struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}
