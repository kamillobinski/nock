package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"kamillobinski.dev/nock/internal/system"
)

var forceRefresh bool

var systemCmd = &cobra.Command{
	Use:     "system",
	Aliases: []string{"sys"},
	Short:   "Display system info",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		systemDir, err := system.Dir()
		if err != nil {
			return err
		}

		cfg, err := system.LoadConfig(systemDir)
		if err != nil {
			return fmt.Errorf("system config not found at %s: %w", systemDir, err)
		}

		cache := system.LoadCache(systemDir)
		defer cache.Save()

		if cfg.ASCII != "" {
			asciiPath := cfg.ASCII
			if !filepath.IsAbs(asciiPath) {
				asciiPath = filepath.Join(systemDir, asciiPath)
			}
			if data, err := os.ReadFile(asciiPath); err == nil {
				fmt.Println()
				for _, line := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
					fmt.Println(" " + line)
				}
				fmt.Println()
			}
		}

		maxLabel := 0
		for _, g := range cfg.Groups {
			for _, item := range g.Items {
				if len(item.Label) > maxLabel {
					maxLabel = len(item.Label)
				}
			}
		}

		// \033[1;37m (7 bytes) + content + \033[0m (4 bytes) = 11 invisible bytes
		const ansiOverhead = 11
		bold := func(s string) string { return "\033[1;37m" + s + "\033[0m" }
		dim := func(s string) string { return "\033[2m" + s + "\033[0m" }

		// fetch all values concurrently then render once everything is ready
		type groupValues struct {
			group  system.Group
			values []string
		}
		all := make([]groupValues, len(cfg.Groups))
		for i, g := range cfg.Groups {
			all[i] = groupValues{g, make([]string, len(g.Items))}
		}

		var wg sync.WaitGroup
		for i, g := range cfg.Groups {
			for j, item := range g.Items {
				wg.Add(1)
				go func(i, j int, item system.Item) {
					defer wg.Done()
					all[i].values[j] = system.FetchValue(item, systemDir, cache, forceRefresh)
				}(i, j, item)
			}
		}
		wg.Wait()

		// global max left-column value width so right column is always aligned
		leftValW := 0
		for _, gv := range all {
			for j := 0; j < len(gv.values); j += 2 {
				if l := len(gv.values[j]); l > leftValW {
					leftValW = l
				}
			}
		}
		leftValW += 4

		for i, gv := range all {
			if i > 0 {
				fmt.Println()
			}
			if gv.group.ShowName {
				fmt.Printf("  %s\n", dim(gv.group.Name))
			}
			for j := 0; j < len(gv.group.Items); j += 2 {
				left := fmt.Sprintf("  %-*s   %-*s", maxLabel+ansiOverhead, bold(gv.group.Items[j].Label), leftValW, gv.values[j])
				if j+1 < len(gv.group.Items) {
					right := fmt.Sprintf("%-*s   %s", maxLabel+ansiOverhead, bold(gv.group.Items[j+1].Label), gv.values[j+1])
					fmt.Println(left + right)
				} else {
					fmt.Println(left)
				}
			}
		}
		fmt.Println()

		return nil
	},
}

func init() {
	root.AddCommand(systemCmd)
	systemCmd.Flags().BoolVarP(&forceRefresh, "force", "f", false, "Ignore cache and refresh all values")
}
