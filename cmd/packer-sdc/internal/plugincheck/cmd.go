package plugincheck

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

var (
	cmdPrefix = "plugin-check"

	//go:embed README.md
	readme string
)

type Command struct {
	Load string
}

func (cmd *Command) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet(cmdPrefix, flag.ExitOnError)
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
