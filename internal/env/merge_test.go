package env

import (
	"testing"
)

func TestMerge_VaultWins(t *testing.T) {
	local := map[string]string{"A": "local_a", "B": "local_b"}
	vault := map[string]string{"A": "vault_a", "C": "vault_c"}

	result := Merge(local, vault, StrategyVaultWins)

	assertMergeEqual(t, result, "A", "vault_a") // vault overwrites local
	assertMergeEqual(t, result, "B", "local_b")  // local-only key preserved
	assertMergeEqual(t, result, "C", "vault_c")  // vault-only key added
}

func TestMerge_LocalWins(t *testing.T) {
	local := map[string]string{"A": "local_a", "B": "local_b"}
	vault := map[string]string{"A": "vault_a", "C": "vault_c"}

	result := Merge(local, vault, StrategyLocalWins)

	assertMergeEqual(t, result, "A", "local_a") // local value kept
	assertMergeEqual(t, result, "B", "local_b") // local-only key preserved
	assertMergeEqual(t, result, "C", "vault_c") // vault-only key added
}

func TestMerge_AddOnly(t *testing.T) {
	local := map[string]string{"A": "local_a"}
	vault := map[string]string{"A": "vault_a", "B": "vault_b"}

	result := Merge(local, vault, StrategyAddOnly)

	assertMergeEqual(t, result, "A", "local_a") // existing local not touched
	assertMergeEqual(t, result, "B", "vault_b") // new key from vault added
}

func TestMerge_EmptyLocal(t *testing.T) {
	local := map[string]string{}
	vault := map[string]string{"X": "1", "Y": "2"}

	result := Merge(local, vault, StrategyVaultWins)

	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}

func TestMerge_EmptyVault(t *testing.T) {
	local := map[string]string{"A": "a", "B": "b"}
	vault := map[string]string{}

	result := Merge(local, vault, StrategyVaultWins)

	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	assertMergeEqual(t, result, "A", "a")
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	local := map[string]string{"A": "original"}
	vault := map[string]string{"A": "changed"}

	_ = Merge(local, vault, StrategyVaultWins)

	if local["A"] != "original" {
		t.Errorf("Merge mutated local map: got %q", local["A"])
	}
}

func assertMergeEqual(t *testing.T, m map[string]string, key, want string) {
	t.Helper()
	got, ok := m[key]
	if !ok {
		t.Errorf("key %q not found in result", key)
		return
	}
	if got != want {
		t.Errorf("key %q: got %q, want %q", key, got, want)
	}
}
