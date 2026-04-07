package scripts

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Script struct {
	Name        string
	Description string
	Path        string
}

func List(dir string) ([]Script, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var scripts []Script
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		description := readDescription(path)

		scripts = append(scripts, Script{
			Name:        name,
			Description: description,
			Path:        path,
		})
	}

	return scripts, nil
}

func readDescription(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimPrefix(line, "#"))
		}
	}

	return ""
}
