package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadFile parses an existing .env file and returns its key-value pairs.
// Lines starting with '#' and blank lines are ignored.
// Returns an empty map if the file does not exist.
func ReadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("env: open file %q: %w", path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("env: %q line %d: missing '='" , path, lineNum)
		}
		k := strings.TrimSpace(line[:idx])
		v := unquoteValue(strings.TrimSpace(line[idx+1:]))
		result[k] = v
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %q: %w", path, err)
	}
	return result, nil
}

// unquoteValue strips surrounding double quotes and unescapes internal quotes.
func unquoteValue(v string) string {
	if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
		v = v[1 : len(v)-1]
		v = strings.ReplaceAll(v, `\"`, `"`)
	}
	return v
}
