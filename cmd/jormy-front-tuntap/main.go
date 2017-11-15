package main

import (
	flag "github.com/spf13/pflag"
)

var (
	device  = "/dev/tun0"
	backend = "/var/run/jormungand.sock"
)

func init() {

	flag.StringVarP(&device, "device", "d", device,
		"tuntap device")
	flag.StringVarP(&backend, "backend", "B", backend,
		"backend listening address")
}

func main() {
	// TODO
}
