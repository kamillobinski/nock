package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var navCmd = &cobra.Command{
	Use:   "nav [shortcut]",
	Short: "Navigate to a saved directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			shortcuts := viper.GetStringMapString("nav")
			for name, path := range shortcuts {
				fmt.Printf("%-20s %s\n", name, path)
			}
			return nil
		}

		shortcut := args[0]
		path := viper.GetString("nav." + shortcut)

		if path == "" {
			return fmt.Errorf("unknown shortcut: %s", shortcut)
		}

		path = expandPath(path)
		fmt.Println(path)
		return nil
	},
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func init() {
	root.AddCommand(navCmd)
}