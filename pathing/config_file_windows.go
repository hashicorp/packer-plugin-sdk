// +build windows

package pathing

const (
	defaultConfigFile = "packer.config"
)

func getDefaultConfigDir() string {
	return "packer.d"
}
