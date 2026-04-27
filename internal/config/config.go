package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Namespace  string
	OutputFile string
	Prefix     string
	DryRun     bool
}

// Load parses flags and environment variables into a Config.
// Flags take precedence over environment variables.
func Load(args []string) (*Config, error) {
	fs := pflag.NewFlagSet("vaultpull", pflag.ContinueOnError)

	vaultAddr := fs.String("vault-addr", "", "Vault server address (env: VAULT_ADDR)")
	vaultToken := fs.String("vault-token", "", "Vault token (env: VAULT_TOKEN)")
	namespace := fs.String("namespace", "", "Vault namespace / secret path prefix to filter")
	outputFile := fs.String("output", ".env", "Output .env file path")
	prefix := fs.String("prefix", "", "Optional key prefix to strip from secret keys")
	dryRun := fs.Bool("dry-run", false, "Print secrets without writing to file")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	cfg := &Config{
		VaultAddr:  resolveString(*vaultAddr, "VAULT_ADDR"),
		VaultToken: resolveString(*vaultToken, "VAULT_TOKEN"),
		Namespace:  *namespace,
		OutputFile: *outputFile,
		Prefix:     *prefix,
		DryRun:     *dryRun,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func resolveString(flagVal, envKey string) string {
	if flagVal != "" {
		return flagVal
	}
	return os.Getenv(envKey)
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("vault address is required (--vault-addr or VAULT_ADDR)")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("vault token is required (--vault-token or VAULT_TOKEN)")
	}
	return nil
}
