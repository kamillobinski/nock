//go:build windows

package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(path string, args ...string) error {
	var c *exec.Cmd
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ps1":
		c = exec.Command("powershell.exe", append([]string{"-File", path}, args...)...)
	case ".bat", ".cmd":
		c = exec.Command("cmd.exe", append([]string{"/C", path}, args...)...)
	default:
		return fmt.Errorf("unsupported script type on Windows: %s", filepath.Ext(path))
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
