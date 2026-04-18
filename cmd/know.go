package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"kamillobinski.dev/nock/internal/knowledge"
)

var knowListFlag bool

var knowCmd = &cobra.Command{
	Use:   "know",
	Short: "Manage your knowledge base",
	RunE: func(cmd *cobra.Command, args []string) error {
		if knowListFlag {
			notes, err := knowledge.ListNotes()
			if err != nil {
				return err
			}
			for _, n := range notes {
				fmt.Printf("%-20s %s\n", n.Name, n.Description)
			}
			return nil
		}
		return cmd.Help()
	},
}

var knowDraftCmd = &cobra.Command{
	Use:   "draft [name]",
	Short: "Create a knowledge draft",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) > 0 {
			name = args[0]
		} else {
			fmt.Print("draft name: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			name = strings.TrimSpace(input)
			if name == "" {
				return fmt.Errorf("draft name cannot be empty")
			}
		}

		name = knowledge.NormalizeName(name)
		path, err := knowledge.CreateDraft(name)
		if err != nil {
			return err
		}

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

func init() {
	knowCmd.Flags().BoolVarP(&knowListFlag, "list", "l", false, "List knowledge notes")
	knowCmd.AddCommand(knowDraftCmd)
	root.AddCommand(knowCmd)
}
