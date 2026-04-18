package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	nocklog "kamillobinski.dev/nock/internal/log"
)

var logListFlag bool
var logShowDateFlag string

var logCmd = &cobra.Command{
	Use:     "log [message]",
	Aliases: []string{"l"},
	Short:   "Log your work",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if logListFlag {
			dates, err := nocklog.ListFiles()
			if err != nil {
				return err
			}
			for _, d := range dates {
				fmt.Println(d)
			}
			return nil
		}

		if len(args) > 0 {
			return nocklog.Append(args[0])
		}

		path := nocklog.TodayFile()
		if err := os.MkdirAll(nocklog.Dir(), 0755); err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		f.Close()

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

var logShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a work log",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if logShowDateFlag != "" {
			path = nocklog.FileForDate(logShowDateFlag)
		} else {
			path = nocklog.TodayFile()
		}

		lines, err := nocklog.ReadFormatted(path)
		if err != nil {
			return err
		}
		for _, l := range lines {
			fmt.Println(l)
		}
		return nil
	},
}

func init() {
	logCmd.Flags().BoolVarP(&logListFlag, "list", "l", false, "List log files")
	logShowCmd.Flags().StringVarP(&logShowDateFlag, "date", "d", "", "Date to show (YYYY-MM-DD)")
	logCmd.AddCommand(logShowCmd)
	root.AddCommand(logCmd)
}
