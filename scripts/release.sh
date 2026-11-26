#!/usr/bin/env bash
# Publishes a new version on GitHub: builds dist/QuickConvert.zip and
# creates the release "v<VERSION>" with the zip attached, so the download
# link in README.md (releases/latest/download/QuickConvert.zip) gets it.
#
#   1. set VERSION in scripts/build.sh
#   2. scripts/release.sh "What changed in this version"
set -euo pipefail
cd "$(dirname "$0")/.."

notes="${1:?usage: scripts/release.sh \"release notes\"}"
version="$(sed -n 's/^VERSION="\(.*\)"/\1/p' scripts/build.sh)"

scripts/build.sh
gh release create "v$version" dist/QuickConvert.zip \
  --title "QuickConvert $version" --notes "$notes" --target master
