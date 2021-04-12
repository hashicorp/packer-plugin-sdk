package main

import (
	"log"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/generate"
	ms "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/mapstructure-to-hcl2"
	se "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/snippet-extractor"
	md "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/struct-markdown"
	"github.com/mitchellh/cli"
)

func main() {
	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,

		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	c := cli.NewCLI("packer-sdc", "1.0.0")

	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"struct-markdown": func() (cli.Command, error) {
			cmd := md.StructMarkdownCMD{
				Ui: ui,
			}
			return &cmd, nil
		},
		"mapstructure-to-hcl2": func() (cli.Command, error) {
			return &ms.MapstructureToHCL2{Ui: ui}, nil
		},
		"generate-docs": func() (cli.Command, error) {
			return &generate.GenerateDocsCMD{Ui: ui}, nil
		},
		"snippet-extractor": func() (cli.Command, error) {
			return &se.SnippetExtractorCMD{Ui: ui}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
