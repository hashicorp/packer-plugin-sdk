package plugincheck

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	cmdPrefix = "plugin-check"

	//go:embed README.md
	readme string
)

type Command struct {
	Docs bool
	Load string
}

func (cmd *Command) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet(cmdPrefix, flag.ExitOnError)
	fs.BoolVar(&cmd.Docs, "docs", false, "flag to indicate that documentation files should be checked.")
	fs.StringVar(&cmd.Load, "load", "", "flag to check if plugin can be loaded by Packer and is compatible with HCL2.")
	return fs
}

func (cmd *Command) Help() string {
	cmd.Flags().Usage()
	return "\n" + readme
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

	if f.NFlag() == 0 {
		cmd.Help()
		return errors.New("No option passed")
	}

	if cmd.Docs {
		if err := checkDocumentation(); err != nil {
			return err
		}
		fmt.Printf("Plugin successfully passed docs check.\n")
	}

	if len(cmd.Load) == 0 {
		return nil
	}

	if err := checkPluginName(cmd.Load); err != nil {
		return err
	}

	// if err := discoverAndLoad(); err != nil {
	// 	fmt.Printf(err.Error())
	// 	os.Exit(2)
	// }
	// fmt.Printf("Plugin successfully passed compatibility check.\n")
	return nil
}

// checkDocumentation looks for the presence of a docs folder with mdx files inside.
// It is not possible to predict the number of mdx files for a given plugin.
// Because of that, finding one file inside the folder is enough to validate the docs existence.
func checkDocumentation() error {
	// TODO: this should be updated once we have defined what's going to be for plugin's docs
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	docsDir := wd + "/docs"
	stat, err := os.Stat(docsDir)
	if err != nil {
		return fmt.Errorf("could not find docs folter: %s", err.Error())
	}
	if !stat.IsDir() {
		return fmt.Errorf("expecting docs do be a directory of mdx files")
	}

	var mdxFound bool
	_ = filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".mdx" {
			mdxFound = true
			return io.EOF
		}
		return nil
	})

	if mdxFound {
		fmt.Printf("a mdx file was found inside the docs folder\n")
		return nil
	}
	return fmt.Errorf("no docs files found, make sure to have the docs in place before releasing")
}

// checkPluginName checks for the possible valid names for a plugin, packer-plugin-* or packer-[builder|provisioner|post-processor]-*.
// If the name is prefixed with `packer-[builder|provisioner|post-processor]-`, packer won't be able to install it,
// therefore a WARNING will be shown.
func checkPluginName(name string) error {
	if strings.HasPrefix(name, "packer-plugin-") {
		return nil
	}
	if strings.HasPrefix(name, "packer-builder-") ||
		strings.HasPrefix(name, "packer-provisioner-") ||
		strings.HasPrefix(name, "packer-post-processor-") {
		fmt.Printf("\n[WARNING] Plugin is named with old prefix `packer-[builder|provisioner|post-processor]-{name})`. " +
			"These will be detected but Packer cannot install them automatically. " +
			"The plugin must be a multi-component plugin named packer-plugin-{name} to be installable through the `packer init` command.\n")
		return nil
	}
	return fmt.Errorf("plugin's name is not valid")
}

func (cmd *Command) Synopsis() string {
	return "Tell wether a plugin release looks valid for Packer."
}
