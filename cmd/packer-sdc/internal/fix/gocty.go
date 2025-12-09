// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package fix

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/mod/modfile"
)

const (
	sdkPath     string = "github.com/hashicorp/packer-plugin-sdk"
	oldPath     string = "github.com/zclconf/go-cty"
	newPath     string = "github.com/nywilken/go-cty"
	newVersion  string = "1.13.3"
	modFilename string = "go.mod"
)

var goctyFix = fix{
	name:        "gocty",
	description: "Adds a replace directive for github.com/zclconf/go-cty to github.com/nywilken/go-cty",
	scan:        modPaths,
	fixer: goCtyFix{
		OldPath:    oldPath,
		NewPath:    newPath,
		NewVersion: newVersion,
	},
}

// modPaths scans the incoming dir for potential go.mod files to fix.
func modPaths(dir string) ([]string, error) {
	paths := []string{
		path.Join(dir, modFilename),
	}
	return paths, nil

}

type goCtyFix struct {
	OldPath, NewPath, NewVersion string
}

func (f goCtyFix) modFileFormattedVersion() string {
	return fmt.Sprintf("v%s", f.NewVersion)
}

// Fix applies a replace directive in a projects go.mod file for f.OldPath to f.NewPath.
// This fix applies to the replacement of github.com/zclconf/go-cty, as described in https://github.com/hashicorp/packer-plugin-sdk/issues/187
// The return data contains the data file with the applied fix. In cases where the fix is already applied or not needed the original data is returned.
func (f goCtyFix) fix(modFilePath string, data []byte) ([]byte, error) {
	if _, err := os.Stat(modFilePath); err != nil {
		return nil, fmt.Errorf("failed to find go.mod file %s", modFilePath)
	}

	mf, err := modfile.Parse(modFilePath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse go.mod file: %v", modFilePath, err)
	}

	// fix doesn't apply to go.mod with no module dependencies
	if len(mf.Require) == 0 {
		return data, nil
	}

	var requiresSDK, requiresGoCty bool
	for _, req := range mf.Require {
		if req.Mod.Path == sdkPath || mf.Module.Mod.Path == sdkPath {
			requiresSDK = true
		}
		if req.Mod.Path == f.OldPath {
			requiresGoCty = true
		}

		if requiresSDK && requiresGoCty {
			break
		}
	}

	if !(requiresSDK && requiresGoCty) {
		return data, nil
	}

	for _, r := range mf.Replace {
		if r.Old.Path != f.OldPath {
			continue
		}

		if r.New.Path != f.NewPath {
			return nil, fmt.Errorf("%s: found unexpected replace for %s", modFilePath, r.Old.Path)
		}

		if r.New.Version == f.modFileFormattedVersion() {
			return data, nil
		}
	}

	if err := mf.DropReplace(f.OldPath, ""); err != nil {
		return nil, fmt.Errorf("%s: failed to drop previously added replacement fix %v", modFilePath, err)
	}

	commentSuffix := " // added by packer-sdc fix as noted in github.com/hashicorp/packer-plugin-sdk/issues/187"
	if err := mf.AddReplace(f.OldPath, "", f.NewPath, f.modFileFormattedVersion()+commentSuffix); err != nil {
		return nil, fmt.Errorf("%s: failed to apply go-cty fix: %v", modFilePath, err)
	}

	newData, err := mf.Format()
	if err != nil {
		return nil, err
	}
	return newData, nil
}
