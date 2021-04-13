package renderdocs

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	fs "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/fs"
	"github.com/pkg/errors"
)

const cmdPrefix = "renderdocs"

type Command struct {
	SrcDir      string
	PartialsDir string
	DstDir      string
}

func (cmd *Command) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet(cmdPrefix, flag.ExitOnError)
	fs.StringVar(&cmd.SrcDir, "src", "docs", "docs/ folder to copy from.")
	fs.StringVar(&cmd.PartialsDir, "partials", "docs-partials", "docs-partials/ folder containing all mdx partials.")
	fs.StringVar(&cmd.DstDir, "dst", ".docs", "output folder.")
	return fs
}

func (cmd *Command) Help() string {
	return "Renders .mdx docs from static files and partials."
}

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
		return errors.Wrap(err, "unable to parse flags")
	}
	log.Printf("Copying %q to %q", cmd.SrcDir, cmd.DstDir)
	if err := fs.SyncDir(cmd.SrcDir, cmd.DstDir); err != nil {
		return errors.Wrap(err, "SyncDir failed")
	}
	log.Printf("Replacing @import '...' calls in %s", cmd.DstDir)

	return RenderDocsFolder(cmd.DstDir, cmd.PartialsDir)
}

func RenderDocsFolder(folder, partials string) error {
	entries, err := ioutil.ReadDir(folder)
	if err != nil {
		return errors.Wrapf(err, "cannot read directory %s", folder)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(folder, entry.Name())
		if entry.IsDir() {
			if err = RenderDocsFolder(entryPath, partials); err != nil {
				return err
			}
		} else {
			if err = renderDocsFile(entryPath, partials); err != nil {
				return errors.Wrap(err, "renderDocsFile")
			}
		}
	}
	return nil
}

var (
	includeStr = []byte("@include '")
)

func renderDocsFile(filePath, partialsDir string) error {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	for i := 0; i+len(includeStr) < len(f); i++ {
		if f[i] != '@' {
			continue
		}
		if diff := bytes.Compare(f[i:i+len(includeStr)], includeStr); diff != 0 {
			continue
		}
		ii := i + len(includeStr)
		for ; ii < len(f); ii++ {
			if f[ii] == '\'' {
				break
			}
		}
		if ii == len(includeStr) || f[ii] != '\'' {
			log.Printf("Unclosed @include quote at %d in %s", ii, filePath)
		}
		partialPath := string(f[i+len(includeStr) : ii])
		partial, err := os.ReadFile(filepath.Join(partialsDir, partialPath))
		if err != nil {
			return err
		}
		f = append(f[:i], append(partial, f[ii+1:]...)...)
	}

	return os.WriteFile(filePath, f, 0)
}

func (cmd *Command) Synopsis() string {
	return ""
}
