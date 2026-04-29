package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// WriteFile writes the given secrets map to a .env file at the specified path.
// Existing file contents are overwritten. Keys are written in sorted order.
func WriteFile(path string, secrets map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("env: create file %q: %w", path, err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", k, quoteValue(secrets[k]))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("env: write key %q: %w", k, err)
		}
	}
	return nil
}

// quoteValue wraps the value in double quotes if it contains spaces or
// special shell characters, escaping any existing double quotes.
func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n\r#$\"\'\\`") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
