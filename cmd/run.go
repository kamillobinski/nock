package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

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
			printScripts(list)
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("no script provided")
		}

		list, err := scripts.List(scriptsDir)
		if err != nil {
			return err
		}

		var matches []scripts.Script
		for _, s := range list {
			if s.Name == args[0] {
				matches = append(matches, s)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("unknown script: %s", args[0])
		}
		if len(matches) > 1 {
			paths := make([]string, len(matches))
			for i, s := range matches {
				paths[i], _ = filepath.Rel(scriptsDir, s.Path)
			}
			return fmt.Errorf("script name %s is used more than once: %s", args[0], strings.Join(paths, ", "))
		}
		target := matches[0]

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
		if strings.Contains(strings.ToLower(s.Name), q) ||
			strings.Contains(strings.ToLower(s.Description), q) ||
			strings.Contains(strings.ToLower(s.Group), q) {
			out = append(out, s)
		}
	}
	return out
}

func printScripts(list []scripts.Script) {
	tty := isTerminal(os.Stdout)
	if len(list) == 0 {
		if isTerminal(os.Stderr) {
			fmt.Fprintln(os.Stderr, "\033[2mno scripts found\033[0m")
		} else {
			fmt.Fprintln(os.Stderr, "no scripts found")
		}
		return
	}

	width := 20
	for _, s := range list {
		n := utf8.RuneCountInString(s.Name)
		if s.Group != "" {
			n += 2
		}
		if n > width {
			width = n
		}
	}

	style := func(code, s string) string {
		if !tty || s == "" {
			return s
		}
		return "\033[" + code + "m" + s + "\033[0m"
	}

	group := ""
	for i, s := range list {
		if s.Group != group {
			if i > 0 {
				fmt.Println()
			}
			fmt.Println(style("1", s.Group))
			group = s.Group
		}
		if s.Group == "" {
			fmt.Printf("%-*s %s\n", width, s.Name, style("2", s.Description))
		} else {
			fmt.Printf("  %-*s %s\n", width-2, s.Name, style("2", s.Description))
		}
	}
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
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
