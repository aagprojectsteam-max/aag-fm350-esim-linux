#!/usr/bin/env bash
set -Eeuo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$ROOT/bin"
for d in "$ROOT"/cmd/*; do
  [ -f "$d/main.go" ] || continue
  name="$(basename "$d")"
  echo "==> building $name"
  (cd "$ROOT" && go build -trimpath -o "$ROOT/bin/aag-esim-$name" "./cmd/$name")
done
