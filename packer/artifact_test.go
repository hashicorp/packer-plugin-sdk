// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package packer

type TestArtifact struct {
	id            string
	state         map[string]interface{}
	destroyCalled bool
}

func (*TestArtifact) BuilderId() string {
	return "bid"
}

func (*TestArtifact) Files() []string {
	return []string{"a", "b"}
}

func (a *TestArtifact) Id() string {
	id := a.id
	if id == "" {
		id = "id"
	}

	return id
}

func (*TestArtifact) String() string {
	return "string"
}

func (a *TestArtifact) State(name string) interface{} {
	value := a.state[name]
	return value
}

func (a *TestArtifact) Destroy() error {
	a.destroyCalled = true
	return nil
}
