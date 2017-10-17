package router

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "auto-router",
	Aliases: []string{"router"},
	Short:   "Using automatically routing China traffic back to China.",

	RunE: main,
}

var (
	ip_info = "/etc/jormungand/ip_info.db"
)

func init() {

	flag := Command.Flags()
	flag.StringVarP(&ip_info, "ip-info", "B", ip_info,
		"ip information database file")
}
