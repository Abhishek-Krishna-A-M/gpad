#!/usr/bin/env bash
set -e

APP="gpad"
OUTPUT_DIR="dist"

mkdir -p $OUTPUT_DIR

# Build matrix
targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
  "windows arm64"
)

echo "Building binaries..."

for t in "${targets[@]}"; do
  os=$(echo $t | awk '{print $1}')
  arch=$(echo $t | awk '{print $2}')
  
  out="$OUTPUT_DIR/${APP}-${os}-${arch}"
  if [ "$os" = "windows" ]; then
    out="$out.exe"
  fi

  echo "→ $os/$arch"

  GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o "$out" ./cmd/gpad
done

echo "Done. Binaries are in ./dist/"

