package env

import (
	"fmt"
	"io"
	"sort"
)

// PrintDiff writes a human-readable summary of a DiffResult to w.
// Output uses +/~/- prefixes consistent with common diff conventions.
func PrintDiff(w io.Writer, d DiffResult) {
	if d.IsEmpty() {
		fmt.Fprintln(w, "No changes.")
		return
	}

	if len(d.Added) > 0 {
		keys := sortedKeys(d.Added)
		for _, k := range keys {
			fmt.Fprintf(w, "+ %s\n", k)
		}
	}

	if len(d.Changed) > 0 {
		keys := sortedKeys(d.Changed)
		for _, k := range keys {
			fmt.Fprintf(w, "~ %s\n", k)
		}
	}

	if len(d.Removed) > 0 {
		keys := sortedKeys(d.Removed)
		for _, k := range keys {
			fmt.Fprintf(w, "- %s\n", k)
		}
	}

	fmt.Fprintf(w, "\nSummary: %d added, %d changed, %d removed\n",
		len(d.Added), len(d.Changed), len(d.Removed))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
