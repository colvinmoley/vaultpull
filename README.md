# vaultpull

> CLI tool to sync HashiCorp Vault secrets to local `.env` files with namespace filtering

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultpull.git
cd vaultpull
go build -o vaultpull .
```

---

## Usage

Set your Vault address and token, then run `vaultpull` against a secret path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull secrets from a path and write to .env
vaultpull pull --path secret/data/myapp --output .env

# Filter by namespace
vaultpull pull --path secret/data/myapp --namespace production --output .env.production
```

The resulting `.env` file will contain key-value pairs sourced directly from Vault:

```
DB_HOST=db.example.com
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Description |
|-------------|--------------------------------------|
| `--path` | Vault secret path to pull from |
| `--output` | Output file path (default: `.env`) |
| `--namespace` | Filter secrets by namespace prefix |
| `--overwrite` | Overwrite existing `.env` file |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)