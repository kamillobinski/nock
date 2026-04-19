package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var logSummaryAgentFlag string

var logSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Generate a standup summary from your work logs",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		today := time.Now()
		yesterday := today.AddDate(0, 0, -1)

		todayLines, err := nocklog.ReadFormatted(nocklog.TodayFile())
		if err != nil {
			return err
		}
		yesterdayLines, err := nocklog.ReadFormatted(nocklog.FileForDate(yesterday.Format("2006-01-02")))
		if err != nil {
			return err
		}

		if len(yesterdayLines) == 0 && len(todayLines) == 0 {
			return fmt.Errorf("no log entries found for today or yesterday")
		}

		var prompt strings.Builder
		prompt.WriteString("Format these work logs into a professional standup update with Yesterday, Today, and Blockers sections. Keep it concise. Do not include timestamps.\n\n")

		if len(yesterdayLines) > 0 {
			prompt.WriteString(fmt.Sprintf("Yesterday (%s):\n", yesterday.Format("2006-01-02")))
			for _, l := range yesterdayLines {
				prompt.WriteString("  " + l + "\n")
			}
			prompt.WriteString("\n")
		}
		if len(todayLines) > 0 {
			prompt.WriteString(fmt.Sprintf("Today (%s):\n", today.Format("2006-01-02")))
			for _, l := range todayLines {
				prompt.WriteString("  " + l + "\n")
			}
		}

		agentID := logSummaryAgentFlag
		if agentID == "" {
			agentID = viper.GetString("onyx.agent_id")
		}

		cmdArgs := []string{"ask"}
		if agentID != "" {
			cmdArgs = append(cmdArgs, "--agent-id", agentID)
		}
		cmdArgs = append(cmdArgs, prompt.String())

		fmt.Println("generating summary...")
		out, err := exec.Command("onyx-cli", cmdArgs...).Output()
		if err != nil {
			return fmt.Errorf("onyx-cli: %w", err)
		}

		fmt.Println()
		fmt.Print(string(out))
		return nil
	},
}

func init() {
	logCmd.Flags().BoolVarP(&logListFlag, "list", "l", false, "List log files")
	logShowCmd.Flags().StringVarP(&logShowDateFlag, "date", "d", "", "Date to show (YYYY-MM-DD)")
	logCmd.AddCommand(logShowCmd)
	logSummaryCmd.Flags().StringVarP(&logSummaryAgentFlag, "agent", "a", "", "Onyx agent ID to use (overrides onyx.agent_id in config)")
	logCmd.AddCommand(logSummaryCmd)
	root.AddCommand(logCmd)
}
