package main

import (
	flag "github.com/spf13/pflag"
)

var (
	addr    = ":8989"
	backend = "/var/run/jormungand.sock"
)

func init() {

	flag.StringVarP(&addr, "addr", "A", addr,
		"frontend serving address")
	flag.StringVarP(&backend, "backend", "B", backend,
		"backend listening address")
}

func main() {
	// TODO
}
