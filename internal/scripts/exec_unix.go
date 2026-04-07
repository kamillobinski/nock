//go:build !windows

package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(path string) error {
	var c *exec.Cmd
	switch strings.ToLower(filepath.Ext(path)) {
	case ".sh", ".bash":
		c = exec.Command("bash", path)
	case ".zsh":
		c = exec.Command("zsh", path)
	case ".py":
		c = exec.Command("python3", path)
	default:
		return fmt.Errorf("unsupported script type: %s", filepath.Ext(path))
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
