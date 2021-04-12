package generate

import (
	"flag"
)

type GenerateDocsCMD struct {
	outputDir string
	extension string
}

func (cmd *GenerateDocsCMD) Help() string {
	return "Generates docs from static files and partials"
}

func (cmd *GenerateDocsCMD) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate-docs", flag.ContinueOnError)
	fs.StringVar(&cmd.outputDir, "output_dir", "./docs/", "directory in which the files will be generated, note that the directory layout is kept")
	fs.StringVar(&cmd.extension, "extension", ".mdx", "extension for generated files")
	return fs
}

func (cmd *GenerateDocsCMD) Run(args []string) int {
	return 0
}

func (cmd *GenerateDocsCMD) Synopsis() string {
	return ""
}
