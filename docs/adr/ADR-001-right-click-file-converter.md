# ADR-001: Right-click file converter for Windows ("QuickConvert")

**Status:** Proposed
**Date:** 2026-11-24
**Deciders:** Project owner (sole user)

## Context

The owner wants to change and clean up files fast (convert formats, shrink PDFs, strip photo
metadata, turn PDFs into Markdown for AI) without opening an app, uploading anywhere, or moving
files around.

What the owner decided:
- **Right-click menu only for now.** A floating widget is out of scope because it has to keep running in the background.
- **Every result is saved next to the original**, in the same folder. The original is kept.
- **Features:** format conversion (video, audio, image, PDF), TXT ↔ code-file conversion, PDF size
  reduction, photo metadata removal, and PDF → Markdown.
- **Metadata removal is for photos only.**
- **PDF → Markdown is for normal text PDFs, to feed into AI tools. No OCR.**

Constraints and forces:
- **Windows 10/11, personal use.** One user, one PC, and no admin rights should be needed.
- **Nothing runs in the background.** The tool starts when a menu item is clicked and exits when the work is done.
- **Windows 11 has two right-click menus.** The short "modern" menu only accepts entries from
  signed, packaged apps. Ordinary registry entries appear under **"Show more options"**
  (or Shift+right-click).
- **Built on Linux, used on Windows.** Development happens in a Linux container, so the tool must
  cross-compile to a Windows .exe, and the owner does the final testing on Windows.
- **Big files.** Large videos can take minutes, so the user needs to see that something is happening.

## Decision

Build **Option B**: a small Go program (`quickconvert.exe`) that registers right-click entries per
file type under the current user's registry and drives bundled command-line engines:

| Engine | Used for |
|--------|----------|
| `ffmpeg.exe` | video, audio and image format conversion |
| `mutool.exe` (MuPDF) | PDF → PNG/JPG, PDF → Markdown (via its structured-text output) |
| `gswin64c.exe` (Ghostscript) | PDF size reduction |
| `exiftool.exe` | photo metadata removal |
| *(none: done in Go)* | TXT ↔ code files |

The program shows a Windows notification when work starts and when it finishes, and writes every
result into the original file's folder.

## Options Considered

### Option A: PowerShell scripts + registry entries
| Dimension | Assessment |
|-----------|------------|
| Complexity | Low |
| Robustness | Low: a console window flashes, and PowerShell takes about 0.5–1 s to start per file |

**Pros:** No compiling, and the owner can edit the scripts by hand.
**Cons:** A console window flashes on every run. Windows notifications from PowerShell are awkward.
Paths with spaces or Unicode characters are fragile. Some PCs block scripts by default (execution policy).
The PDF → Markdown logic would be painful to write in PowerShell.

### Option B: Small compiled Go helper + registry entries ✅ chosen
| Dimension | Assessment |
|-----------|------------|
| Complexity | Medium-low |
| Robustness | High: one .exe, no console window, no runtime to install |
| Buildability here | Go cross-compiles to Windows from Linux with a single command |

**Pros:** Silent and fast to start. Real notifications. Clean handling of names, duplicates and errors.
The PDF → Markdown logic lives in testable Go code.
**Cons:** Entries sit under "Show more options" on Windows 11.

### Option C: Native Windows 11 menu (packaged app with a shell extension)
| Dimension | Assessment |
|-----------|------------|
| Complexity | High |
| Cost | Needs a code-signing certificate, or self-signed certificates with developer mode |
| Buildability here | Poor: the packaging and signing tools are Windows-only |

**Pros:** Entries appear in the short Windows 11 menu, with no extra click.
**Cons:** Several times the work. Certificate hassle. Can't be built or tested from Linux.

## Trade-off Analysis

The goal is **speed with no friction**. Option A brings back the "an app is opening" feel (console
flash, lag) that the owner wants to avoid. Option C removes one click on Windows 11 but costs signing,
packaging and Windows-only tooling. Option B gets almost all of the benefit for a fraction of the work.
Option C can be layered on later using the same converter .exe.

## Design details

### Menu
Each supported file gets the entries that make sense for it, under the classic menu:

```
movie.mov   →  Convert to ▸  MP4 · MP3 (audio only) · GIF
photo.jpg   →  Convert to ▸  PNG · JPG · WEBP
               Remove metadata
report.pdf  →  Convert to ▸  PNG · JPG · MD
               Reduce size
notes.txt   →  Convert to ▸  .py · .m · .c · .cpp · .h · .js · .ts · .html · .css · .java · .json · .md · .bat · .sh
script.py   →  Convert to ▸  .txt
```

### Feature table
| Input | Action | Output (same folder) | Engine |
|-------|--------|----------------------|--------|
| Video: .mov .mp4 .mkv .avi .webm .wmv | → MP4 / MP3 / GIF | `trip.mp4`, `trip.mp3`, `trip.gif` | ffmpeg |
| Audio: .wav .m4a .flac .ogg .aac .wma | → MP3 / WAV | `song.mp3` | ffmpeg |
| Image: .png .jpg .jpeg .webp .bmp .heic | → PNG / JPG / WEBP | `photo.png` | ffmpeg |
| Photo: .jpg .jpeg .png .heic .webp | Remove metadata | `photo (clean).jpg` | exiftool |
| PDF | → PNG / JPG | `report-1.png`, `report-2.png`, … (single page: `report.png`) | mutool |
| PDF | → MD | `report.md` | mutool + Go |
| PDF | Reduce size | `report (small).pdf` | Ghostscript |
| .txt | → code extension | `notes.py` (identical content) | Go |
| Code files (the extensions above) | → .txt | `script.txt` (identical content) | Go |

### Output rules (apply to every feature)
- The result goes in the original's folder with the original's base name.
- If the name is taken, add ` (1)`, ` (2)`, …, e.g. `trip (1).mp4` or `photo (clean) (1).jpg`.
- **Never overwrite, never delete, never modify the original.**
- Write to a temporary `*.part` name first and rename on success, so a failed or interrupted run
  leaves no half-written file that looks real. On failure, the temporary file is removed.

### Feature specifics
**Video/audio (my calls):**
- **MOV/MKV → MP4:** if the streams are already MP4-compatible (H.264/HEVC + AAC), the file is rewrapped
  without re-encoding, which takes seconds. Otherwise it is re-encoded to H.264/AAC. HEVC is tagged
  `hvc1` so Windows and phones play it.
- **→ MP3:** 192 kbps.
- **→ GIF:** capped at 480 px wide and 15 fps so files stay a sane size.

**Images:** converted at high quality. When converting to JPG, transparency is filled with white.

**Photo metadata removal:**
- Removes all metadata: GPS location, camera/lens, dates, software, thumbnails and comments.
- **Keeps the orientation tag**, so photos don't appear rotated.
- Lossless: pixels are not re-encoded.

**PDF → PNG/JPG:** rendered at 150 DPI.

**PDF reduce size:**
- Ghostscript `/ebook` profile: images downsampled to about 150 DPI and recompressed. Text stays sharp.
- If the result is not at least 10% smaller than the original, discard it and notify:
  "Already optimized, no smaller file made".
- The finish notification shows the size change, e.g. "4.2 MB → 0.9 MB".

**PDF → Markdown (for AI use):**
- Extract structured text with font size and style for each line (mutool structured-text output). Go code then turns it into Markdown:
  - Headings come from font size relative to the body text (`#`, `##`, `###`).
  - **Bold** and *italic* are kept.
  - Bulleted and numbered lists become Markdown lists.
  - Paragraphs are re-joined across line breaks, and words hyphenated at line ends are rejoined.
  - Headers, footers and page numbers that repeat on every page are removed (noise for AI).
  - Tables become Markdown tables where the column layout is clear. Otherwise they are left as plain lines.
  - Images are skipped.
- **No OCR.** If the PDF has no extractable text (a scan), make no file and notify:
  "This PDF has no selectable text (scanned?)".

**TXT ↔ code:** a byte-for-byte copy with the new extension. No content changes. Windows hides
extensions by default, which is why a menu entry is useful here.

### Feedback
- A notification on start ("Converting trip.mov → MP4…"). It is skipped for actions that finish in under about 1 second, such as the TXT copy.
- A notification on finish ("Done: trip.mp4") or failure ("Couldn't convert trip.mov", with a short reason).
- Errors are logged to `%LOCALAPPDATA%\QuickConvert\log.txt`.

### Selecting many files
Windows launches the program once per selected file. At most **2 heavy jobs** (video/audio
encoding, PDF size reduction) run at once, and the rest wait their turn. Light jobs are not limited.

### Install / uninstall
- `install.cmd` copies the files to `%LOCALAPPDATA%\QuickConvert\` and adds the menu entries under
  `HKEY_CURRENT_USER`. No admin rights needed.
- `uninstall.cmd` removes the menu entries and the folder.
- Delivered as one .zip (about 150 MB, mostly ffmpeg) containing the .exe, the engines, both scripts and `README.txt`.

## Goals / Non-goals / Anti-goals
- **Goals:** right-click → result in the same folder, quickly, with no window to manage and no internet.
- **Non-goals (for now):**
  - the floating widget
  - OCR for scanned PDFs
  - metadata removal for anything other than photos
  - editing (trim, crop, resize)
  - quality choices in the menu
  - Word/Excel conversion
- **Anti-goals:**
  - anything that overwrites, modifies or deletes the original
  - anything that uploads files
  - anything that stays running in the background
  - a console window popping up

## Consequences
- **Easier:** adding a format or action is a table entry plus a rebuild.
- **Harder:** the extra "Show more options" click on Windows 11 remains.
- **Licensing:** MuPDF and Ghostscript are AGPL. That's fine for personal use; selling or distributing the tool would need a rethink.
- **Revisit:**
  - Native Windows 11 menu (Option C).
  - The floating widget, as a separate project reusing the converter.
  - OCR, if scanned PDFs come up.
  - Quality options, if the defaults are wrong.

## Builder notes (left to the implementer)
- Register menus via `HKCU\Software\Classes\SystemFileAssociations\<ext>\shell\…` with cascading
  `ExtendedSubCommandsKey` submenus. This doesn't take over the default app for any file type.
- Build with `GOOS=windows` and the `-H windowsgui` linker flag. Start the engines with a hidden window
  (`CREATE_NO_WINDOW`) so no console flashes.
- Use a named Windows semaphore for the 2-heavy-jobs limit.
- Testing:
  - Unit-test naming rules, command building and PDF → Markdown conversion on Linux, with sample PDFs as fixtures.
  - Run the real engine commands on Linux with the same arguments.
  - The owner does the Windows smoke test.

## Action Items
1. [ ] Owner approves this ADR.
2. [ ] Write the implementation plan.
3. [ ] Build: core (output naming, temporary-file rename, job runner, logging, notifications), with tests.
4. [ ] Build: ffmpeg conversions (video, audio, image).
5. [ ] Build: PDF features (→ PNG/JPG, reduce size, → MD).
6. [ ] Build: photo metadata removal and TXT ↔ code.
7. [ ] Build: install/uninstall registry scripts.
8. [ ] Package the .zip.
9. [ ] Owner smoke-tests on Windows 10/11.
