package knowledge

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Dir() string {
	dir := viper.GetString("knowledge.dir")
	if dir == "" {
		cfgDir, _ := os.UserConfigDir()
		dir = filepath.Join(cfgDir, "nock", "knowledge")
	}
	return expandPath(dir)
}

func PendingDir() string {
	return filepath.Join(Dir(), "pending")
}

type Note struct {
	Name        string
	Description string
}

func ListNotes() ([]Note, error) {
	path := filepath.Join(Dir(), "INDEX.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var notes []Note
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- [") {
			continue
		}
		end := strings.Index(line, "]")
		if end == -1 {
			continue
		}
		name := line[3:end]
		if name == "" {
			continue
		}
		var desc string
		parts := strings.SplitN(line[end+1:], " — ", 3)
		if len(parts) == 3 {
			desc = parts[2]
		}
		notes = append(notes, Note{Name: name, Description: desc})
	}
	return notes, nil
}

func CreateDraft(name string) (string, error) {
	if err := os.MkdirAll(PendingDir(), 0755); err != nil {
		return "", err
	}
	path := filepath.Join(PendingDir(), name+".md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			return "", err
		}
	}
	return path, nil
}

func NormalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
