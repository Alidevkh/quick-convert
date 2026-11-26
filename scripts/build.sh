#!/usr/bin/env bash
# Builds dist/QuickConvert.zip, ready to copy to Windows.
#
#   1. run all tests
#   2. build quickconvert.exe for Windows (no console window)
#   3. write install.cmd / uninstall.cmd
#   4. add README.txt and the engines (if scripts/fetch-engines.sh got them)
#   5. zip it
set -euo pipefail
cd "$(dirname "$0")/.."

VERSION="0.4.0"
OUT="dist/QuickConvert"

go test ./...

rm -rf dist && mkdir -p "$OUT"
GOOS=windows GOARCH=amd64 go build -trimpath \
  -ldflags "-H windowsgui -s -w -X main.version=$VERSION" \
  -o "$OUT/quickconvert.exe" ./cmd/quickconvert
go run ./cmd/genscripts "$OUT"
sed 's/$/\r/' packaging/README.txt > "$OUT/README.txt"   # Windows line endings for Notepad
sed 's/$/\r/' packaging/selftest.ps1 > "$OUT/selftest.ps1"

if [ -f engines/windows/ffmpeg.exe ]; then
  cp -r engines/windows/. "$OUT/"
else
  echo "NOTE: engines/windows is empty - the zip has no engines (see README.txt, 'Getting the engines yourself')."
fi

(cd dist && zip -qr QuickConvert.zip QuickConvert)
echo "Built dist/QuickConvert.zip"
