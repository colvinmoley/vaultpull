package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFile_SortsKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"ZEBRA": "last",
		"APPLE": "first",
		"MANGO": "middle",
	}

	if err := WriteFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APPLE=") {
		t.Errorf("expected first line to start with APPLE=, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA=") {
		t.Errorf("expected last line to start with ZEBRA=, got %q", lines[2])
	}
}

func TestWriteFile_QuotesSpecialValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"KEY": "hello world",
	}

	if err := WriteFile(path, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), `KEY="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(content))
	}
}

func TestWriteFile_EmptySecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := WriteFile(path, map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}
}

func TestQuoteValue(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"has space", `"has space"`},
		{"has#hash", `"has#hash"`},
		{`has"quote`, `"has\"quote"`},
	}
	for _, c := range cases {
		got := quoteValue(c.input)
		if got != c.expected {
			t.Errorf("quoteValue(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}
