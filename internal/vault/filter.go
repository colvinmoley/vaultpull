package vault

import "strings"

// FilterByNamespace returns only the keys from keys that match the given
// namespace prefix. An empty namespace matches everything.
//
// A key is considered to belong to a namespace when it equals the namespace
// or starts with "<namespace>/".
func FilterByNamespace(keys []string, namespace string) []string {
	if namespace == "" {
		return keys
	}
	prefix := strings.TrimRight(namespace, "/") + "/"
	var matched []string
	for _, k := range keys {
		if k == namespace || strings.HasPrefix(k, prefix) {
			matched = append(matched, k)
		}
	}
	return matched
}

// MergeSecrets merges multiple SecretData maps into one.
// Later values overwrite earlier ones on key collision.
func MergeSecrets(sets ...SecretData) SecretData {
	merged := make(SecretData)
	for _, s := range sets {
		for k, v := range s {
			merged[k] = v
		}
	}
	return merged
}
