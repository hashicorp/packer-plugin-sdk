## `packer-sdc fix`

Fix rewrites parts of the plugin codebase to address known issues or common workarounds used within plugins consuming the Packer plugin SDK.

Options:

 -diff      If the -diff flag is set, no files are rewritten. Instead, fix prints the differences a rewrite would introduce.

Available Fixes:

 gocty      Adds a replace directive for github.com/zclconf/go-cty to github.com/nywilken/go-cty


### Related Issues
Use `packer-sdc fix` to resolve the [cty.Value does not implement gob.GobEncoder](https://github.com/hashicorp/packer-plugin-sdk/issues/187)

