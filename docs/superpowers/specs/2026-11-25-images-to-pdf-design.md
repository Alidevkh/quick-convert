# Images to PDF and Merge PDFs — design

**Date:** 2026-11-25 · **Status:** approved by owner (2026-11-25); additions below added the same day · **Size:** small feature

## Why
The owner photographs paper documents with a phone and wants all the photos
combined into **one** PDF, one photo per page, quickly from Windows Explorer.

## What the owner decided
- **Several photos → one PDF.** Every selected photo becomes one page.
- **Entry point: "Send to".** Select photos → right-click → Show more options
  → Send to → **Images to PDF**. Windows passes all selected files to one
  program run, so every photo is guaranteed to end up in the same PDF.
- **Page size = the photo's own size.** No A4 and no fitting. Photos with
  different sizes give pages with different sizes.
- **Orientation as shown on screen.** Phone photos that are stored sideways
  with an "display rotated" tag are turned upright, so each page matches what
  Windows shows.

## Approaches considered
| | How | Verdict |
|---|---|---|
| **Send to (chosen)** | A shortcut in the Send To folder; Windows gives all files in one call | Simple and 100% reliable |
| Right-click entry | Windows starts one process per photo; they must find each other and elect one to combine | Same clicks, but more complex, and a slow start can split a photo into its own PDF |

## Behaviour (owner's decisions plus my calls, marked *(my call)*)
- **Result:** `<first photo name>.pdf`, saved in the first photo's folder, e.g.
  `IMG_2041.jpg` + `IMG_2042.jpg` + `IMG_2043.jpg` → `IMG_2041.pdf` (3 pages).
  Same no-overwrite rule as everything else: `IMG_2041 (1).pdf` if taken. *(my call)*
- **Page order:** by file name, with numbers compared as numbers
  (`IMG_9` before `IMG_10`). Phones name photos in shooting order. Windows does
  not tell programs the order in which you clicked the files. *(my call)*
- **Accepted files:** .jpg .jpeg .png .webp .bmp .heic. Other selected files are
  skipped. If none are left, the message is "No pictures selected". *(my call)*
- **Size of the PDF:** each photo is re-saved as a good-quality JPEG with its
  longest side at most 2480 px. That is sharp enough to read and print an A4
  page, and it keeps the PDF around 0.5–1 MB per page instead of 3–5 MB.
  Smaller photos are never enlarged. *(my call, from the earlier proposal the
  owner did not object to)*
- **Physical page size:** the page has the photo's proportions, measured at
  300 pixels per inch. So a 2480 px wide photo gives an A4-width page, and
  pages keep their sizes relative to each other. *(my call)*
- **Transparency** (PNG) becomes white, as in the existing JPG conversion.
- **Notification:** "Done: IMG_2041.pdf (3 pages)".
- **One photo selected:** a one-page PDF.

## How it is built *(my calls)*
- **ffmpeg** (already bundled) prepares each photo: upright, scaled down,
  white background, saved as JPEG. It reuses the existing orientation and JPG
  logic in `internal/media`.
- **Go code writes the PDF itself.** A PDF holding JPEG pages is a small,
  well-defined file. Writing it ourselves gives exact control of page sizes, and
  the JPEGs are stored as they are, without re-compression. No new engine.
- **Program call:** `quickconvert.exe imgpdf <file1> <file2> ...`. The app gets a
  "many inputs" path next to today's single-file path, with the same temporary
  folder, cleanup, notifications and log.
- **Install / uninstall:** `install.cmd` creates the shortcut
  `%APPDATA%\Microsoft\Windows\SendTo\Images to PDF.lnk` (via one PowerShell
  command, since cmd cannot make shortcuts). `uninstall.cmd` deletes it.

## Additions requested with the approval
1. **Merge PDFs.** Select PDFs → Send to → **Merge PDFs**. The result is one PDF
   holding all pages: `<first PDF name> (merged).pdf`, saved next to the first
   PDF. Files are joined in file-name order, with numbers compared as numbers,
   as for photos. Pages are copied as they are, using `mutool merge` (bundled),
   so nothing is re-compressed. A password-protected PDF gives
   "This PDF is password-protected". With only one PDF selected, the message is
   "Select at least two PDFs to merge". *(my calls: name, order, messages)*
2. **Image → PDF in "Convert to".** Every picture type gets a **PDF** entry in
   its Convert to submenu. It makes a one-page PDF, `photo.pdf`, with the same
   page rules as above. Selecting several photos and using Convert to gives one
   PDF per photo, because Windows runs each separately. **Send to** is the way to
   combine them.

## Not included
Auto-cropping the paper edges, black-and-white "scanner" look, choosing the
page order by hand, and OCR. They can come later.

## Known limits
- **Number of photos:** Windows limits a program's command line to about
  32,000 characters. That is a few hundred photos with normal folder names.
- **HEIC:** HEIC photos depend on the bundled Windows ffmpeg being able to read
  them. This is the same open question as the existing HEIC conversion.
