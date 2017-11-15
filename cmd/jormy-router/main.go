package main

import (
	flag "github.com/spf13/pflag"
)

var (
	ip_info = "/etc/jormungand/ip_info.db"
)

func init() {

	flag.StringVarP(&ip_info, "ip-info", "B", ip_info,
		"ip information database file")
}

func main() {
	// TODO
}
