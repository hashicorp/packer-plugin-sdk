package main

import (
	"log"
	"os"

	mapstructure_to_hcl2 "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/mapstructure-to-hcl2"
	"github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/renderdocs"
	struct_markdown "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/struct-markdown"
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
		"renderdocs": func() (cli.Command, error) {
			return &renderdocs.Command{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
