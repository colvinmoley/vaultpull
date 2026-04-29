package env

// MergeStrategy defines how to handle conflicts when merging secrets.
type MergeStrategy int

const (
	// StrategyVaultWins overwrites local values with Vault values.
	StrategyVaultWins MergeStrategy = iota
	// StrategyLocalWins keeps local values when a key exists in both.
	StrategyLocalWins
	// StrategyAddOnly only adds keys that are missing locally; never overwrites.
	StrategyAddOnly
)

// Merge combines existing local secrets with incoming Vault secrets
// according to the given MergeStrategy. It returns a new map that
// can be written back to the .env file.
func Merge(local, vault map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(local))

	// Copy all local keys first.
	for k, v := range local {
		result[k] = v
	}

	switch strategy {
	case StrategyVaultWins:
		// Vault values overwrite everything.
		for k, v := range vault {
			result[k] = v
		}

	case StrategyLocalWins:
		// Only add keys that don't already exist locally.
		for k, v := range vault {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}

	case StrategyAddOnly:
		// Same as LocalWins — never overwrite existing local values.
		for k, v := range vault {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}

	return result
}
