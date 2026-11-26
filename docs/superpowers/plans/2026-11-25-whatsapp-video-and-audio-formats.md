# WhatsApp video and audio formats — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Add a "MP4 (WhatsApp / any device)" entry for videos, and M4A/FLAC/OGG (plus MP3/WAV) for audio.
**Architecture:** New actions in `internal/actions`. New ffmpeg argument builders in `internal/media`. `Probe` also reads the video's colour transfer, to detect HDR. Dispatch routes the new IDs.
**Spec:** `docs/superpowers/specs/2026-11-25-whatsapp-video-and-audio-formats-design.md`

## Global Constraints
- Readable code: every engine flag has a comment.
- WhatsApp action: ID `whatsapp`, menu label `MP4 (WhatsApp / any device)`, tag ` (WhatsApp)`, kind Heavy.
- Audio actions (all Heavy):
  - `mp3` (.mp3), `m4a` (.m4a), `wav` (.wav), `flac` (.flac), `ogg` (.ogg)
- Video submenu order:
  1. MP4
  2. MP4 (WhatsApp / any device)
  3. MP3 (audio only)
  4. M4A (audio only)
  5. GIF
- Audio submenu order: MP3, M4A, WAV, FLAC, OGG.
- Audio inputs: `.wav .m4a .flac .ogg .aac .wma .mp3 .opus`.
- Same-format items are hidden only for untagged actions, so `.mp4` still gets the WhatsApp entry.

## Review Focus
1. **Portrait iPhone video** (rotation tag 90°) → output must be portrait 1080×1920. *Test `TestWhatsAppPortrait`.*
2. **HDR PQ video** at about 100 nits must come out bright (luma > 180), not grey (about 128). *Test `TestWhatsAppHDRToneMapped`.*
3. **Video with no sound** → the WhatsApp MP4 still works (no audio track). *Test `TestWhatsAppNoAudio`.*
4. **iPhone AAC sound → M4A** is copied without loss. *Test `TestVideoToM4ACopiesAAC`.*
5. **A `.mp3` input** must offer M4A, WAV, FLAC and OGG, but not MP3. *Test `TestSubmenuAudioFormats`.*

### Task 1: Actions and menu
- Tests:
  - `TestSubmenuVideo`: updated to the new order.
  - `TestSubmenuAudioFormats`: `.mp3` gives `[m4a wav flac ogg]`, `.wav` gives `[mp3 m4a flac ogg]`.
  - `TestWhatsAppOfferedForMP4`.
  - menu `TestInstallCascadeKey`: `02mp3` becomes `03mp3`.
- Run `go test ./internal/actions/ ./internal/menu/`. Expected: FAIL, then PASS after the change.

### Task 2: Audio formats
- `func AudioArgs(in, out, format string, i Info) []string`. `MP3Args` and `WAVArgs` become wrappers around it.
- Tests (`internal/media`):
  - `TestAudioFormats`: wav → m4a gives aac, → flac gives flac, → ogg gives vorbis; mp3 → wav gives pcm_s16le.
  - `TestVideoToM4ACopiesAAC`.
- `Audio` probes the input; `Video` handles `m4a`.

### Task 3: WhatsApp MP4
- `Info.ColorTransfer`. `func (i Info) HDR() bool` is true for `smpte2084` and `arib-std-b67`.
- `func WhatsAppArgs(in, out string, i Info) []string`
- Tests:
  - `TestWhatsAppFromHEVC4K`: 3840×2160 HEVC gives h264, yuv420p, 1920×1080, aac, named `clip (WhatsApp).mp4`.
  - `TestWhatsAppPortrait`
  - `TestWhatsAppHDRToneMapped`
  - `TestWhatsAppNoAudio`
- Dispatch `whatsapp` → `media.Video`.
- Then run the full suite with both ffmpeg versions, rebuild the zip, update README.
