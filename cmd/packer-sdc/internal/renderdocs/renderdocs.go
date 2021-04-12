package renderdocs

import (
	"flag"
	"log"
	"os"
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
	if err := os.MkdirAll(cmd.DstDir, 0644); err != nil {
		log.Fatalf("mkdir: %s", err)
	}
	// cp -r src dst
	// while rendering stuff in dst: render stuff in dst
	return 0
}

func (cmd *Command) Synopsis() string {
	return ""
}
