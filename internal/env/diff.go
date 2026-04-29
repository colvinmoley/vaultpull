package env

// DiffResult holds the changes between two secret maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string // key -> new value
}

// Diff compares an existing (local) map with an incoming (remote) map and
// returns a DiffResult describing what would change if the incoming map were
// written to disk.
func Diff(existing, incoming map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]string),
	}

	for k, v := range incoming {
		if existing == nil {
			result.Added[k] = v
			continue
		}
		if old, ok := existing[k]; !ok {
			result.Added[k] = v
		} else if old != v {
			result.Changed[k] = v
		}
	}

	for k := range existing {
		if _, ok := incoming[k]; !ok {
			result.Removed[k] = existing[k]
		}
	}

	return result
}

// IsEmpty returns true when there are no additions, removals, or changes.
func (d DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}
