package front

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "front-socks5",
	Aliases: []string{"socks5"},
	Short:   "Using Socks5 protocol for user communicating.",

	RunE: main,
}

var (
	addr    = ":8989"
	backend = "/var/run/jormungand.sock"
)

func init() {

	flag := Command.Flags()
	flag.StringVarP(&addr, "addr", "A", addr,
		"frontend serving address")
	flag.StringVarP(&backend, "backend", "B", backend,
		"backend listening address")
}
