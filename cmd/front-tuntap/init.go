package front

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "front-tuntap",
	Aliases: []string{"tuntap"},
	Short:   "Using tuntap device for user communicating.",

	RunE: main,
}

var (
	device  = "/dev/tun0"
	backend = "/var/run/jormungand.sock"
)

func init() {

	flag := Command.Flags()
	flag.StringVarP(&device, "device", "d", device,
		"tuntap device")
	flag.StringVarP(&backend, "backend", "B", backend,
		"backend listening address")
}
