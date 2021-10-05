//go:generate packer-sdc struct-markdown

package communicator

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// Configure screenshot capture
type ScreenshotConfig struct {
	// The directory in which screenshots will be placed
	// Defaults to ./screenshots/.
	ScreenshotDirectory string `mapstructure:"screenshot_dir"`

	// Skip taking a screenshot on failure.
	// Defaults to false.
	ScreenshotSkip bool `mapstructure:"skip_screenshot"`
}

// The ConfigSpec funcs are used by the Packer core to parse HCL2 templates.
func (c *ScreenshotConfig) ConfigSpec() hcldec.ObjectSpec { return c.FlatMapstructure().HCL2Spec() }

// Configure parses the json template into the Config structs
func (c *ScreenshotConfig) Configure(raws ...interface{}) ([]string, error) {
	err := config.Decode(c, nil, raws...)
	return nil, err
}

func (c *ScreenshotConfig) Prepare(ctx *interpolate.Context) []error {
	if c.ScreenshotDirectory == "" {
		c.ScreenshotDirectory = "./screenshots/"
	}

	// Validation
	var errs []error

	// Not really anything to validate here yet.

	return errs
}