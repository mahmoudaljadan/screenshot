#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="$ROOT_DIR/frontend/src"
DIST_DIR="$ROOT_DIR/frontend/dist"

mkdir -p "$DIST_DIR"
cp "$SRC_DIR/index.html" "$DIST_DIR/index.html"
cp "$SRC_DIR/main.js" "$DIST_DIR/main.js"
cp "$SRC_DIR/styles.css" "$DIST_DIR/styles.css"

echo "Synced frontend/src -> frontend/dist"
