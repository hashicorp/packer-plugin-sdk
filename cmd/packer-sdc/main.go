package main

import (
	"log"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/generate"
	se "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/snippet-extractor"
	"github.com/mitchellh/cli"
)

func main() {

	c := cli.NewCLI("packer-sdc", "1.0.0")

	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"struct-markdown": func() (cli.Command, error) {
			cmd := struct_markdown.Command{}
			return &cmd, nil
		},
		"mapstructure-to-hcl2": func() (cli.Command, error) {
			return &mapstructure_to_hcl2.CMD{}, nil
		},
		"generate-docs": func() (cli.Command, error) {
			return &generate.GenerateDocsCMD{}, nil
		},
		"snippet-extractor": func() (cli.Command, error) {
			return &se.SnippetExtractorCMD{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
