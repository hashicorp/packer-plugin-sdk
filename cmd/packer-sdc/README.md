
# packer-sdc

The packer software development command is meant for plugin maintainers and
Packer maintainers, it helps generate docs and code. For Packer and for Packer
Plugins.

At the top of many Packer go files, you can see the following:

```go
//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,CustomerEncryptionKey
```
This will generate multiple files.

## `struct-markdown`

`struct-markdown` will read all go structs that have fields with a mapstructure
tag, for Example:


```
//go:generate packer-sdc struct-markdown

// Config helps configure things
type Config struct {

	// The JSON file containing your account credentials. Not required if you
	// run Packer on a GCE instance with a service account. Instructions for
	// creating the file or using service accounts are above.
	AccountFile string `mapstructure:"account_file" required:"false"`

    // Foo is an example field.
    Foo string `mapstructure:"foo" required:"true"`
```

This will generate a `Config.mdx` file containing the header docs of the Config
struct, a `Config-not-required.mdx` file containing the docs for the the
`account_file` field and a `Config-required.mdx` file containing the docs for
the `foo` field. This is quite helpful in the sense that the code now becomes
the single source of truth for docs. In Packer, many common structs are reused
in internal and external plugins, this binary makes it possible to use these the
docs too. See the documentation for `renderdocs` for further info.

## `mapstructure-to-hcl`

`mapstructure-to-hcl` Helps generate the code necessary for any plugin,
including the core plugins to communicate the HCL2 layout of a plugin with
Packer, for Example, with the following `config.go` file:


```
go:generate packer-sdc mapstructure-to-hcl2 -type Config

// Config helps configure things
type Config struct {
```

Will generate a `config.hcl2spec.go` file containing a `FlatConfig` struct,
which is a struct with all the mapstructure nested fields 'flattened' into the
`FlatConfig`, so nothing is nested. The `FlatConfig` struct will get a
`HCL2Spec` function that describes its HCL2 layout. This will be used to read
and validate actual HCL2 files. The `config.hcl2spec.go` will also add a
`FlatMapstructure` function to the `Config` struct. That function returns a
`FlatConfig`. These functions together define an interface meant for a plugin
component to 'speak' the HCL2 language with the Packer core.

Before HCL2, Packer JSON heavily relied on the mapstructure decoding library to
load/parse user config files, making this part of the code very tested. To go to
HCL2 this command was created. 

Here are a few differences/gaps betweens hcl2 and mapstructure:
 * in HCL2 all basic struct fields (string/int/struct) that are not pointers
  are required ( must be set ). In mapstructure everything is optional.
 * mapstructure allows to 'squash' fields
 (ex: Field CommonStructType `mapstructure:",squash"`) this allows to
 decorate structs and reuse configuration code. HCL2 parsing libs don't have
 anything similar.
`mapstructure-to-hcl2` parses go files for a plugin and generates the HCL2
compliant code that allows to not change in
order to softly move to HCL2

## `renderdocs`

`renderdocs` is meant to be used in Packer plugins. It renders the docs by
replacing any `@include 'partial.mdx'` call with its actual content, for
example:

`packer-sdc renderdocs -src ./docs-src -dst ./docs-rendered -partials ./docs-partials`

Will first copy the contents of the `./docs-src` dir into the `./docs-rendered`
dir (any file in `./docs-rendered` that is not present it `./docs-src` will be
removed), then each file in `./docs-rendered` will be parsed and any 
`@include 'partial.mdx'` call will be replaced by the contents of the 
`partial.mdx` file.
