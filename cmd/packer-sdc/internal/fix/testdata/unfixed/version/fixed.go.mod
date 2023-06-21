module packer-plugin-scaffolding

go 1.20

require (
	github.com/hashicorp/hcl/v2 v2.13.0
	github.com/hashicorp/packer-plugin-sdk v0.3.1
	github.com/zclconf/go-cty v1.10.0
)

replace github.com/zclconf/go-cty => github.com/nywilken/go-cty v1.13.1 // add by packer-sdc fix as noted in github.com/hashicorp/packer-plugin-sdk/issues/187
