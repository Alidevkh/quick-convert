# WhatsApp-ready video and more audio formats — design

**Date:** 2026-11-25 · **Status:** approved by owner · **Size:** small change

## Why
- **Video:** the owner sends iPhone videos on WhatsApp and wants a file that
  plays on any device. Today's "MP4" keeps the iPhone's HEVC video (fast, no
  quality loss), but HEVC does not play everywhere. iPhone HDR also looks
  washed out or dark on normal screens when converted naively.
- **Audio:** the owner wants the popular audio formats, not only MP3.

## Approaches (video)
| | Verdict |
|---|---|
| **New entry next to MP4 (chosen)** | Keeps the fast, lossless MP4 and adds a compatible one |
| Make MP4 always compatible | One entry, but always slow and never lossless |

## Decisions
**Video: "MP4 (WhatsApp / any device)"** (owner approved the approach and these proposals):
- Output `<name> (WhatsApp).mp4`: H.264 (High profile) video, 8-bit 4:2:0
  colours, AAC 128 kbps stereo audio, "fast start" (plays while loading).
- 4K is reduced to **1080p**: the short side is at most 1080 px, so portrait
  videos stay portrait. Smaller videos keep their size.
- The frame rate is kept (30 or 60).
- HDR (PQ or HLG, including iPhone Dolby Vision) is tone-mapped to normal
  colours (BT.709) with ffmpeg's zscale + tonemap, so it looks right on
  every screen.
- Offered for every video type, including `.mp4` (an HEVC `.mp4` needs it too).

**Audio** *(my calls, stated to the owner)*:
- Audio files get **MP3, M4A, WAV, FLAC, OGG** (never their own format).
  `.mp3` and `.opus` (WhatsApp voice notes) become accepted inputs too.
- Videos get **MP3 (audio only)** and **M4A (audio only)**. The menu stays short.
- Quality:
  - MP3: 192 kbps.
  - M4A: AAC 192 kbps. If the source sound is already AAC (iPhone videos),
    it is copied without loss.
  - WAV: 16-bit.
  - FLAC: lossless.
  - OGG: Vorbis at about 192 kbps.

## Not included
Shrinking a video to a target file size, trimming, and Opus output.
