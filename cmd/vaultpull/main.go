// Package main is the entry point for the vaultpull CLI tool.
// It wires together configuration loading, Vault client initialization,
// secret fetching, namespace filtering, and .env file writing.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/vault"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Define CLI flags.
	var (
		vaultAddr  = flag.String("vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
		vaultToken = flag.String("vault-token", "", "Vault token (overrides VAULT_TOKEN)")
		secretPath = flag.String("path", "", "KV secret path to read from (e.g. secret/myapp)")
		namespace  = flag.String("namespace", "", "Comma-separated namespace prefixes to include (e.g. APP_,DB_)")
		outputFile = flag.String("output", ".env", "Output .env file path")
		showVer    = flag.Bool("version", false, "Print version and exit")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("vaultpull %s\n", version)
		return nil
	}

	if *secretPath == "" {
		return fmt.Errorf("--path is required (e.g. --path secret/myapp)")
	}

	// Load configuration, merging flags and environment variables.
	cfg, err := config.Load(*vaultAddr, *vaultToken)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Build the Vault client.
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	// Fetch secrets from the given path.
	secrets, err := client.ReadSecret(*secretPath)
	if err != nil {
		return fmt.Errorf("reading secret at %q: %w", *secretPath, err)
	}

	// Apply namespace filtering when prefixes are provided.
	var filtered map[string]string
	if *namespace != "" {
		prefixes := splitAndTrim(*namespace, ",")
		filtered = vault.FilterByNamespace(secrets, prefixes)
	} else {
		filtered = secrets
	}

	if len(filtered) == 0 {
		fmt.Fprintln(os.Stderr, "warning: no secrets matched the given criteria")
	}

	// Write the filtered secrets to the output .env file.
	if err := writeEnvFile(*outputFile, filtered); err != nil {
		return fmt.Errorf("writing env file %q: %w", *outputFile, err)
	}

	fmt.Printf("wrote %d secret(s) to %s\n", len(filtered), *outputFile)
	return nil
}

// writeEnvFile serialises a map of key/value pairs into a .env-formatted file.
// Each line is written as KEY=VALUE. Existing file content is replaced.
func writeEnvFile(path string, secrets map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range secrets {
		// Quote values that contain spaces or special characters.
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

// splitAndTrim splits s by sep and trims whitespace from each element.
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
