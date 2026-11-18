#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/contract"
forge build
jq -r '.abi' "$ROOT/artifacts/Counter.sol/Counter.json" >"$ROOT/artifacts/Counter.abi"
jq -r '.bytecode.object' "$ROOT/artifacts/Counter.sol/Counter.json" >"$ROOT/artifacts/Counter.bin"
abigen --abi "$ROOT/artifacts/Counter.abi" \
  --bin "$ROOT/artifacts/Counter.bin" \
  --pkg bindings --type Counter \
  --out "$ROOT/bindings/counter.go"
echo "bindings/counter.go updated"