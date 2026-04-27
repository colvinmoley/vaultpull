package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is a minimal HashiCorp Vault HTTP client.
type Client struct {
	addr       string
	token      string
	httpClient *http.Client
}

// SecretData holds the key/value pairs returned from a Vault KV secret.
type SecretData map[string]string

// NewClient creates a new Vault client with the given address and token.
func NewClient(addr, token string) *Client {
	return &Client{
		addr:  strings.TrimRight(addr, "/"),
		token: token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ListSecrets returns the keys available under the given KV v2 path.
func (c *Client) ListSecrets(path string) ([]string, error) {
	url := fmt.Sprintf("%s/v1/%s?list=true", c.addr, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: build list request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: list request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: list %s returned %d", path, resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("vault: decode list response: %w", err)
	}
	return result.Data.Keys, nil
}

// ReadSecret fetches the key/value data for a KV v2 secret at path.
func (c *Client) ReadSecret(path string) (SecretData, error) {
	url := fmt.Sprintf("%s/v1/%s", c.addr, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: build read request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: read request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: read %s returned %d", path, resp.StatusCode)
	}

	var result struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("vault: decode read response: %w", err)
	}

	secrets := make(SecretData, len(result.Data.Data))
	for k, v := range result.Data.Data {
		secrets[k] = fmt.Sprintf("%v", v)
	}
	return secrets, nil
}
