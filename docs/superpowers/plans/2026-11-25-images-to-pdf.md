# Images to PDF and Merge PDFs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add three features:
- "Send to → Images to PDF": many photos become one PDF.
- "Send to → Merge PDFs": many PDFs become one PDF.
- "Convert to → PDF" for every picture type.

**Architecture:**
- **Many files in one run.** A Send-to shortcut runs `quickconvert.exe <action> <file1> <file2> ...`. A new `app.RunMany` filters the files, sorts them in natural order, and links them all into the job's temporary folder as `in-001.<ext>`, `in-002.<ext>`, and so on. Then it reuses the existing handler, commit, notification and log flow.
- **Images to PDF.** ffmpeg prepares each photo as a JPEG: upright, at most 2480 px on its longest side, with a white background. A small pure-Go writer (`pdfwriter`) puts the JPEGs into a PDF without re-compressing them.
- **Merge PDFs.** Merging uses `mutool merge`.

**Tech Stack:** Go standard library only (`image/jpeg` for reading sizes). Engines: ffmpeg, exiftool, mutool (all already bundled).

**Spec:** `docs/superpowers/specs/2026-11-25-images-to-pdf-design.md`

## Global Constraints

- The readability rule from the first plan still applies:
  - every file starts with a comment saying what it does;
  - every exported function has a plain-English doc comment;
  - every engine flag gets a comment explaining it.
- Never overwrite, modify or delete originals. The output goes in the **first** input's folder.
- **Images to PDF** output: `<first photo base>.pdf`. Accepted extensions: `.jpg .jpeg .png .webp .bmp .heic`.
- **Merge PDFs** output: `<first PDF base> (merged).pdf`.
- Order of files: natural sort on the file name, compared case-insensitively, with digit runs compared as numbers.
- Photo pages:
  - longest side at most **2480 px**, never enlarged;
  - orientation applied, so the page looks as it does on screen;
  - transparency filled with white;
  - page size = pixels × 72 / **300** points.
- Messages, exact text:
  - "No pictures selected"
  - "No PDFs selected"
  - "Select at least two PDFs to merge"
  - "This PDF is password-protected"
- Done notification details:
  - images: `"3 pages"` (or `"1 page"`)
  - merge: `"3 files"`
- Send-to shortcut names:
  - `Images to PDF.lnk` runs `quickconvert.exe imgpdf`
  - `Merge PDFs.lnk` runs `quickconvert.exe merge`

  Both live in `%APPDATA%\Microsoft\Windows\SendTo`.

## Review Focus

1. **A Send-to selection that mixes photos with other files** (a .txt, a folder, a .pdf). Only the accepted files are used. If none are left, the user gets the "No … selected" message. *Test: Task 1 `TestRunManyFiltersAndSorts`, `TestRunManyNoneAccepted`.*
2. **`IMG_9.jpg` and `IMG_10.jpg` in the same selection.** IMG_9 must come first. *Test: Task 1 `TestSortNatural`.*
3. **A phone photo stored sideways** (EXIF orientation 6). Its PDF page must be portrait, as shown on screen. *Test: Task 3 `TestImagesToPDFRotated`.*
4. **A transparent PNG.** Its page background must be white, not black. *Test: Task 3 `TestImagesToPDFTransparentWhite`.*
5. **A password-protected PDF inside a merge selection.** The user gets the friendly message, and no half-written file is left. *Test: Task 4 `TestMergeEncrypted`.*

---

### Task 1: Many-input runs

**Files:**
- Modify:
  - `internal/actions/actions.go`
  - `internal/outpath/outpath.go`
  - `internal/app/app.go`
  - `cmd/quickconvert/main.go`
  - `cmd/quickconvert/dispatch.go`
- Create: `internal/app/natural.go`
- Tests:
  - `internal/actions/actions_test.go`
  - `internal/outpath/outpath_test.go`
  - `internal/app/app_test.go`

**Interfaces — Produces:**

- Actions:
  - **"pdf"** — Label `"PDF"`, Light, OutExt `".pdf"`. Added as the last item of the image group's submenu.
  - **"imgpdf"** — Label `"Images to PDF"`, Light, OutExt `".pdf"`.
  - **"merge"** — Label `"Merge PDFs"`, Light, OutExt `".pdf"`, Tag `" (merged)"`.
- `type actions.SendTo struct{ ActionID, Label string; Exts []string; NoneMsg string }`
- `func actions.SendToItems() []SendTo`. Returns two entries, imgpdf then merge:
  - imgpdf: `Exts` = the image extensions, `NoneMsg` = "No pictures selected".
  - merge: `Exts` = `[".pdf"]`, `NoneMsg` = "No PDFs selected".
- `func actions.SendToFor(id string) (SendTo, bool)`
- `Work.Ins []string` holds the work-folder names of all inputs. `LinkInput` sets it to `[]string{In}`.
- `func (w *Work) LinkInputs(srcs []string) error` links each file as `in-%03d<lower ext>` and sets `In = Ins[0]`.
- `func app.SortNatural(paths []string)` sorts in place by base name.
- `func app.RunMany(ctx context.Context, actionID string, inputs []string, d Deps) int`. Its steps:
  1. Unknown or non-send-to action → exit 2.
  2. Keep only regular, non-empty files whose extension is in `Exts`.
  3. If none are left → fail with `NoneMsg`.
  4. Sort the files with `SortNatural`.
  5. The first file sets the folder, `Base` and `Ext`.
  6. `Job.Inputs` = the sorted absolute paths.
  7. Use the same timer, `NewWork`, `LinkInputs`, handler, notify and log flow as `Run`. The log line lists all inputs joined by `", "`.
- `Job.Inputs []string`. `Run` sets it to `[]string{input}`.
- `main`: when `actions.SendToFor(args[1])` succeeds, call `RunMany` with `os.Args[2:]`, which needs at least 1 file. Otherwise the existing single-file rule applies.

- [ ] **Step 1: Write failing tests:**
  - `TestSortNatural`: `["b/IMG_10.jpg","a/IMG_9.jpg","IMG_2.JPG","img_1.jpg"]` sorts to `["img_1.jpg","IMG_2.JPG","a/IMG_9.jpg","b/IMG_10.jpg"]`.
  - `TestSendToItems`: two entries, `imgpdf` and `merge`. `.heic` is in imgpdf's `Exts`.
  - `TestImageSubmenuHasPDF`: `Submenu(".jpg")` ends with `{"pdf","PDF"}`. `Submenu(".pdf")` does not contain `pdf`.
  - `TestLinkInputs`: three sources give `Ins` = `[in-001.jpg in-002.png in-003.jpg]`, with matching content, and `In == "in-001.jpg"`.
  - `TestRunManyFiltersAndSorts`: inputs are `IMG_10.jpg`, `notes.txt`, `IMG_9.jpg`, a folder, and an empty `IMG_1.jpg`. The fake handler records `j.Inputs` and sees exactly `[IMG_9.jpg IMG_10.jpg]`, with `j.Base == "IMG_9"`.
  - `TestRunManyNoneAccepted`: `imgpdf` with only `notes.txt` → exit 1, body "No pictures selected".
  - `TestRunManyUnknownAction`: `RunMany(…,"mp4",…)` → exit 2.
- [ ] **Step 2:** Run `go test ./internal/...`. Expected: FAIL (undefined).
- [ ] **Step 3:** Implement. Map `imgpdf`, `merge` and `pdf` to a temporary `notYet` handler in `dispatch.go`. It is replaced in Tasks 3–4.
- [ ] **Step 4:** Run `go test ./... && GOOS=windows go vet ./...`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: run one action on many files (Send to)`

### Task 2: PDF writer for JPEG pages

**Files:**
- Create: `internal/pdfwriter/pdfwriter.go`, `internal/pdfwriter/pdfwriter_test.go`

**Interfaces — Produces:**
- `type Page struct{ JPEG []byte; Width, Height int }`. Width and Height are in pixels.
- `const PixelsPerInch = 300`
- `func Write(w io.Writer, pages []Page) error`. It writes a PDF 1.4 file:
  - a catalog and a page tree;
  - for each page: a MediaBox of `[0 0 W*72/300 H*72/300]`, and a content stream that draws the image over the whole page;
  - the image stored as an XObject with `/Filter /DCTDecode /ColorSpace /DeviceRGB /BitsPerComponent 8`;
  - a correct xref table and trailer.

  An empty `pages` list returns an error.

- [ ] **Step 1: Write failing tests** (`testutil.Need(t,"ffmpeg","mutool")`, using a JPEG made with `testutil.MakeImage`):
  - `TestWriteTwoPages`: pages of 600×300 and 300×600 px. Then `mutool info` or `mutool pages out.pdf` reports 2 pages, with MediaBoxes `[0 0 144 72]` and `[0 0 72 144]`.
  - `TestWriteKeepsJPEGBytes`: the output contains the input JPEG bytes unchanged (use `bytes.Contains`).
  - `TestWriteRenders`: `mutool draw -r 300 -o p%d.png out.pdf` page 1 is 600 px wide (±1).
  - `TestWriteEmpty`: returns an error.
- [ ] **Step 2:** Run `go test ./internal/pdfwriter/`. Expected: FAIL.
- [ ] **Step 3:** Implement `Write`. Collect object byte offsets while writing, then emit the xref.
- [ ] **Step 4:** Run `go test ./internal/pdfwriter/`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: minimal PDF writer for JPEG pages`

### Task 3: Images to PDF (Send to + Convert to)

**Files:**
- Create: `internal/imagepdf/imagepdf.go`, `internal/imagepdf/imagepdf_test.go`
- Modify: `internal/media/media.go` (add PageArgs), `cmd/quickconvert/dispatch.go`

**Interfaces:**
- Consumes:
  - `media.Orientation`
  - `pdfwriter.Write`
  - `Job.Work.Ins`
- Produces:
  - `func media.PageArgs(in, out string, orientation int) []string`. It is like `ImageArgs(…,"jpg",…)`, with an extra scale step before the white flatten:

    ```
    scale='if(gte(iw,ih),min(2480,iw),-2)':'if(gte(iw,ih),-2,min(2480,ih))'
    ```

    plus `-pix_fmt yuvj420p`.
  - `func imagepdf.Convert(ctx context.Context, j app.Job) (app.Result, error)`. For each name in `j.Work.Ins`:
    1. Run ffmpeg to make `page-%03d.jpg`.
    2. Read the size with `jpeg.DecodeConfig`.

    Then `pdfwriter.Write` writes `out.pdf`, which is committed as `j.Base + ".pdf"`. `Detail` is `"N pages"`, or `"1 page"`.
  - Dispatch: both `"pdf"` (from an image extension) and `"imgpdf"` go to `imagepdf.Convert`.

- [ ] **Step 1: Write failing tests** (they need ffmpeg, ffprobe, exiftool and mutool):
  - `TestImagesToPDFThreePages`: three JPGs → one PDF named after the first. `mutool` shows 3 pages, and `Detail == "3 pages"`.
  - `TestImagesToPDFDownscales`: a 4000×3000 JPG becomes a page of MediaBox `[0 0 595.2 446.4]`. That is 2480×1860 px at 300 ppi, within ±0.5.
  - `TestImagesToPDFNoEnlarge`: a 600×400 JPG becomes MediaBox `[0 0 144 96]`.
  - `TestImagesToPDFRotated`: a 40×20 JPG with Orientation 6 becomes a page taller than it is wide.
  - `TestImagesToPDFTransparentWhite`: a transparent 16×16 PNG becomes a page; render it with mutool, and the first pixel is ≥ 245 in every channel.
  - `TestConvertToPDFSingle`: `RunHandler(Convert, "photo.jpg", "pdf")` gives `photo.pdf` with 1 page.
- [ ] **Step 2:** Run `go test ./internal/imagepdf/`. Expected: FAIL.
- [ ] **Step 3:** Implement and wire the dispatch.
- [ ] **Step 4:** Run `go test ./internal/... ./cmd/...`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: images to PDF`

### Task 4: Merge PDFs

**Files:**
- Create: `internal/pdf/merge.go`
- Modify: `internal/pdf/pdf_test.go`, `cmd/quickconvert/dispatch.go`

**Interfaces — Produces:**
- `func pdf.Merge(ctx context.Context, j app.Job) (app.Result, error)`:
  - If `len(j.Work.Ins) < 2`, return `UserError` "Select at least two PDFs to merge".
  - Otherwise run `mutool merge -o out.pdf in-001.pdf in-002.pdf …`, passing errors through `Friendly`.
  - Commit the result as `j.Base + " (merged)"` with ext `.pdf`. `Detail` is `"N files"`.
- Dispatch: `"merge"` goes to `pdf.Merge`.

- [ ] **Step 1: Write failing tests:**
  - `TestMergeThree`: PDFs of 1, 2 and 3 pages, run through `app.RunMany`. The result is `a (merged).pdf` with 6 pages, and the first page's text matches the first file in natural order. Check it with `mutool draw -F txt`.
  - `TestMergeOnlyOne`: gives "Select at least two PDFs to merge".
  - `TestMergeEncrypted`: one encrypted PDF among the inputs gives "This PDF is password-protected", and no `.qc-*.part` folder or `(merged)` file is left.
- [ ] **Step 2:** Run `go test ./internal/pdf/`. Expected: FAIL.
- [ ] **Step 3:** Implement, wire it, and remove the `notYet` placeholder.
- [ ] **Step 4:** Run `go test ./...`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: merge PDFs`

### Task 5: Installer, docs, package

**Files:**
- Modify:
  - `internal/menu/menu.go`
  - `internal/menu/menu_test.go`
  - `packaging/README.txt`
  - `docs/CODE-GUIDE.md`

**Interfaces — Consumes:** `actions.SendToItems()`.

`install.cmd` gets one line per Send-to item, after the registry lines. The shortcut creation is:

```
powershell -NoProfile -ExecutionPolicy Bypass -Command "$s=(New-Object -ComObject WScript.Shell).CreateShortcut([Environment]::GetFolderPath('SendTo')+'\<Label>.lnk'); $s.TargetPath=$env:DEST+'\quickconvert.exe'; $s.Arguments='<id>'; $s.Save()"
```

`uninstall.cmd` gets `del "%APPDATA%\Microsoft\Windows\SendTo\<Label>.lnk" 2>nul` for each item.

- [ ] **Step 1: Write failing tests:**
  - `TestInstallCreatesSendToShortcuts`: the install script contains `'\Images to PDF.lnk'` with `$s.Arguments='imgpdf'`, and `'\Merge PDFs.lnk'` with `$s.Arguments='merge'`.
  - `TestUninstallRemovesSendToShortcuts`: the uninstall script contains both `del` lines, before the `rmdir` line.
  - `TestInstallJpgHasPDFEntry`: the `.jpg` submenu has a `…\shell\04pdf\command` line.
- [ ] **Step 2:** Run `go test ./internal/menu/`. Expected: FAIL.
- [ ] **Step 3:** Implement. Update README.txt with a "Send to" section, and update CODE-GUIDE.md (Arabic) with the new files.
- [ ] **Step 4:** Run `go test ./... && scripts/build.sh`. Expected: PASS, and the zip builds with the engines.
- [ ] **Step 5:** Commit: `feat: Send to shortcuts in installer; docs`
