package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile_ParsesKeyValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := "# comment\nDB_HOST=localhost\nDB_PORT=5432\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", got["DB_HOST"], "localhost")
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT = %q, want %q", got["DB_PORT"], "5432")
	}
}

func TestReadFile_NonExistentReturnsEmpty(t *testing.T) {
	got, err := ReadFile("/nonexistent/.env")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestReadFile_UnquotesValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	content := `KEY="hello world"` + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "hello world" {
		t.Errorf("KEY = %q, want %q", got["KEY"], "hello world")
	}
}

func TestReadFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	original := map[string]string{
		"PLAIN":  "value",
		"SPACED": "hello world",
		"QUOTED": `say "hi"`,
	}

	if err := WriteFile(path, original); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	for k, want := range original {
		if got[k] != want {
			t.Errorf("key %q: got %q, want %q", k, got[k], want)
		}
	}
}
