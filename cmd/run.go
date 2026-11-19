package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"kamillobinski.dev/nock/internal/scripts"
)

var listFlag bool

var runCmd = &cobra.Command{
	Use:     "run [script] [args...]",
	Aliases: []string{"r"},
	Short:   "Run a script",
	Args:    cobra.ArbitraryArgs,
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
			if len(args) > 0 {
				list = filterScripts(list, strings.Join(args, " "))
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

		err = scripts.Run(target.Path, args[1:]...)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// The script reported its own error; just pass its exit code through.
			os.Exit(exitErr.ExitCode())
		}
		return err
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 || listFlag {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		dir, err := scriptsDir()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		list, err := scripts.List(dir)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var out []string
		for _, s := range list {
			if !strings.HasPrefix(s.Name, toComplete) {
				continue
			}
			if s.Description != "" {
				out = append(out, s.Name+"\t"+s.Description)
			} else {
				out = append(out, s.Name)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	},
}

func filterScripts(list []scripts.Script, query string) []scripts.Script {
	q := strings.ToLower(query)
	var out []scripts.Script
	for _, s := range list {
		if strings.Contains(strings.ToLower(s.Name), q) || strings.Contains(strings.ToLower(s.Description), q) {
			out = append(out, s)
		}
	}
	return out
}

func scriptsDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "nock", "scripts"), nil
}

func init() {
	runCmd.Flags().SetInterspersed(false)
	runCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available scripts, optionally filtered: -l [filter]")
	root.AddCommand(runCmd)
}
