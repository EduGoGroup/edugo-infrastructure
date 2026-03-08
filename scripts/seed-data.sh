#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR/postgres"
go run ./cmd/seed all

cd "$ROOT_DIR/mongodb"
go run ./cmd/seed all
