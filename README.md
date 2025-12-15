# sma-go

Secure messaging service in Go. Identities are keyed by a human-friendly UniqueString and a PGP public key. Messages reference identities by UniqueString and are stored as OpenPGP ciphertext.

## Features

- Identity model with unique UniqueString and required PGP public key.
- Messages with FK to Identity.UniqueString (referential integrity via GORM).
- OpenPGP encryption via gopenpgp (supports unarmored binary + base64 to avoid ASCII-armor headers).
- Simple .env-driven configuration.

## Tech

- Go, GORM, PostgreSQL
- ProtonMail gopenpgp v3

## Setup

1. Prerequisites

   - Go 1.24+
   - PostgreSQL 17+

2. Configuration

   - Copy `.env.sample` to `.env` and set:
     - `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`, `DB_PORT`

3. Install and run (Windows, PowerShell)
   - `go mod tidy`
   - `go run .`

## Models

- Identity
  - PublicKey: string, not null
  - UniqueString: string, unique, not null (indexed)
- Message
  - FromUnique, ToUnique: string, indexed, not null
  - FKs: `foreignKey:FromUnique/ToUnique` → `references:UniqueString` with cascade update, set null on delete
  - Text: text, not null

## Encryption notes

- To drop PGP armor headers/footers, use binary messages and base64-encode for transport; decrypt with `crypto.NoArmor`.

## Environment

- See `.env.sample` for DB variables.

## Development

- Auto-migrate models with GORM on startup (recommended).
- Run tests: `go test ./...`

## License

- Choose and add a LICENSE file.
