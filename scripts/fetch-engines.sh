#!/usr/bin/env bash
# Downloads the Windows versions of the four converter engines into
# engines/windows/, ready to be put in the zip by scripts/build.sh.
#
#   ffmpeg.exe, ffprobe.exe   video, audio, images      (gyan.dev "essentials" build)
#   mutool.exe                PDF pages and text         (MuPDF)
#   gswin64c.exe, gsdll64.dll PDF size reduction          (Ghostscript)
#   exiftool.exe + exiftool_files/   photo metadata      (ExifTool)
#
# Safe to run again: engines already present are skipped.
# Needs: curl, unzip, 7z (package p7zip-full, to unpack the Ghostscript installer).
set -euo pipefail

FFMPEG_URL="https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
MUPDF_VERSION="1.24.10"
MUPDF_URL="https://mupdf.com/downloads/archive/mupdf-${MUPDF_VERSION}-windows.zip"
GS_VERSION="10040"   # Ghostscript 10.04.0
GS_URL="https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs${GS_VERSION}/gs${GS_VERSION}w64.exe"
# exiftool.org only hosts the newest version; ver.txt says which one that is.
EXIFTOOL_VERSION_URL="https://exiftool.org/ver.txt"

cd "$(dirname "$0")/.."
OUT="engines/windows"
TMP="engines/download"
mkdir -p "$OUT" "$TMP"

download() { # url file
  echo "Downloading $1"
  curl -fL --retry 3 -o "$2" "$1"
}

if [ ! -f "$OUT/ffmpeg.exe" ]; then
  download "$FFMPEG_URL" "$TMP/ffmpeg.zip"
  rm -rf "$TMP/ffmpeg" && unzip -q "$TMP/ffmpeg.zip" -d "$TMP/ffmpeg"
  cp "$TMP"/ffmpeg/*/bin/ffmpeg.exe "$TMP"/ffmpeg/*/bin/ffprobe.exe "$OUT/"
fi

if [ ! -f "$OUT/mutool.exe" ]; then
  download "$MUPDF_URL" "$TMP/mupdf.zip"
  rm -rf "$TMP/mupdf" && unzip -q "$TMP/mupdf.zip" -d "$TMP/mupdf"
  cp "$TMP"/mupdf/*/mutool.exe "$OUT/"
fi

if [ ! -f "$OUT/gswin64c.exe" ]; then
  download "$GS_URL" "$TMP/gs.exe"
  rm -rf "$TMP/gs" && 7z x -y -o"$TMP/gs" "$TMP/gs.exe" >/dev/null
  cp "$TMP/gs/bin/gswin64c.exe" "$TMP/gs/bin/gsdll64.dll" "$OUT/"
fi

if [ ! -f "$OUT/exiftool.exe" ]; then
  version="$(curl -fsSL "$EXIFTOOL_VERSION_URL" | tr -d '[:space:]')"
  download "https://exiftool.org/exiftool-${version}_64.zip" "$TMP/exiftool.zip"
  rm -rf "$TMP/exiftool" && unzip -q "$TMP/exiftool.zip" -d "$TMP/exiftool"
  cp "$TMP"/exiftool/*/"exiftool(-k).exe" "$OUT/exiftool.exe"
  rm -rf "$OUT/exiftool_files" && cp -r "$TMP"/exiftool/*/exiftool_files "$OUT/"
fi

# Check everything is there.
for f in ffmpeg.exe ffprobe.exe mutool.exe gswin64c.exe gsdll64.dll exiftool.exe exiftool_files; do
  [ -e "$OUT/$f" ] || { echo "MISSING: $OUT/$f" >&2; exit 1; }
done
rm -rf "$TMP"
echo "All engines ready in $OUT"
