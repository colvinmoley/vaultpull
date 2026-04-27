package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpull/internal/vault"
)

func newListResponse(keys []string) []byte {
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{"keys": keys},
	})
	return body
}

func newReadResponse(data map[string]interface{}) []byte {
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{"data": data},
	})
	return body
}

func TestListSecrets_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(newListResponse([]string{"db", "api"}))
	}))
	defer server.Close()

	c := vault.NewClient(server.URL, "test-token")
	keys, err := c.ListSecrets("secret/metadata/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 || keys[0] != "db" || keys[1] != "api" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestListSecrets_NonOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	c := vault.NewClient(server.URL, "bad-token")
	_, err := c.ListSecrets("secret/metadata/myapp")
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func TestReadSecret_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(newReadResponse(map[string]interface{}{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		}))
	}))
	defer server.Close()

	c := vault.NewClient(server.URL, "test-token")
	data, err := c.ReadSecret("secret/data/myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["DB_HOST"] != "localhost" {
		t.Fatalf("expected DB_HOST=localhost, got %q", data["DB_HOST"])
	}
	if data["DB_PORT"] != "5432" {
		t.Fatalf("expected DB_PORT=5432, got %q", data["DB_PORT"])
	}
}

func TestReadSecret_NonOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := vault.NewClient(server.URL, "test-token")
	_, err := c.ReadSecret("secret/data/myapp/missing")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}
