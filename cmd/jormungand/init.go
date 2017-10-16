package main

import (
	socks5 "github.com/setekhid/jormungand/cmd/front-socks5"
	tuntap "github.com/setekhid/jormungand/cmd/front-tuntap"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "jormungand",
	Aliases: []string{"jorm"},
	Short:   "",

	Run: func(_ *cobra.Command, _ []string) {},
}

func init() {

	Command.AddCommand(
		socks5.Command,
		tuntap.Command,
	)
}
