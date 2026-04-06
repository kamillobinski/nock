package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kamillobinski.dev/nock/internal/config"
)

var root = &cobra.Command{
	Use: "nock",
	Short: "Personal CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init()
	},
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}