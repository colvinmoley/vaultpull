package env

import "fmt"

// ParseMergeStrategy converts a string flag value into a MergeStrategy.
// Accepted values: "vault-wins", "local-wins", "add-only".
func ParseMergeStrategy(s string) (MergeStrategy, error) {
	switch s {
	case "vault-wins":
		return StrategyVaultWins, nil
	case "local-wins":
		return StrategyLocalWins, nil
	case "add-only":
		return StrategyAddOnly, nil
	default:
		return StrategyVaultWins, fmt.Errorf(
			"unknown merge strategy %q: must be one of vault-wins, local-wins, add-only", s,
		)
	}
}

// String returns the canonical string representation of a MergeStrategy.
// Implements the fmt.Stringer interface.
func (m MergeStrategy) String() string {
	switch m {
	case StrategyVaultWins:
		return "vault-wins"
	case StrategyLocalWins:
		return "local-wins"
	case StrategyAddOnly:
		return "add-only"
	default:
		return "unknown"
	}
}
