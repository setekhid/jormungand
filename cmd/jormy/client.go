package main

import (
	"github.com/spf13/cobra"
)

func NewClientCommand() *cobra.Command {

	command := &cobra.Command{
		Use:  "client",
		RunE: clientMain,
	}

	// TODO

	return command
}

func clientMain(cmd *cobra.Command, args []string) error {
	// TODO
	return nil
}
