// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fix

import (
	_ "embed" //used for embedding the command README
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const cmdPrefix string = "fix"

var (
	// availableFixes to apply to a plugin - refer to init func
	availableFixes fixes

	//go:embed README.md
	readme string
)

// Command is the base entry for the fix sub-command
type Command struct {
	Dir   string
	Check bool
}

func (cmd *Command) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet(cmdPrefix, flag.ExitOnError)
	fs.BoolVar(&cmd.Check, "check", false, "Check plugin for potential fixes [dry-run].")
	return fs
}

// Help displays usage for the command.
func (cmd *Command) Help() string {
	return "\n" + readme
}

// Run executes the command
func (cmd *Command) Run(args []string) int {
	if err := cmd.run(args); err != nil {
		log.Printf("%v", err)
		return 1
	}
	return 0
}

func (cmd *Command) run(args []string) error {
	f := cmd.Flags()
	err := f.Parse(args)
	if err != nil {
		return errors.Wrap(err, "[packer-sdc fix] unable to parse flags")
	}

	if f.NArg() != 1 {
		cmd.Help()
		return errors.New("a plugin root directory must be specified")
	}

	dir := f.Arg(0)
	if dir == "." {
		dir, _ = os.Getwd()
	}

	info, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		return errors.New("a plugin root directory must be specified")
	}

	if !info.IsDir() {
		return errors.New("a plugin root directory must be specified")
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return errors.New("unable to determine the absolute path for the plugin root directory")
	}
	cmd.Dir = dir

	for _, f := range availableFixes {
		err = f.Fix(cmd.Dir, cmd.Check)
		if err != nil {
			return errors.Errorf("failed to apply %s fix: %s", f.Name, err)
		}
	}
	return nil
}

func (cmd *Command) Synopsis() string {
	return "Rewrites parts of the plugin codebase to address known issues or common workarounds within the plugin API."
}

// Fixer applies all defined fixes on the provide plugin dir.
// A dryRun argument is available to check if a fix will be applied without executing.
// A Fixer should be idempotent.
type Fixer interface {
	Fix(pluginRootDir string, dryRun bool) error
}
type fix struct {
	Name, Description string
	Fixer
}

type fixes []fix

func init() {
	availableFixes = []fix{
		{
			Name:        "gocty",
			Description: "Adds a replace directive for github.com/zclconf/go-cty to github.com/nywilken/go-cty",
			Fixer:       NewGoCtyFixer(),
		},
	}
}
