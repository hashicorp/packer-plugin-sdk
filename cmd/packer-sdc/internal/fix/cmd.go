// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package fix

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/diff"
)

const cmdPrefix string = "fix"

// Command is the base entry for the fix sub-command.
type Command struct {
	Dir  string
	Diff bool
}

// Flags contain the default flags for the fix sub-command.
func (cmd *Command) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet(cmdPrefix, flag.ExitOnError)
	fs.BoolVar(&cmd.Diff, "diff", false, "if set prints the differences a rewrite would introduce.")
	return fs
}

// Help displays usage for the command.
func (cmd *Command) Help() string {
	var s strings.Builder
	for _, fix := range availableFixes {
		s.WriteString(fmt.Sprintf("  %s\t\t%s\n", fix.name, fix.description))
	}

	helpText := `
Usage: packer-sdc fix [options] directory

  Fix rewrites parts of the plugin codebase to address known issues or
  common workarounds used within plugins consuming the Packer plugin SDK.

Options:
  -diff		If the -diff flag is set fix prints the differences an applied fix would introduce.

Available fixes:
%s`
	return fmt.Sprintf(helpText, s.String())
}

// Run executes the command
func (cmd *Command) Run(args []string) int {
	if err := cmd.run(args); err != nil {
		fmt.Printf("%v", err)
		return 1
	}
	return 0
}

func (cmd *Command) run(args []string) error {
	f := cmd.Flags()
	err := f.Parse(args)
	if err != nil {
		return errors.New("unable to parse flags for fix command")
	}

	if f.NArg() != 1 {
		err := fmt.Errorf("packer-sdc fix: missing directory argument\n%s", cmd.Help())
		return err
	}

	dir := f.Arg(0)
	if dir == "." || dir == "./..." {
		dir, _ = os.Getwd()
	}

	info, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		return errors.New("a plugin root directory must be specified or a dot for the current directory")
	}

	if !info.IsDir() {
		return errors.New("a plugin root directory must be specified or a dot for the current directory")
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return errors.New("unable to determine the absolute path for the provided plugin root directory")
	}
	cmd.Dir = dir

	return processFiles(cmd.Dir, cmd.Diff)
}

func (cmd *Command) Synopsis() string {
	return "Rewrites parts of the plugin codebase to address known issues or common workarounds within plugins consuming the Packer plugin SDK."
}

func processFiles(rootDir string, showDiff bool) error {
	srcFiles := make(map[string][]byte)
	fixedFiles := make(map[string][]byte)

	var cmdApplyErrs error
	for _, f := range availableFixes {
		matches, err := f.scan(rootDir)
		if err != nil {
			return fmt.Errorf("failed to apply %s fix: %s", f.name, err)
		}

		//matches contains all files to apply the said fix on
		for _, filename := range matches {
			if _, ok := srcFiles[filename]; !ok {
				bs, err := os.ReadFile(filename)
				if err != nil {
					cmdApplyErrs = multierror.Append(cmdApplyErrs, err)
				}
				srcFiles[filename] = append([]byte{}, bs...)
			}

			fixedData, ok := fixedFiles[filename]
			if !ok {
				fixedData = append([]byte{}, srcFiles[filename]...)
			}

			fixedData, err := f.fix(filename, fixedData)
			if err != nil {
				cmdApplyErrs = multierror.Append(cmdApplyErrs, err)
				continue
			}
			if bytes.Equal(fixedData, srcFiles[filename]) {
				continue
			}
			fixedFiles[filename] = fixedData
		}
	}

	if cmdApplyErrs != nil {
		return cmdApplyErrs
	}

	if showDiff {
		for filename, fixedData := range fixedFiles {
			diff.Text(filename, "Fixed: "+filename, string(srcFiles[filename]), string(fixedData), os.Stdout)
		}
		return nil
	}

	for filename, fixedData := range fixedFiles {
		fmt.Println(filename)
		info, _ := os.Stat(filename)
		os.WriteFile(filename, fixedData, info.Mode())
	}

	return nil
}
