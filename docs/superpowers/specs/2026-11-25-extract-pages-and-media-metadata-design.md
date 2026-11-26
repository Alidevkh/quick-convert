# Extract PDF pages, and remove metadata from video and audio — design

**Date:** 2026-11-25 · **Status:** approved by owner · **Size:** small change

## Extract pages… (PDF)
- **Why:** the owner wants to pick some pages of a PDF into one new PDF.
- **Menu:** right-click a PDF → **Extract pages…**. A small Windows input box asks
  "Which pages? (1–N)", e.g. `1-3, 5`. This is the only entry that opens a window.
- **Input rules:**
  - Ranges like `3-5` and single pages; separators can be commas or spaces.
  - Order is kept, so `5, 1-3` puts page 5 first. A reversed range `5-3` gives 5, 4, 3.
  - Arabic-Indic digits (٣) are accepted. *(my call)*
  - Pages outside 1–N or unreadable text give a clear message.
  - Cancel or an empty answer does nothing.
- **Output:** one PDF, `report (pages 1-3, 5).pdf`, next to the original. A very long
  selection gives `report (selected pages).pdf`. *(my call)*
- **How:** `mutool merge -o out.pdf in.pdf 1-3,5` copies the pages as they are.
  The PDF is checked to open first; a password-protected one gives the usual message.
- **Input box:** Windows' standard InputBox, shown through a hidden PowerShell, with
  text passed in environment variables. *(my call)*

## Remove metadata for video and audio
- **Why:** the owner's phone videos carry the GPS location. The owner wants metadata
  removal for videos and audio "as for photos, in general — not only WhatsApp".
- **Menu:** right-click any video or audio file → **Remove metadata** →
  `<name> (clean).<ext>`.
- **How:** ffmpeg copies the picture and sound as they are (no quality loss) and
  drops all container and stream metadata (`-map_metadata -1`), chapters, and
  data tracks (iPhone location and timed-metadata tracks). Checked: the
  rotation is kept, so portrait videos stay portrait.
- **Album covers:** a music file's cover picture is kept (it isn't personal data).
  *(my call)*
- **WhatsApp MP4:** now also drops all metadata, since it is meant for sharing.
