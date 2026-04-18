package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const timeFormat = "15:04"
const dateFormat = "2006-01-02"

type Entry struct {
	Time    time.Time
	Message string
}

func Dir() string {
	dir := viper.GetString("log.dir")
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, "Documents", "logs")
	}
	return expandPath(dir)
}

func TodayFile() string {
	return filepath.Join(Dir(), time.Now().Format(dateFormat)+".md")
}

func FileForDate(date string) string {
	return filepath.Join(Dir(), date+".md")
}

func Append(message string) error {
	if err := os.MkdirAll(Dir(), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(TodayFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(time.Now().Format(timeFormat) + "  " + message + "\n")
	return err
}

func ListFiles() ([]string, error) {
	entries, err := os.ReadDir(Dir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var dates []string
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		if name != "" {
			dates = append(dates, name)
		}
	}
	return dates, nil
}

func ReadFormatted(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	entries := parse(string(data))
	if len(entries) == 0 {
		return nil, nil
	}

	var lines []string
	for i, e := range entries {
		var tag string
		if i < len(entries)-1 {
			tag = e.Time.Format(timeFormat) + " (" + formatDuration(entries[i+1].Time.Sub(e.Time)) + ")"
		} else {
			tag = e.Time.Format(timeFormat) + " (ongoing)"
		}
		lines = append(lines, fmt.Sprintf("%-20s %s", tag, e.Message))
	}
	return lines, nil
}

func parse(content string) []Entry {
	var entries []Entry
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			continue
		}
		t, err := time.Parse(timeFormat, strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		entries = append(entries, Entry{Time: t, Message: strings.TrimSpace(parts[1])})
	}
	return entries
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
