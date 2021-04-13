package renderdocs

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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
				return errors.Wrap(err, "copying directory failed")
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
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	i := int64(0)
	for ; i < fileStat.Size(); i++ {
		ctt := make([]byte, 1)
		if _, err = file.ReadAt(ctt, i); err != nil {
			return err
		}

		// byte by byte, look for an '@'
		if ctt[0] != '@' {
			continue
		}
		//detected an @
		found := make([]byte, len(includeStr))
		if _, err := file.ReadAt(found, i); err != nil {
			continue
		}
		if diff := bytes.Compare(includeStr, found); diff != 0 {
			continue
		}
		// found `@include '`
		var path []byte
		ii := i + int64(len(includeStr))
		for i < fileStat.Size() {
			if _, err = file.ReadAt(ctt, ii); err != nil {
				return err
			}
			if ctt[0] == '\'' {
				ii++
				// found end of path
				break
			}
			if ctt[0] == '\n' {
				return fmt.Errorf("Unclosed @include quote at %d in %s", ii, filePath)
			}
			path = append(path, ctt...)
			ii++
			if ii == fileStat.Size() {
				return fmt.Errorf("Unclosed @include quote at %d in %s", ii, filePath)
			}
		}
		partialPath := string(path)
		partialFile, err := os.Open(filepath.Join(partialsDir, partialPath))
		if err != nil {
			return errors.Wrapf(err, "Failed to open partial at %q", partialPath)
		}
		partialFileStat, err := partialFile.Stat()
		if err != nil {
			return err
		}

		// copy rest of text from after @include '...' backwards onto file
		offset := partialFileStat.Size() - int64(len(includeStr)+len(partialPath)+1)
		for iii := fileStat.Size() - 1; ii < iii; iii-- {
			if _, err = file.ReadAt(ctt, iii); err != nil {
				return err
			}
			if _, err := file.WriteAt(ctt, iii+offset); err != nil {
				return err
			}
		}
		if err := file.Sync(); err != nil {
			return err
		}
		fileStat, err = file.Stat()
		if err != nil {
			return err
		}

		// Write content of partial here
		if _, err := file.Seek(i, 0); err != nil {
			return err
		}
		if _, err := io.Copy(file, partialFile); err != nil {
			return err
		}
		if err := file.Sync(); err != nil {
			return err
		}

		i += partialFileStat.Size()
	}

	return nil
}

func (cmd *Command) Synopsis() string {
	return ""
}
