// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package version helps plugin creators set and track the plugin version using
// the same convenience functions used by the Packer core.
package version

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

// The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Package version helps plugin creators set and track the sdk version using
var Version = "0.6.4"

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
var VersionPrerelease = "dev"

// The metadata for the version, this is optional information to add around
// a particular release.
//
// This has no impact on the ordering of plugins, and is ignored for non-human eyes.
var VersionMetadata = ""

// SDKVersion is used by the plugin set to allow Packer to recognize
// what version of the sdk the plugin is.
var SDKVersion = NewPluginVersion(Version, VersionPrerelease, VersionMetadata)

// InitializePluginVersion initializes the SemVer and returns a version var.
//
// Deprecated: InitializePluginVersion does not support metadata out of the
// box, and should be replaced by either NewPluginVersion or NewRawVersion.
func InitializePluginVersion(vers, versionPrerelease string) *PluginVersion {
	return NewPluginVersion(vers, versionPrerelease, "")
}

// NewRawVersion is made for more freeform version strings. It won't accept
// much more than what `NewPluginVersion` already does, but is another
// convenient form to create a version if preferred.
//
// As NewRawVersion, if the version is invalid, it will panic.
func NewRawVersion(rawSemVer string) *PluginVersion {
	vers := version.Must(version.NewVersion(rawSemVer))

	if len(vers.Segments()) != 3 {
		panic(fmt.Sprintf("versions should only have 3 segments, %q had %d", rawSemVer, len(vers.Segments())))
	}

	return &PluginVersion{
		version:           vers.Core().String(),
		versionPrerelease: vers.Prerelease(),
		versionMetadata:   vers.Metadata(),
		semVer:            vers,
	}
}

// NewPluginVersion initializes the SemVer and returns a PluginVersion from it.
// If the provided "version" string is not valid, the call to version.Must
// will panic.
//
// This function should always be called in a package init() function to make
// sure that plugins are following proper semantic versioning and to make sure
// that plugins which aren't following proper semantic versioning crash
// immediately rather than later.
//
// If the core version number is empty, it will default to 0.0.0.
func NewPluginVersion(vers, versionPrerelease, versionMetadata string) *PluginVersion {
	var versionRawString = vers

	if versionRawString == "" {
		// Defaults to "0.0.0". Useful when binary is created for development purpose.
		versionRawString = "0.0.0"
	}

	if versionPrerelease != "" {
		versionRawString = fmt.Sprintf("%s-%s", versionRawString, versionPrerelease)
	}

	if versionMetadata != "" {
		versionRawString = fmt.Sprintf("%s+%s", versionRawString, versionMetadata)
	}

	return NewRawVersion(versionRawString)
}

type PluginVersion struct {
	// The main version number that is being run at the moment.
	version string
	// A pre-release marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	versionPrerelease string
	// Extra metadata that can be part of the version.
	//
	// This is legal in semver, and has to be the last part of the version
	// string, starting with a `+`.
	versionMetadata string
	// The Semantic Version of the plugin. Used for version constraint comparisons
	semVer *version.Version
}

func (p *PluginVersion) SetMetadata(meta string) {
	p.versionMetadata = meta
}

func (p *PluginVersion) FormattedVersion() string {
	versionString := p.semVer.String()

	// Given there could be some metadata already, we add the commit to the
	// reported version as part of the metadata, with a `-` spearator if
	// the metadata is already there, otherwise we make it the metadata
	if GitCommit != "" {
		versionString = fmt.Sprintf("%s (%s)", versionString, GitCommit)
	}

	return versionString
}

func (p *PluginVersion) SemVer() *version.Version {
	return p.semVer
}

func (p *PluginVersion) GetVersion() string {
	return p.version
}

func (p *PluginVersion) GetVersionPrerelease() string {
	return p.versionPrerelease
}

func (p *PluginVersion) GetMetadata() string {
	return p.versionMetadata
}

// String returns the complete version string, including prerelease
func (p *PluginVersion) String() string {
	return p.semVer.String()
}
