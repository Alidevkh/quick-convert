# QuickConvert Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A Windows right-click tool that converts, shrinks, cleans and turns files into Markdown, and saves the result next to the original.

**Architecture:** One Go program, `quickconvert.exe <action> <file>`, launched by registry context-menu verbs. Each run creates a temporary work folder inside the file's own folder and runs a bundled engine there (ffmpeg/ffprobe, mutool, Ghostscript, exiftool), always using short ASCII relative names. It then "commits" the results to unique final names with a no-overwrite hard link, and deletes the work folder. One action table drives both the dispatch and the generated `install.cmd`/`uninstall.cmd`, so the menu and the program can't drift apart.

**Tech Stack:** Go 1.24 standard library only (no third-party modules). Windows APIs are reached through `syscall`. Engines: ffmpeg/ffprobe, MuPDF `mutool`, Ghostscript, ExifTool. Tests run on Linux with the Linux builds of those engines (installed in this container: ffmpeg 6.1, mutool 1.23, gs 10.02, exiftool 12.76). Windows-only code is checked by cross-compiling.

**Spec:** `docs/adr/ADR-001-right-click-file-converter.md`

## Global Constraints

- Module path `quickconvert`. Go 1.24. **No third-party Go modules.**
- **Never overwrite, modify or delete the original. Never overwrite any existing file.**
- Output goes in the original's folder. Duplicates become `name (1).ext`, `name (2).ext`, …; the suffix comes after any tag, e.g. `photo (clean) (1).jpg`.
- Output names by feature:
  - conversions: `trip.mp4`
  - photo cleaning: `photo (clean).jpg`
  - PDF shrink: `report (small).pdf`
  - PDF pages: `report-1.png`, `report-2.png`, …, or `report.png` for a single page
  - text: `notes.py`
- Work happens in a temporary folder `<dir>/.qc-<random>.part`. Results move out only on success, and the folder is always removed.
- No console window may ever appear. Build with `-ldflags "-H windowsgui"`, and start engines with `CREATE_NO_WINDOW`.
- At most **2 heavy jobs** run at once (video/audio/GIF encoding and PDF shrink), across processes, using the named semaphore `Local\QuickConvertHeavy`.
- Notifications:
  - Start: `Converting <file> → <LABEL>…`, shown only if the job is still running after 1 s.
  - Finish: `Done: <outname>`.
  - Failure: `Couldn't convert <file>` plus a short reason.
- Log file: `%LOCALAPPDATA%\QuickConvert\log.txt`. Use `os.UserCacheDir()` + `QuickConvert/log.txt`. Truncate it when it exceeds 1 MB.
- No admin rights: registry entries go under `HKCU\Software\Classes` only.
- Fixed values:
  - Video → MP3: 192 kbps.
  - Video → GIF: 15 fps, at most 480 px wide.
  - PDF → images: 150 DPI.
  - PDF shrink: Ghostscript `/ebook`; the result is kept only if it is at least 10% smaller.
  - Re-encoded MP4: H.264 + AAC 192k.
  - Remuxed HEVC gets the `hvc1` tag.

## Review Focus

1. **Non-ASCII, spaces and apostrophes in file or folder names** (e.g. `C:\Users\سارة\Desktop\رحلة 'صيف'.mov`). Every feature must work, and the notification must show the name correctly. Engines only ever see relative ASCII names, and notification text is passed through environment variables. *Test owners: Task 1 (naming), Task 5 (one end-to-end conversion in such a folder), Task 4 (toast script uses `$env:` only).*
2. **20+ files selected at once, including two inputs that map to the same output** (`a.mov` and `a.mkv` both → MP4). Every file gets converted, and the results get distinct names (`a.mp4`, `a (1).mp4`). *Test owner: Task 1 (concurrent commits to the same name).*
3. **Phone media:**
   - iPhone MOVs (HEVC plus extra data/subtitle tracks) → MP4 must remux and play.
   - EXIF-rotated JPEGs → PNG/JPG/WEBP must come out upright, with no double rotation on newer ffmpeg.
   - Photo cleaning keeps the colour profile, so P3 photos don't shift colour.

   *Test owners: Task 5 (extra track and rotation tests), Task 6 (ICC kept).*
4. **Bad input leaves nothing behind:** a zero-byte or corrupt file, a password-protected PDF, a video with no audio → MP3. Each gives a friendly message, no stray `.part` folder, and an untouched original. *Test owners: Task 3 (cleanup on failure), Task 5 (no audio), Task 7 (password PDF).*
5. **Read-only folder or failed hard link** (USB stick formatted FAT32/exFAT, a network share). A read-only folder gives the message "Can't save in this folder". When a hard link fails, the commit falls back to a no-overwrite copy. *Test owner: Task 1 (link-failure fallback via injectable link func; the read-only test is skipped when running as root).*

---

## File Structure

```
quick-convert/
  go.mod
  cmd/quickconvert/main.go          # arg parsing, wiring real deps, exit code
  cmd/genscripts/main.go            # writes install.cmd / uninstall.cmd into a folder
  internal/outpath/outpath.go       # naming rules, Work folder, LinkInput, Commit, CommitPages
  internal/actions/actions.go       # THE action table: ids, labels, extensions, heavy flag, menu items
  internal/engines/engines.go       # locate engines, Run (cwd = work dir), EngineError
  internal/engines/sysproc_windows.go / sysproc_other.go   # CREATE_NO_WINDOW
  internal/app/app.go               # Run(): validate, timer toast, semaphore, dispatch, log, messages
  internal/app/errors.go            # UserError, Info
  internal/platform/notify_windows.go / notify_other.go    # toast via hidden PowerShell
  internal/platform/limit_windows.go / limit_other.go      # named semaphore
  internal/media/media.go           # Probe + ffmpeg arg builders + handlers (video/audio/image)
  internal/photo/clean.go           # exiftool metadata removal
  internal/textcopy/textcopy.go     # txt <-> code byte copy
  internal/pdf/render.go            # PDF -> PNG/JPG
  internal/pdf/shrink.go            # PDF size reduction + FormatSize
  internal/pdfmd/stext.go           # parse mutool stext XML into Lines
  internal/pdfmd/markdown.go        # Lines -> Markdown
  internal/pdfmd/pdfmd.go           # handler: run mutool, convert, commit
  internal/menu/menu.go             # generate install/uninstall .cmd text from actions
  internal/testutil/testutil.go     # Need(t, engine), fixture generators
  scripts/fetch-engines.sh          # download + extract Windows engines into engines/windows/
  scripts/build.sh                  # cross-compile, gen scripts, assemble dist/QuickConvert.zip
  packaging/README.txt              # user-facing instructions shipped in the zip
```

Shared handler contract (defined in Task 3, used by Tasks 5–8):

```go
// internal/app
type Job struct {
    Input  string            // absolute path of the original
    Base   string            // original file name without extension, e.g. "trip"
    Ext    string            // lower-case original extension with dot, e.g. ".mov"
    Action actions.Action
    Work   *outpath.Work     // temp folder; input already linked as Work.In ("in"+Ext)
    Eng    *engines.Set
}
type Result struct {
    Outputs []string // final committed paths
    Detail  string   // optional extra for the Done toast, e.g. "4.2 MB → 0.9 MB"
}
type Handler func(ctx context.Context, j Job) (Result, error)
```

---

### Task 1: Project scaffold + output naming and commit rules

**Files:**
- Create: `go.mod`, `internal/outpath/outpath.go`, `internal/outpath/outpath_test.go`, `.gitignore` (`dist/`, `engines/`)

**Interfaces:**
- Produces:
  - `func Candidate(base string, n int, ext string) string`: `n==0` gives `base+ext`, otherwise `base+" ("+n+")"+ext`.
  - `func PageName(base string, n, page int, ext string) string`: `base[+" (n)"]+"-"+page+ext`.
  - `type Work struct{ Dir, Parent, In string; link func(old, new string) error }`
  - `func NewWork(parent string) (*Work, error)`: creates `parent/.qc-<8 random hex>.part`. If the folder can't be created, returns `ErrNotWritable`.
  - `var ErrNotWritable = errors.New("Can't save in this folder")`
  - `func (w *Work) LinkInput(src string) error`: hard-links `src` to `w.Dir/in<ext>` and sets `w.In` to the relative name `in<ext>`. Falls back to a byte copy.
  - `func (w *Work) Commit(tmpName, base, ext string) (string, error)`: moves `w.Dir/tmpName` to the first free `Candidate(base,n,ext)` in `w.Parent`. A free name is claimed with `w.link` (which fails if the target exists), retrying n+1 on `os.IsExist`. If link is unsupported (any other error), it falls back to `os.OpenFile(O_CREATE|O_EXCL)` + copy. Returns the final path.
  - `func (w *Work) CommitPages(tmpNames []string, base, ext string) (app.Result, error)`: picks the smallest n for which every `PageName(base,n,i,ext)` (i = 1…len) is free, then commits each page.
  - `func (w *Work) Cleanup()`: `os.RemoveAll(w.Dir)`.

- [ ] **Step 1:** `git init` is already done in `/home/user/quick-convert`. Run `go mod init quickconvert` and create `.gitignore`.
- [ ] **Step 2: Write failing tests** in `outpath_test.go`:
  - `TestCandidate`: `Candidate("trip",0,".mp4")=="trip.mp4"`, `Candidate("photo (clean)",1,".jpg")=="photo (clean) (1).jpg"`.
  - `TestPageName`: `PageName("report",0,2,".png")=="report-2.png"`, `PageName("report",1,1,".png")=="report (1)-1.png"`.
  - `TestCommitSkipsExisting`: with `trip.mp4` and `trip (1).mp4` already in the parent folder, Commit returns `…/trip (2).mp4`, both existing files are byte-identical to before, and the temp file is gone from `w.Dir`.
  - `TestCommitConcurrentSameName`: 10 goroutines, each with its own Work in the same parent, all commit `base "a", ext ".mp4"`. The test expects 10 distinct final paths and 10 files on disk.
  - `TestCommitLinkUnsupportedFallsBackToCopy`: set `w.link` to return `errors.New("not supported")`. Commit still produces `a.mp4` with the right bytes, and with `a.mp4` already present it produces `a (1).mp4`.
  - `TestCommitPagesPicksCommonN`: with only `report-2.png` existing, a 3-page commit gives `report (1)-1.png`, `report (1)-2.png` and `report (1)-3.png`.
  - `TestUnicodeNames`: in parent folder `مجلد أ`, Commit of base `رحلة 'صيف'` with ext `.mp4` gives that exact file name.
  - `TestNewWorkNotWritable`: `chmod 0555` on the parent, then `errors.Is(err, ErrNotWritable)`. Call `t.Skip` when `os.Geteuid()==0`.
  - `TestCleanupRemovesDir`.
- [ ] **Step 3:** Run `cd /home/user/quick-convert && go test ./internal/outpath/`. Expected: FAIL (undefined).
- [ ] **Step 4:** Implement `outpath.go` to the interfaces above. `NewWork` sets `link: os.Link`.
- [ ] **Step 5:** Run `go test -race ./internal/outpath/`. Expected: PASS.
- [ ] **Step 6:** Commit: `git add -A && git commit -m "feat: output naming and no-overwrite commit"`

---

### Task 2: Action table

**Files:**
- Create: `internal/actions/actions.go`, `internal/actions/actions_test.go`

**Interfaces:**
- Produces:
  - `type Kind int`, with constants `Light` and `Heavy`.
  - `type Action struct{ ID, Label string; Kind Kind; OutExt, Tag string }`. `Tag` is the name suffix, e.g. `" (clean)"`, `" (small)"`, or `""`.
  - `type MenuItem struct{ ActionID, Label string }`
  - `func Lookup(id string) (Action, bool)`
  - `func Submenu(ext string) []MenuItem`: the "Convert to ▸" items, in order.
  - `func Extras(ext string) []MenuItem`: the top-level verbs ("Remove metadata", "Reduce size").
  - `func Supports(id, ext string) bool`
  - `func AllExts() []string`: sorted.

Table values, from the spec:

| Group | Exts | Submenu (ID:Label) | Kind |
|---|---|---|---|
| video | .mov .mp4 .mkv .avi .webm .wmv | mp4:MP4, mp3:"MP3 (audio only)", gif:GIF | Heavy |
| audio | .wav .m4a .flac .ogg .aac .wma | mp3:MP3, wav:WAV | Heavy |
| image | .png .jpg .jpeg .webp .bmp .heic | png:PNG, jpg:JPG, webp:WEBP | Light |
| pdf | .pdf | png:PNG, jpg:JPG, md:MD | Light |
| text | .txt | code-py:.py code-m:.m code-c:.c code-cpp:.cpp code-h:.h code-js:.js code-ts:.ts code-html:.html code-css:.css code-java:.java code-json:.json code-md:.md code-bat:.bat code-sh:.sh | Light |
| code | the 14 code exts above | txt:.txt | Light |

Extras:
- `clean`: "Remove metadata", Light, tag `" (clean)"`, OutExt = the input's ext. Applies to .jpg .jpeg .png .heic .webp.
- `shrink`: "Reduce size", Heavy, tag `" (small)"`, OutExt `.pdf`. Applies to .pdf.

Rules:
- An item whose output equals the input format is left out. `.jpeg` counts as `jpg`.
- `png`/`jpg` are the same action IDs for image and PDF inputs. Dispatch picks the handler by input group (Task 4).
- The label for `mp3` depends on the group, which is why `Submenu` returns its own labels.

- [ ] **Step 1: Write failing tests:**
  - `TestSubmenuVideo`: `Submenu(".mov")` gives the IDs `[mp4 mp3 gif]` and labels `[MP4, MP3 (audio only), GIF]`, and `Submenu(".mp4")` gives `[mp3 gif]`.
  - `TestSubmenuImageExcludesSame`: `Submenu(".jpeg")` gives `[png webp]`.
  - `TestSubmenuText`: `.txt` has 14 items, the first being `code-py:.py`, and `Submenu(".py")` is `[txt:.txt]`.
  - `TestExtras`: `Extras(".pdf")` is `[shrink:Reduce size]`, `Extras(".heic")` is `[clean:Remove metadata]`, and `Extras(".mov")` is empty.
  - `TestSupports`: `Supports("mp4",".txt")==false`, `Supports("clean",".png")==true`.
  - `TestLookupCode`: `Lookup("code-py")` gives `OutExt ".py"`, Light.
  - `TestUppercaseExt`: `Submenu(".MOV")` equals `Submenu(".mov")`.
- [ ] **Step 2:** Run `go test ./internal/actions/`. Expected: FAIL.
- [ ] **Step 3:** Implement the table as package-level data, plus the functions above. Extensions are matched case-insensitively.
- [ ] **Step 4:** Run `go test ./internal/actions/`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: action table`

---

### Task 3: Engines runner + app core (errors, cleanup, messages)

**Files:**
- Create: `internal/engines/engines.go`, `internal/engines/sysproc_windows.go`, `internal/engines/sysproc_other.go`, `internal/engines/engines_test.go`
- Create: `internal/app/app.go`, `internal/app/errors.go`, `internal/app/app_test.go`
- Create: `internal/testutil/testutil.go`

**Interfaces:**
- Consumes: `outpath.*` (Task 1) and `actions.*` (Task 2).
- Produces (engines):
  - `type Set struct{ Dirs []string }`
  - `func Default() *Set`: searches the exe's folder first, then (Windows only) `%ProgramFiles%\gs\gs*\bin`, then `PATH`.
  - `func (s *Set) Path(name string) (string, error)`: `name` is one of `ffmpeg ffprobe mutool gs exiftool`. On Windows, `gs` maps to `gswin64c.exe` and the others get `.exe`. A missing engine returns `&MissingError{Name}`, whose `Error()` is `"<exe> is missing — reinstall QuickConvert"`.
  - `func (s *Set) Run(ctx context.Context, dir, name string, args ...string) ([]byte, error)`: runs the engine with `cmd.Dir = dir` and returns stdout. A non-zero exit returns `&EngineError{Name string; Tail string}`, where Tail is the last non-empty stderr line, trimmed to 160 chars.
  - `func hide(cmd *exec.Cmd)`: on Windows, sets `SysProcAttr{HideWindow:true, CreationFlags:0x08000000}`; elsewhere it does nothing.
- Produces (app):
  - The `Job` and `Handler` types shown in File Structure.
  - `type UserError struct{ Msg string }`: a friendly failure, shown as the reason.
  - `type Info struct{ Msg string }`: not a failure. It means "nothing created" (e.g. a PDF that's already optimized), and the notification title is `QuickConvert` with Msg as the body.
  - `type Deps struct{ Notify func(title, body string); Acquire func() (release func()); Log io.Writer; Eng *engines.Set; Handlers func(a actions.Action, ext string) Handler; StartDelay time.Duration }`
  - `func Run(ctx context.Context, actionID, input string, d Deps) int`. It returns the exit code: 0 for success or Info, 1 for failure, 2 for usage. Its sequence:
    1. Validate the action and that `Supports(id, ext)` is true. Otherwise notify failure with reason `"<LABEL> isn't available for <ext> files"`.
    2. Check the input exists and is a regular file.
    3. Start a `time.AfterFunc(StartDelay)` start notification.
    4. For Heavy actions, call `Acquire`.
    5. `NewWork(dir)` followed by `defer Cleanup()`.
    6. `LinkInput`.
    7. Call the handler.
    8. Stop the timer.
    9. Notify:
       - success: `"Done: <base names joined by ', '>"` + `" ("+Detail+")"` when Detail is non-empty. More than 3 outputs gives `"Done: <first> + N more"`.
       - failure: title `"Couldn't convert <file name>"` with the reason taken from UserError.Msg, ErrNotWritable, MissingError, or `EngineError.Tail`.
    10. Append one log line: `RFC3339 action input -> outputs | error`.

- [ ] **Step 1:** `testutil.go`: `func Need(t *testing.T, engines ...string)` skips the test when an engine isn't on PATH.
- [ ] **Step 2: Write failing tests:**
  - `engines_test.go`:
    - `TestRunUsesDir`: `Run(ctx, tmp, "sh"…)` is not allowed, so add the test-only helper `s.runExe(ctx, dir, exe, args…)` shared by Run. Running `pwd` must report tmp.
    - `TestEngineErrorTail`: an `ffmpeg -i missing.mov` failure gives an `*EngineError` with a non-empty Tail and Name `"ffmpeg"` (needs ffmpeg).
    - `TestMissingEngine`: with `Set{Dirs: []string{emptyDir}}` and PATH emptied via `t.Setenv("PATH","")`, the error is a `*MissingError`.
  - `app_test.go`, using fake handlers and a recording Notify:
    - `TestRunSuccessNotifies`: a handler that commits a file gives exit 0 and the last notify title `"Done: x.py"`.
    - `TestRunFailureCleansUp`: a handler that writes into `j.Work.Dir` and then returns `UserError{"boom"}` gives exit 1 and notify `("Couldn't convert x.txt","boom")`. No `.qc-*.part` folder is left in the parent, and the original is unchanged (same bytes and mtime).
    - `TestRunUnsupported`: `Run(ctx,"mp4","/tmp/x.txt")` gives exit 1 with reason `"MP4 isn't available for .txt files"`.
    - `TestRunDetail`: handler returning `Result{Outputs:[x.pdf], Detail:"2 MB → 1 MB"}` → body `"Done: x.pdf (2 MB → 1 MB)"`.
  - `TestRunInfo`: a handler returning `Info{"Already optimized — no smaller file made"}` gives exit 0, title `QuickConvert`, and that body.
    - `TestStartNotificationOnlyWhenSlow`: with StartDelay 50 ms, a handler sleeping 120 ms produces 2 notifications, the first being `"Converting x.mov → MP4…"`. A handler that returns at once produces exactly 1.
    - `TestHeavyAcquires`: for a Heavy action, Acquire is called once and release is called. For a Light action, Acquire is not called.
    - `TestLogLine`: Log receives one line containing the action, input and output.
- [ ] **Step 3:** Run `go test ./internal/engines/ ./internal/app/`. Expected: FAIL.
- [ ] **Step 4:** Implement both packages.
- [ ] **Step 5:** Run `go test -race ./internal/...` and `GOOS=windows go vet ./internal/...`. Expected: PASS and no output.
- [ ] **Step 6:** Commit: `feat: engine runner and app core`

---

### Task 4: Platform glue + main

**Files:**
- Create: `internal/platform/notify_windows.go`, `notify_other.go`, `limit_windows.go`, `limit_other.go`, `internal/platform/platform_test.go`
- Create: `cmd/quickconvert/main.go`, `cmd/quickconvert/dispatch.go`, `cmd/quickconvert/dispatch_test.go`

**Interfaces:**
- Consumes: `app.Run`, `app.Deps`, `app.Handler`, and the handlers from Tasks 5–8, which are wired in as they land. Until then, the dispatch returns a handler that gives `UserError{"Not built yet"}`.
- Produces:
  - `func platform.Notify(title, body string)`. On Windows it runs `powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand <base64 UTF-16LE of ToastScript>` with env `QC_TITLE` and `QC_BODY`, hidden, using `cmd.Start()` with no Wait. Elsewhere it writes `title: body` to stderr.
  - `const platform.ToastScript string`: a ToastText02 template sent via `CreateToastNotifier('{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\WindowsPowerShell\v1.0\powershell.exe')`. Text comes **only** from `$env:QC_TITLE` and `$env:QC_BODY`.
  - `func platform.AcquireHeavy() (release func())`. On Windows it uses kernel32 `CreateSemaphoreW(nil, 2, 2, "Local\\QuickConvertHeavy")` + `WaitForSingleObject(INFINITE)`, and release calls `ReleaseSemaphore`. Elsewhere it does nothing.
  - `func handlerFor(a actions.Action, ext string) app.Handler`, in package main: `clean`, `shrink`, `md`, `code-*`/`txt`, and `png`/`jpg` go to `pdf.Render` when ext is `.pdf` and to `media.Image` otherwise.
  - `main`:
    - `quickconvert <action> <file>` calls `os.Exit(app.Run(...))` with StartDelay 1 s, Log = the log file (opened for append and truncated above 1 MB), and `Eng: engines.Default()`.
    - Wrong arg count notifies `"QuickConvert"`/`"Usage: quickconvert <action> <file>"` and exits 2.
    - `--version` prints the version.

- [ ] **Step 1: Write failing tests:**
  - `TestToastScriptUsesEnvOnly`: ToastScript contains `$env:QC_TITLE` and `$env:QC_BODY`, and contains no `%s` or `{0}`.
  - `TestHandlerForRoutesPdfImages`: `handlerFor(png, ".pdf")` and `handlerFor(png, ".jpg")` are different functions. Compare with `reflect.ValueOf(h).Pointer()`.
  - `TestHandlerForAllMenuItems`: for every ext in `AllExts()` and every Submenu/Extras item, `handlerFor` is non-nil.
- [ ] **Step 2:** Run `go test ./internal/platform/ ./cmd/...`. Expected: FAIL.
- [ ] **Step 3:** Implement.
- [ ] **Step 4:** Run `go test ./... && GOOS=windows go build -ldflags "-H windowsgui" -o /tmp/qc.exe ./cmd/quickconvert && file /tmp/qc.exe`. Expected: PASS, and the output contains `PE32+ executable (GUI) x86-64`.
- [ ] **Step 5:** Commit: `feat: notifications, heavy-job limit, main entrypoint`

---

### Task 5: Video, audio and image conversion (ffmpeg)

**Files:**
- Create: `internal/media/media.go`, `internal/media/media_test.go`
- Modify: `cmd/quickconvert/dispatch.go` (wire the handlers)
- Modify: `internal/testutil/testutil.go` (fixture generators)

**Interfaces:**
- Consumes: `app.Job`, `engines.Set.Run`, and `outpath.Work.Commit`.
- Produces:
  - `type Info struct{ VideoCodec, AudioCodec string; HasVideo, HasAudio bool }`
  - `func Probe(ctx context.Context, eng *engines.Set, dir, in string) (Info, error)`: uses `ffprobe -v error -show_streams -of json`. Takes the first video stream and the first audio stream. Ignores `attached_pic` video streams (cover art).
  - `func MP4Args(in, out string, i Info) []string`:
    - Remux, when the video is h264/hevc and the audio is aac/mp3/none: `-map 0:v:0 -map 0:a? -c copy -movflags +faststart`, plus `-tag:v hvc1` for hevc.
    - Otherwise: `-map 0:v:0 -map 0:a? -c:v libx264 -preset medium -crf 20 -pix_fmt yuv420p -c:a aac -b:a 192k -movflags +faststart`.
    - Always prefix `-hide_banner -nostdin -y -i in` and add `-map_metadata 0`.
  - `func MP3Args(in, out string) []string`: `-vn -map 0:a:0 -c:a libmp3lame -b:a 192k`.
  - `func WAVArgs(in, out string) []string`: `-vn -map 0:a:0 -c:a pcm_s16le`.
  - `func GIFArgs(in, out string) []string`: `-vf "fps=15,scale='min(480,iw)':-1:flags=lanczos,split[a][b];[a]palettegen[p];[b][p]paletteuse" -loop 0`.
  - `func ImageArgs(in, out, format string, orientation int) []string`:
    - Always `-noautorotate` before `-i`, then `-frames:v 1 -map_metadata -1`.
    - Orientation filter: 2=`hflip`, 3=`hflip,vflip`, 4=`vflip`, 5=`transpose=0`, 6=`transpose=1`, 7=`transpose=3`, 8=`transpose=2`.
    - jpg: `-q:v 2` and transparency flattened onto white. The filter choice is the implementer's, pinned by the test.
    - webp: `-c:v libwebp -quality 90`.
    - png: defaults.
  - `func Orientation(ctx, eng, dir, in string) int`: `exiftool -n -s3 -Orientation in`. Returns 1 on any error or empty output.
  - Handlers:
    - `func Video(ctx, j app.Job) (app.Result, error)`, for mp4/mp3/gif. For mp3 when `!HasAudio`, it returns `UserError{"This video has no audio"}`.
    - `func Audio(ctx, j app.Job) (app.Result, error)`, for mp3/wav.
    - `func Image(ctx, j app.Job) (app.Result, error)`, for png/jpg/webp.
  - Every handler writes to `out<OutExt>` inside `j.Work.Dir` (relative), then calls `j.Work.Commit("out"+ext, j.Base+j.Action.Tag, ext)`.
- `testutil` generators, all created with the Linux engines:
  - `MakeVideo(t, dir, name, vcodec string, withAudio bool)`: uses lavfi `testsrc=d=1:s=320x240` + `sine`.
  - `MakeVideoWithSubs(t, dir, name)`: an h264 + aac + mov_text subtitle MOV.
  - `MakeImage(t, dir, name string, w, h int, transparent bool)`.
  - `SetOrientation(t, path, n)`: uses exiftool.
  - `Dims(t, path) (w, h int)`: uses ffprobe.

- [ ] **Step 1: Write failing tests** (`testutil.Need(t,"ffmpeg","ffprobe")`; tests that use orientation also need exiftool):
  - `TestMP4ArgsRemuxHEVC`: the args contain `-c copy` and `-tag:v hvc1`, and do not contain `libx264`.
  - `TestMP4ArgsReencodeVP9`: the args contain `libx264` and `192k`.
  - `TestVideoMovToMp4Remux`: an h264+aac `clip.mov` gives `clip.mp4` next to it, and its probe shows h264 + aac.
  - `TestVideoWithSubtitleTrack`: `MakeVideoWithSubs` → mp4 succeeds, and the output has no subtitle stream.
  - `TestVideoToMp3`: the output's audio codec is `mp3`. Bitrate is not asserted.
  - `TestVideoNoAudioToMp3`: gives the UserError `"This video has no audio"`, and no `clip.mp3` exists.
  - `TestVideoToGifWidthCap`: a 1280x720 input gives a GIF 480 px wide.
  - `TestAudioWavToMp3`.
  - `TestImageRotated`: a 40x20 JPEG with Orientation=6 → png gives dimensions 20x40.
  - `TestImageTransparentToJpgWhite`: a fully transparent 16x16 PNG → jpg, then `ffmpeg -i out.jpg -f rawvideo -pix_fmt rgb24 -` gives a first pixel with every channel ≥ 245.
  - `TestImageToWebp`: the output codec is `webp`.
  - `TestUnicodeFolderEndToEnd`: `MakeVideo` in folder `مجلد 'أ'`, named `رحلة صيف.mov` → mp3 gives `رحلة صيف.mp3`.
- [ ] **Step 2:** Run `go test ./internal/media/`. Expected: FAIL.
- [ ] **Step 3:** Implement, and wire Video/Audio/Image in `dispatch.go`.
- [ ] **Step 4:** Run `go test ./internal/media/ ./cmd/...`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: video, audio and image conversion`

---

### Task 6: Photo metadata removal + TXT ↔ code

**Files:**
- Create: `internal/photo/clean.go`, `internal/photo/clean_test.go`, `internal/textcopy/textcopy.go`, `internal/textcopy/textcopy_test.go`
- Modify: `cmd/quickconvert/dispatch.go`

**Interfaces:**
- Consumes: `app.Job`, `engines.Set.Run`, `outpath.Work.Commit`.
- Produces:
  - `func photo.Clean(ctx, j app.Job) (app.Result, error)`: runs `exiftool -q -q -charset filename=utf8 -all= -tagsFromFile @ -Orientation -ICC_Profile -o out<ext> in<ext>` in `j.Work.Dir`, then commits as `j.Base+" (clean)"`.
  - `func textcopy.Copy(ctx, j app.Job) (app.Result, error)`: copies `j.Work.In` byte-for-byte to `out<OutExt>`, then commits as `j.Base`.

- [ ] **Step 1: Write failing tests:**
  - `TestCleanRemovesGPSAndCamera`: a JPEG tagged with `-GPSLatitude=25.2 -GPSLatitudeRef=N -Make=Apple -Model="iPhone 15" -DateTimeOriginal="2026:01:01 10:00:00" -Orientation#=6`, plus an embedded ICC profile copied from a stock sRGB profile:
    - Generate one with `exiftool -icc_profile<=` from `/usr/share/color/icc` if present. Otherwise skip that one assertion.
    - In `photo (clean).jpg`, `exiftool -GPSLatitude -Make -Model -DateTimeOriginal` gives empty values, and `Orientation#` is `6`.
    - `ICC_Profile` is present when the input had one.
    - Decoded pixels match: the `ffmpeg -noautorotate -i X -f md5 -` output is the same for input and output.
    - The original still has its GPS tag.
  - `TestCleanPNG`: a PNG tagged with `-Make=Test` gives a clean output without Make.
  - `TestCopyIdenticalBytes`: `notes.txt` with content `"print('hi')\r\n"` → `code-py` gives `notes.py` with identical bytes, and `notes.txt` is unchanged.
  - `TestCopyCollision`: with `notes.py` already existing, the result is `notes (1).py`.
  - `TestCopyCodeToTxt`: `script.m` → `txt` gives `script.txt`.
- [ ] **Step 2:** Run `go test ./internal/photo/ ./internal/textcopy/`. Expected: FAIL.
- [ ] **Step 3:** Implement, and wire `clean`, `code-*` and `txt`.
- [ ] **Step 4:** Run the same command. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: photo metadata removal and text/code copy`

---

### Task 7: PDF → PNG/JPG + PDF size reduction

**Files:**
- Create: `internal/pdf/render.go`, `internal/pdf/shrink.go`, `internal/pdf/pdf_test.go`
- Modify: `internal/testutil/testutil.go`, `cmd/quickconvert/dispatch.go`

**Interfaces:**
- Consumes: `app.Job`, `app.UserError`, `app.Info`, the engines, and `outpath.Work.Commit`/`CommitPages`.
- Produces:
  - `func Render(ctx, j app.Job) (app.Result, error)`:
    - Runs `mutool draw -q -r 150 -F png -o page-%d.png in.pdf`.
    - Collects the `page-N.png` files sorted by N.
    - For `jpg`, converts each page with `ffmpeg -i page-N.png -q:v 2 page-N.jpg`.
    - One page: `Commit(…, j.Base, ext)`. Several pages: `CommitPages`.
    - If mutool's stderr contains `password` (case-insensitive), returns `UserError{"This PDF is password-protected"}`.
  - `func Shrink(ctx, j app.Job) (app.Result, error)`:
    - Runs `gs -q -dSAFER -dBATCH -dNOPAUSE -sDEVICE=pdfwrite -dCompatibilityLevel=1.5 -dPDFSETTINGS=/ebook -sOutputFile=out.pdf in.pdf`.
    - If `size(out) > 0.9*size(in)`, returns `Info{"Already optimized — no smaller file made"}`.
    - Otherwise commits `j.Base+" (small)"` and returns `app.Result{Outputs, Detail: FormatSize(in)+" → "+FormatSize(out)}`. The success notification body becomes `"Done: report (small).pdf (4.2 MB → 0.9 MB)"`.
    - Password-protected PDFs get the same UserError as Render.
  - `func FormatSize(n int64) string`: one decimal place in MB at or above 1 MiB (`"4.2 MB"`), otherwise whole KB (`"850 KB"`).
  - `testutil`:
    - `MakePDF(t, dir, name string, pages int)`: uses gs from PostScript, with text + the page number on each page.
    - `MakeImagePDF(t, dir, name)`: one page carrying a large noisy embedded image (ffmpeg `noise` filter → JPEG at 300 DPI, then `gs` `viewjpeg.ps`, or a `mutool create` equivalent), so it can shrink.
    - `MakeEncryptedPDF(t, dir, name)`: uses `gs … -sOwnerPassword=o -sUserPassword=u`.

- [ ] **Step 1: Write failing tests:**
  - `TestFormatSize`: `4404019` gives `"4.2 MB"` and `870400` gives `"850 KB"`.
  - `TestRenderSinglePage`: gives `doc.png`, 1240 px wide (A4 at 150 DPI = 1240 ± 2).
  - `TestRenderMultiPage`: a 3-page PDF gives `doc-1.png`, `doc-2.png` and `doc-3.png`.
  - `TestRenderJpg`: a 2-page PDF gives `doc-1.jpg` and `doc-2.jpg`, with codec `mjpeg`.
  - `TestRenderEncrypted`: gives the UserError `"This PDF is password-protected"`, and no `.qc-*.part` folder is left (run through `app.Run`).
  - `TestShrinkImagePDF`: the output exists and is smaller than 0.9× the input, and `Result.Detail` contains `"→"`.
  - `TestShrinkTextPDFAlreadyOptimal`: a `MakePDF` 1-page input gives `Info`, and no `(small)` file exists.
- [ ] **Step 2:** Run `go test ./internal/pdf/`. Expected: FAIL.
- [ ] **Step 3:** Implement, then wire `png`/`jpg` for `.pdf` and `shrink`.
- [ ] **Step 4:** Run `go test ./internal/... ./cmd/...`. Expected: PASS.
- [ ] **Step 5:** Commit: `feat: PDF render and shrink`

---

### Task 8: PDF → Markdown

**Files:**
- Create: `internal/pdfmd/stext.go`, `internal/pdfmd/markdown.go`, `internal/pdfmd/pdfmd.go`, `internal/pdfmd/*_test.go`, `internal/pdfmd/testdata/*.xml`
- Modify: `cmd/quickconvert/dispatch.go`

**Interfaces:**
- Consumes: `app.Job`, the engines, `outpath.Work.Commit`, `app.UserError`.
- Produces:
  - `type Span struct{ Text string; Bold, Italic bool }`
  - `type Line struct{ Page, Block int; PageH, X0, X1, Y, Size float64; Spans []Span }`. `Y` is the baseline of the first char, and `Size` is the font size weighted by character count.
  - `func ParseStext(r io.Reader) ([]Line, error)`: parses mutool's `-F stext` XML (`document>page[@height]>block>line>font[@name,@size]>char[@x,@y,@c]`) and ignores `<image>`. A font name containing bold/black/heavy/semibold/demi counts as Bold, and one containing italic/oblique counts as Italic (both case-insensitive). Adjacent chars with the same style merge into one Span.
  - `func ToMarkdown(lines []Line) (string, error)`: returns `ErrNoText` when there are 0 non-space chars.
  - `var ErrNoText = errors.New("This PDF has no selectable text (scanned?)")`
  - `func Convert(ctx, j app.Job) (app.Result, error)`: runs `mutool draw -q -F stext -o out.xml in.pdf`, then ParseStext and ToMarkdown, writes `out.md` as UTF-8 without a BOM with a trailing `\n`, and commits it as `j.Base`. ErrNoText becomes `UserError{ErrNoText.Error()}`.

`ToMarkdown` rules. The values are exact; the algorithm is spelled out because tests alone don't pin it.
1. **Body size:** the most common `Size` rounded to 0.5, weighted by character count.
2. **Header/footer removal:** a line is in the band when `Y < 0.08*PageH` or `Y > 0.92*PageH`. Its key is the text lower-cased and trimmed, with every digit run replaced by `#`. Drop band lines whose key appears on ≥ 50% of pages when the document has ≥ 3 pages. Always drop band lines whose key matches `^(page\s*)?#(\s*(of|/)\s*#)?$`.
3. **Tables (per page):**
   - Group lines whose `Y` values are within 2 pt into rows. A row with ≥ 2 lines is a candidate, with its cells sorted by X0.
   - ≥ 2 consecutive candidate rows with the same cell count, and every cell's X0 within 6 pt of the first row's, form a table.
   - Emit `| c1 | c2 |`, then `| --- | --- |` after the first row. Cells are plain text with `|` escaped as `\|`.
   - The table is emitted at the position of its first line in reading order, and its other lines are skipped.
4. **Headings:**
   - Collect the distinct sizes ≥ `1.15 × body` (rounded to 0.5) and sort them in descending order. The 1st is `#`, the 2nd `##`, and the 3rd and below `###`.
   - Consecutive heading lines in the same block with the same size join into one heading, separated by spaces.
   - Heading text is plain (no emphasis markers).
5. **Lists:**
   - A line starting with one of `• ◦ ▪ ▫ ■ □ ● ○ ‣ ⁃ ∙ · - – *` followed by whitespace becomes `- rest`.
   - `^\d{1,3}[.)]\s+` becomes `N. rest`.
   - A line whose whole text is one bullet glyph merges with the next line in reading order when that line is within 2 pt in Y.
   - Lines that continue a list item (same block, not themselves a list start, X0 > the item's X0 + 3) join the item.
6. **Paragraphs:**
   - Lines in one block join into one paragraph. When the previous line ends in letter+`-` and the next starts with a lower-case letter, drop the hyphen and add no space. Otherwise join with a single space.
   - Blocks are separated by a blank line. List items in a run are separated by a single newline.
7. **Emphasis:** a Bold span is wrapped in `**…**`, Italic in `*…*`, and both in `***…***`. Whitespace stays outside the markers, and whitespace-only spans get no markers. A literal `*` in the text becomes `\*`. A `#` or `>` at the start of a paragraph is escaped with a backslash.
8. **Output:** collapse three or more newlines into two, trim, and end with a single `\n`.

- [ ] **Step 1: Write the XML fixtures** in `testdata/`, each a small handwritten stext document:
  - `headings.xml`: a 24 pt title, a 16 pt subtitle, and 11 pt body lines in 2 blocks, one line ending in `"exam-"` followed by `"ple text"`.
  - `lists.xml`: bullet lines `"• Apples"` and `"• Pears"`, a numbered `"2) Second"`, and a separate-glyph bullet (`"•"` at x=72 and `"Item"` at x=86 on the same Y).
  - `table.xml`: 3 rows × 3 cells aligned at x=72/200/330, the header row being `Name, Qty, Price`.
  - `headerfooter.xml`: 4 pages, each with `"ACME Report"` at Y=30, `"Page N of 4"` at Y=820 (PageH 842) and one body line.
  - `styles.xml`: one line made of spans `"This is "` (Helvetica), `"bold"` (Helvetica-Bold), `" and "` and `"italic"` (Helvetica-Oblique), plus `"2*3"`.
  - `empty.xml`: 2 pages with no chars.
- [ ] **Step 2: Write failing tests:**
  - `TestHeadings`: the output starts with `"# Title\n\n## Subtitle\n\n"` and contains `"example text"`.
  - `TestLists`: contains `"- Apples\n- Pears"`, `"2. Second"` and `"- Item"`.
  - `TestTable`: contains `"| Name | Qty | Price |\n| --- | --- | --- |\n"` followed by two data rows.
  - `TestHeaderFooterRemoved`: contains neither `"ACME Report"` nor `"Page"`, and contains all 4 body lines.
  - `TestStyles`: contains `"This is **bold** and *italic*"` and `"2\*3"`.
  - `TestEmpty`: returns `ErrNoText`.
  - `TestParseStextFontFlags`: the font name `"ABCDEF+Arial-BoldItalicMT"` gives Bold and Italic.
  - `TestConvertEndToEnd` (needs gs + mutool): a PostScript-built PDF with a 20 pt `"Big Title"` and 11 pt body becomes `doc.md` starting with `"# Big Title"`.
  - `TestConvertScanned` (needs gs, mutool, ffmpeg): an image-only PDF from `MakeImagePDF` gives the UserError `"This PDF has no selectable text (scanned?)"`, and no `.md` is created.
- [ ] **Step 3:** Run `go test ./internal/pdfmd/`. Expected: FAIL.
- [ ] **Step 4:** Implement `stext.go`, then `markdown.go`, then `pdfmd.go`, and wire `md`.
- [ ] **Step 5:** Run `go test ./internal/pdfmd/ ./cmd/...`. Expected: PASS.
- [ ] **Step 6:** Commit: `feat: PDF to Markdown`

---

### Task 9: Install / uninstall script generation

**Files:**
- Create: `internal/menu/menu.go`, `internal/menu/menu_test.go`, `cmd/genscripts/main.go`

**Interfaces:**
- Consumes: `actions.AllExts`, `actions.Submenu`, `actions.Extras`.
- Produces:
  - `func InstallScript() string`
  - `func UninstallScript() string`
  - Both use CRLF line endings throughout. `genscripts <dir>` writes `<dir>/install.cmd` and `<dir>/uninstall.cmd`.

Registry layout (cascade via an empty `SubCommands` value + a nested `shell` key; all paths under `HKCU\Software\Classes\SystemFileAssociations\<ext>\shell\`):
- `QuickConvert`: `MUIVerb="Convert to"`, `SubCommands=""`.
- `QuickConvert\shell\NN<id>`, where NN is a 2-digit menu order starting at 01: `MUIVerb=<Label>`, `MultiSelectModel="Player"`.
- `QuickConvert\shell\NN<id>\command`: default value `"<DEST>\quickconvert.exe" <id> "%1"`.
- `QuickConvert.<id>` for extras: `MUIVerb=<Label>` and `MultiSelectModel="Player"`, with a `\command` key like the one above.

`install.cmd`:
1. `@echo off`, `setlocal`, `set "DEST=%LOCALAPPDATA%\QuickConvert"`.
2. Run the full uninstall `reg delete` list (so stale entries go).
3. `robocopy "%~dp0." "%DEST%" /E /NFL /NDL /NJH /NJS /NP >nul`. Skip this step when `"%~dp0"=="%DEST%\"`.
4. All the `reg add … /f >nul` lines. In the batch text the command data is written as `/d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\""`.
5. Finish with `echo QuickConvert installed. Right-click a file ^> Show more options ^> Convert to.` and `pause`.

`uninstall.cmd`:
1. Run `reg delete "…\<ext>\shell\QuickConvert" /f >nul 2>&1` (and the same for each `QuickConvert.<id>`) for every ext.
2. `echo` that the menu entries were removed.
3. Remove the install folder last: `(goto) 2>nul & rmdir /S /Q "%LOCALAPPDATA%\QuickConvert"`. This avoids the "batch file cannot be found" error when the script is run from inside the folder it deletes.

- [ ] **Step 1: Write failing tests:**
  - `TestCRLF`: every `\n` in both scripts is preceded by `\r`.
  - `TestInstallHasMovMp4`: contains `reg add "HKCU\Software\Classes\SystemFileAssociations\.mov\shell\QuickConvert\shell\01mp4\command" /ve /d "\"%DEST%\quickconvert.exe\" mp4 \"%%1\"" /f >nul`.
  - `TestInstallNoBarePercent1`: the text contains no `"%1"` that isn't preceded by `%`.
  - `TestInstallExtras`: contains `…\.pdf\shell\QuickConvert.shrink" /v MUIVerb /d "Reduce size"` and `…\.heic\shell\QuickConvert.clean" /v MUIVerb /d "Remove metadata"`.
  - `TestEveryMenuItemRegistered`: for every ext and every Submenu/Extras item, a matching `command` line exists.
  - `TestUninstallCoversInstall`: every `…\shell\QuickConvert…` key created by install appears as a `reg delete` key in uninstall.
- [ ] **Step 2:** Run `go test ./internal/menu/`. Expected: FAIL.
- [ ] **Step 3:** Implement.
- [ ] **Step 4:** Run `go test ./internal/menu/ && go run ./cmd/genscripts /tmp/qcs && head -5 /tmp/qcs/install.cmd`. Expected: PASS, and the first line is `@echo off`.
- [ ] **Step 5:** Commit: `feat: generate install/uninstall scripts`

---

### Task 10: Windows engines, build and package

**Files:**
- Create: `scripts/fetch-engines.sh`, `scripts/build.sh`, `packaging/README.txt`

**Interfaces:**
- Consumes: everything above.
- Produces:
  - `dist/QuickConvert.zip` containing a `QuickConvert/` folder with:
    - `quickconvert.exe`, `ffmpeg.exe`, `ffprobe.exe`
    - `mutool.exe`
    - `gswin64c.exe`, `gsdll64.dll`
    - `exiftool.exe` + `exiftool_files/`
    - `install.cmd`, `uninstall.cmd`, `README.txt`

Steps:

- [ ] **Step 1:** `fetch-engines.sh` downloads pinned versions into `engines/windows/` and extracts the files listed above:
  - ffmpeg: the gyan.dev "release essentials" zip.
  - MuPDF: the Windows zip from mupdf.com.
  - Ghostscript: the `gs*w64.exe` NSIS installer from the ArtifexSoftware/ghostpdl-downloads GitHub release, extracted with `7z x` (install `p7zip-full`).
  - ExifTool: the `exiftool-*_64.zip` from exiftool.org. Rename `exiftool(-k).exe` to `exiftool.exe`.

  The script must be idempotent and verify each file exists after extraction.

  **Known blocker:** this container's network currently refuses gyan.dev, mupdf.com, exiftool.org and github.com release downloads. When the script hits that, call `get_environment` and propose allowing exactly the hosts the script needs. If the owner declines, `packaging/README.txt` gets a "Getting the engines yourself" section with the four download links and where to put each file. Then build the zip without the engines and tell the owner.
- [ ] **Step 2:** `build.sh`:
  1. `go test ./...`
  2. `GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-H windowsgui -s -w -X main.version=0.1.0" -o dist/QuickConvert/quickconvert.exe ./cmd/quickconvert`
  3. `go run ./cmd/genscripts dist/QuickConvert`
  4. Copy the engines and `README.txt` in.
  5. Run `cd dist && zip -r QuickConvert.zip QuickConvert`.
- [ ] **Step 3:** Write `packaging/README.txt` in plain English. Sections:
  - Install: unzip anywhere, double-click `install.cmd`.
  - Use: right-click a file → Show more options (or Shift+right-click) → Convert to / Remove metadata / Reduce size.
  - Where results go: the same folder.
  - What each option does.
  - Uninstall: run `uninstall.cmd`.
  - The log file location.
  - The note that notifications show "Windows PowerShell" as the sender.
- [ ] **Step 4:** Run `scripts/build.sh && unzip -l dist/QuickConvert.zip`. Expected: tests pass and the listing shows every file above (or every file except the engines, in the declined case).
- [ ] **Step 5:** Commit `scripts/`, `packaging/` (never `engines/` or `dist/`): `build: Windows packaging`
- [ ] **Step 6: Owner smoke test on Windows.** Send the owner a short checklist:
  1. Install.
  2. On an iPhone `.mov`: MP4 and MP3.
  3. On a phone photo: PNG and Remove metadata.
  4. On a scanned PDF and a text PDF: PNG, MD and Reduce size.
  5. On `notes.txt`: `.py`.
  6. Select 20 files at once and convert.
  7. Uninstall.

  Also flag the HEIC input risk: if Windows ffmpeg can't decode a given `.heic` file, the owner reports it and we revisit.
