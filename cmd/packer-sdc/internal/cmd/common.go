package cmd

import "strings"

// Meta command to be share between all commands
type Meta struct {
}

type ProjectKind int

const (
	UnknownProjectKind ProjectKind = iota
	Packer
	Plugin
)

func (*Meta) ProjectKind(path string) ProjectKind {
	if strings.Contains(path, "github.com/hashicorp/packer/") {
		return Packer
	}
	if strings.Contains(path, "github.com/hashicorp/packer-plugin") {
		return Plugin
	}
	return UnknownProjectKind
}
