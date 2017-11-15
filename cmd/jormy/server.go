package main

import (
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {

	command := &cobra.Command{
		Use:  "server",
		RunE: serverMain,
	}

	// TODO

	return command
}

func serverMain(cmd *cobra.Command, args []string) error {
	// TODO
	return nil
}
