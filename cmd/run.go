package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kamillobinski.dev/nock/internal/scripts"
)

var listFlag bool

var runCmd = &cobra.Command{
	Use:   "run [script]",
	Short: "Run a script",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptsDir, err := scriptsDir()
		if err != nil {
			return err
		}

		if listFlag {
			list, err := scripts.List(scriptsDir)
			if err != nil {
				return err
			}
			for _, s := range list {
				fmt.Printf("%-20s %s\n", s.Name, s.Description)
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("no script provided")
		}

		list, err := scripts.List(scriptsDir)
		if err != nil {
			return err
		}

		var target *scripts.Script
		for i, s := range list {
			if s.Name == args[0] {
				target = &list[i]
				break
			}
		}

		if target == nil {
			return fmt.Errorf("unknown script: %s", args[0])
		}

		return scripts.Run(target.Path)
	},
}

func scriptsDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nock", "scripts"), nil
}

func init() {
	runCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available scripts")
	root.AddCommand(runCmd)
}
