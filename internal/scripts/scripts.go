package scripts

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Script struct {
	Name        string
	Group       string
	Description string
	Path        string
}

func List(dir string) ([]Script, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var root, grouped []Script
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.IsDir() {
			group, err := readDir(filepath.Join(dir, entry.Name()), entry.Name())
			if err != nil {
				return nil, err
			}
			grouped = append(grouped, group...)
			continue
		}
		root = append(root, newScript(dir, entry.Name(), ""))
	}

	return append(root, grouped...), nil
}

func readDir(dir, group string) ([]Script, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var scripts []Script
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		scripts = append(scripts, newScript(dir, entry.Name(), group))
	}
	return scripts, nil
}

func newScript(dir, file, group string) Script {
	path := filepath.Join(dir, file)
	return Script{
		Name:        strings.TrimSuffix(file, filepath.Ext(file)),
		Group:       group,
		Description: readDescription(path),
		Path:        path,
	}
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
