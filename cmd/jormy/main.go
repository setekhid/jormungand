package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewJormyCommand() *cobra.Command {

	command := &cobra.Command{
		Use: "jormy",
	}

	command.AddCommand(
		NewClientCommand(),
		NewServerCommand(),
	)

	return command
}

func main() {
	if err := NewJormyCommand().Execute(); err != nil {
		log.Fatalln(err)
	}
}
