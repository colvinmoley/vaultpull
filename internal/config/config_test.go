package config

import (
	"os"
	"testing"
)

func TestLoad_FlagsOverrideEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://env-vault:8200")
	t.Setenv("VAULT_TOKEN", "env-token")

	cfg, err := Load([]string{
		"--vault-addr", "http://flag-vault:8200",
		"--vault-token", "flag-token",
		"--namespace", "secret/myapp",
		"--output", "secrets.env",
		"--prefix", "MYAPP_",
		"--dry-run",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "VaultAddr", cfg.VaultAddr, "http://flag-vault:8200")
	assertEqual(t, "VaultToken", cfg.VaultToken, "flag-token")
	assertEqual(t, "Namespace", cfg.Namespace, "secret/myapp")
	assertEqual(t, "OutputFile", cfg.OutputFile, "secrets.env")
	assertEqual(t, "Prefix", cfg.Prefix, "MYAPP_")
	if !cfg.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestLoad_FallsBackToEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://env-vault:8200")
	t.Setenv("VAULT_TOKEN", "env-token")

	cfg, err := Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "VaultAddr", cfg.VaultAddr, "http://env-vault:8200")
	assertEqual(t, "VaultToken", cfg.VaultToken, "env-token")
	assertEqual(t, "OutputFile", cfg.OutputFile, ".env") // default
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	t.Setenv("VAULT_TOKEN", "some-token")

	_, err := Load([]string{})
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
}

func TestLoad_MissingVaultToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault:8200")
	os.Unsetenv("VAULT_TOKEN")

	_, err := Load([]string{})
	if err == nil {
		t.Fatal("expected error for missing vault token")
	}
}

func assertEqual(t *testing.T, field, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}
