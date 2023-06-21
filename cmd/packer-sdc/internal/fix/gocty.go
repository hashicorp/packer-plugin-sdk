// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fix

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
)

const (
	oldPath    string = "github.com/zclconf/go-cty"
	newPath    string = "github.com/nywilken/go-cty"
	newVersion string = "1.12.1"
)

type GoCtyFix struct {
	OldPath, NewPath, NewVersion string
	data                         []byte
}

// NewGoCtyFixer is an entry to the go-cty fix command
func NewGoCtyFixer() GoCtyFix {
	return GoCtyFix{
		OldPath:    oldPath,
		NewPath:    newPath,
		NewVersion: newVersion,
	}
}

func (f GoCtyFix) modFileFormmatedVersion() string {
	return fmt.Sprintf("v%s", f.NewVersion)
}

func (f GoCtyFix) isUnfixed(dir string) (bool, error) {
	modFilePath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(modFilePath); err != nil {
		return false, errors.Wrap(err, "failed to find go.mod file")
	}
	data, err := os.ReadFile(modFilePath)
	if err != nil {
		return false, errors.Wrap(err, "failed to read go.mod file")
	}
	mf, err := modfile.Parse(modFilePath, data, nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse go.mod file")
	}

	// Basic go.mod with no module dependencies
	if len(mf.Require) == 0 {
		return false, nil
	}

	// Brute force will be better to Sort then search
	var found bool
	for _, req := range mf.Require {
		if req.Mod.Path == f.OldPath {
			found = true
			break
		}
	}

	if !found {
		return false, nil
	}

	if len(mf.Replace) == 0 {
		return true, nil
	}
	// what happens with multiple replace
	for _, r := range mf.Replace {
		if r.Old.Path != f.OldPath {
			continue
		}
		if r.New.Path != f.NewPath {
			return false, errors.New("found unexpected replace for " + r.Old.Path)
		}
		return r.New.Version != f.modFileFormmatedVersion(), nil
	}
	return true, nil
}

func (f GoCtyFix) Fix(dir string, check bool) error {
	ok, err := f.isUnfixed(dir)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if check {
		fmt.Printf("%s %5s\n", "gocty", "Unfixed!")
		return nil
	}

	modFilePath := filepath.Join(dir, "go.mod")
	info, err := os.Stat(modFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to find plugin go.mod file")
	}
	data, err := os.ReadFile(modFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read plugin go.mod file")
	}
	mf, err := modfile.Parse(modFilePath, data, nil)
	if err != nil {
		return errors.Wrap(err, "failed to parse plugin go.mod file")
	}

	if err := mf.DropReplace(f.OldPath, ""); err != nil {
		return errors.Wrap(err, "failed to apply gocty fix")
	}
	commentSuffix := " // add by packer-sdc fix as noted in github.com/hashicorp/packer-plugin-sdk/issues/187"
	if err := mf.AddReplace(f.OldPath, "", f.NewPath, f.modFileFormmatedVersion()+commentSuffix); err != nil {
		return errors.Wrap(err, "failed to apply gocty fix")
	}

	bytes, _ := mf.Format()
	return os.WriteFile(modFilePath, bytes, info.Mode())
}
