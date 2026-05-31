#!/usr/bin/env bash
# Build the Next.js frontend and refresh the embedded bundle.
# Run this before `go build` whenever the frontend changes.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT/frontend"

echo "→ installing frontend deps"
npm install --no-audit --no-fund

echo "→ building static export"
npm run build

echo "→ refreshing internal/web/webui/dist"
rm -rf "$ROOT/internal/web/webui/dist"
cp -R out "$ROOT/internal/web/webui/dist"

echo "✓ UI bundle embedded. Now run: go build ./..."
