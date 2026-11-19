package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Install shell integration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var configFile, completion string
		switch args[0] {
		case "zsh":
			configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
			completion = zshCompletion
		case "bash":
			configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
			completion = bashCompletion
		default:
			return fmt.Errorf("unsupported shell: %s (supported: zsh, bash)", args[0])
		}

		current, err := os.ReadFile(configFile)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		block := shellIntegrationStart + shellFunc + completion + shellIntegrationEnd
		content := string(current)

		start := strings.Index(content, shellIntegrationStart)
		end := strings.Index(content, shellIntegrationEnd)

		if start != -1 && end != -1 {
			content = content[:start] + block + content[end+len(shellIntegrationEnd):]
		} else {
			content = content + "\n" + block
		}

		if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
			return err
		}

		fmt.Printf("nock shell integration written to %s\n", configFile)
		fmt.Printf("run: source %s\n", configFile)
		return nil
	},
}

const shellIntegrationStart = "# nock:start\n"
const shellIntegrationEnd = "# nock:end\n"

const shellFunc = `nock() {
  if [ "$1" = "nav" ] && [ "$2" = "back" ]; then
    cd "$NOCK_PREV_DIR"
  elif [ "$1" = "nav" ] && [ -n "$2" ] && [[ "$2" != -* ]]; then
    local dir
    dir=$(command nock "$@")
    if [ $? -eq 0 ]; then
      export NOCK_PREV_DIR="$PWD"
      cd "$dir"
    fi
  elif [ "$1" = "run" ] || [ "$1" = "r" ]; then
    local cd_file dir ret
    cd_file=$(mktemp)
    NOCK_CD_FILE="$cd_file" command nock "$@"
    ret=$?
    dir=$(cat "$cd_file")
    rm -f "$cd_file"
    if [ $ret -eq 0 ] && [ -n "$dir" ]; then
      export NOCK_PREV_DIR="$PWD"
      cd "$dir"
    fi
    return $ret
  else
    command nock "$@"
  fi
}
`

const zshCompletion = `if ! (( $+functions[compdef] )); then
  autoload -Uz compinit && compinit
fi
eval "$(command nock completion zsh)"
compdef _nock nock
`

const bashCompletion = `eval "$(command nock completion bash)"
`

func init() {
	root.AddCommand(initCmd)
}
