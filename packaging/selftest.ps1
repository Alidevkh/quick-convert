# QuickConvert self-test for Windows.
#
# Run selftest.cmd (next to this file) from the unzipped QuickConvert folder.
# It installs QuickConvert, makes sample files, runs every kind of
# conversion, checks that nothing is left running or left behind, measures
# time and memory, and writes selftest-report.txt next to this file.
#
# Many "Done" notifications will pop up while it runs; that is expected.
# Your own files are never used or touched: everything happens in a new
# folder under %TEMP%, deleted at the end.
#
#   selftest.cmd              install + test (QuickConvert stays installed)
#   selftest.cmd -Uninstall   also test uninstall.cmd at the end

param([switch]$Uninstall)

$ErrorActionPreference = 'Stop'
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
$dest = Join-Path $env:LOCALAPPDATA 'QuickConvert'
$exe = Join-Path $dest 'quickconvert.exe'
$report = New-Object System.Collections.Generic.List[string]
$script:failures = 0

function Log([string]$line) { Write-Host $line; $report.Add($line) }
function Pass([string]$what) { Log ("  OK    " + $what) }
function Fail([string]$what) { Log ("  FAIL  " + $what); $script:failures++ }
function Check([bool]$ok, [string]$what) { if ($ok) { Pass $what } else { Fail $what } }

# Engine processes started from the QuickConvert folder (not other copies
# of ffmpeg the user may have).
function Engines {
    $names = 'ffmpeg', 'ffprobe', 'mutool', 'gswin64c', 'exiftool', 'perl'
    @(Get-Process -Name $names -ErrorAction SilentlyContinue | Where-Object {
        $_.Path -and ($_.Path.StartsWith($dest, [StringComparison]::OrdinalIgnoreCase) -or $_.Name -eq 'gswin64c')
    })
}
function FFmpegs { @(Engines | Where-Object { $_.Name -eq 'ffmpeg' }) }
function QCProcesses { @(Get-Process -Name 'quickconvert' -ErrorAction SilentlyContinue) }
function Leftovers([string]$dir) { @(Get-ChildItem -LiteralPath $dir -Force -Directory -Filter '.qc-*' -ErrorAction SilentlyContinue) }

# Runs quickconvert.exe and waits; returns exit code, seconds and the peak
# memory (MB) of QuickConvert plus its engines.
function RunQC([string[]]$arguments) {
    $quoted = $arguments | ForEach-Object { '"' + $_ + '"' }
    $start = Get-Date
    $p = Start-Process -FilePath $exe -ArgumentList $quoted -PassThru
    $null = $p.Handle # without this, PowerShell 5.1 may lose the exit code
    $peak = 0
    while (-not $p.HasExited) {
        $mem = 0
        foreach ($e in @(Engines) + @(QCProcesses)) { try { $mem += $e.WorkingSet64 } catch {} }
        if ($mem -gt $peak) { $peak = $mem }
        Start-Sleep -Milliseconds 100
    }
    $p.WaitForExit()
    [pscustomobject]@{ Code = $p.ExitCode; Seconds = ((Get-Date) - $start).TotalSeconds; PeakMB = [math]::Round($peak / 1MB) }
}

function Tool([string]$name) { Join-Path $dest $name }
function Make([string]$tool, [string[]]$arguments) {
    $ErrorActionPreference = 'Continue' # engine messages on stderr are not errors
    & (Tool $tool) @arguments 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) { throw "$tool failed: $arguments" }
}

# ---------------------------------------------------------------- system
Log "QuickConvert self-test  $(Get-Date -Format 'yyyy-MM-dd HH:mm')"
$os = Get-CimInstance Win32_OperatingSystem
$cpu = Get-CimInstance Win32_Processor | Select-Object -First 1
Log ("Windows : " + $os.Caption + " " + $os.Version + " (" + $os.OSArchitecture + ")")
Log ("CPU     : " + $cpu.Name.Trim() + ", " + $cpu.NumberOfLogicalProcessors + " threads")
Log ("Memory  : " + [math]::Round($os.TotalVisibleMemorySize / 1MB, 1) + " GB")
Log ("PowerShell " + $PSVersionTable.PSVersion + ", language mode " + $ExecutionContext.SessionState.LanguageMode)
Log ""

# ---------------------------------------------------------------- install
Log "1. Install"
if (-not (Test-Path (Join-Path $here 'install.cmd'))) { Fail "install.cmd not found next to selftest.ps1"; exit 1 }
cmd /c "echo.| `"$(Join-Path $here 'install.cmd')`"" | Out-Null
Check (Test-Path $exe) "program copied to $dest"
$size = (Get-ChildItem -LiteralPath $dest -Recurse -File | Measure-Object Length -Sum).Sum
Log ("        installed size: " + [math]::Round($size / 1MB) + " MB")
$zone = Get-Item -LiteralPath $exe -Stream Zone.Identifier -ErrorAction SilentlyContinue
Check ($null -eq $zone) "no 'downloaded from the internet' mark (no SmartScreen warning)"
$root = 'HKCU:\Software\Classes\SystemFileAssociations'
foreach ($k in '.mov\shell\QuickConvert', '.heic\shell\QuickConvert', '.pdf\shell\QuickConvert.extract', '.gif\shell\QuickConvert', '.epub\shell\QuickConvert', '.mp3\shell\QuickConvert.clean') {
    Check (Test-Path (Join-Path $root $k)) "menu entry $k"
}
Check (-not (Test-Path (Join-Path $root '.mkv\shell\QuickConvert.clean'))) "no Remove metadata on .mkv"
$sendto = [Environment]::GetFolderPath('SendTo')
foreach ($s in 'Images to PDF.lnk', 'Merge PDFs.lnk') { Check (Test-Path (Join-Path $sendto $s)) "Send to: $s" }
foreach ($e in 'ffmpeg.exe', 'ffprobe.exe', 'mutool.exe', 'exiftool.exe') { Check (Test-Path (Tool $e)) "engine $e" }
$gs = (Test-Path (Tool 'gswin64c.exe')) -or (Get-ChildItem "$env:ProgramFiles\gs\gs*\bin\gswin64c.exe" -ErrorAction SilentlyContinue)
Check ([bool]$gs) "Ghostscript (needed for Reduce size)"
Log ""

# ---------------------------------------------------------------- samples
Log "2. Sample files"
# A folder with Arabic letters and spaces, to check such names work.
$arabic = -join ([char[]](0x062A, 0x062C, 0x0631, 0x0628, 0x0629))
$work = Join-Path $env:TEMP ("QuickConvert selftest " + $arabic + " " + (Get-Random))
New-Item -ItemType Directory -Path $work | Out-Null
Push-Location $work
try {
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'testsrc2=d=3:s=1280x720:r=30', '-f', 'lavfi', '-i', 'sine=d=3',
        '-c:v', 'libx264', '-pix_fmt', 'yuv420p', '-c:a', 'aac', '-shortest', '-metadata', 'location=+25.2+055.2/', 'video.mov')
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'testsrc2=d=3:s=1920x1080:r=30', '-f', 'lavfi', '-i', 'sine=d=3',
        '-c:v', 'libx265', '-x265-params', 'log-level=error', '-pix_fmt', 'yuv420p10le', '-color_primaries', 'bt2020',
        '-color_trc', 'arib-std-b67', '-colorspace', 'bt2020nc', '-tag:v', 'hvc1', '-c:a', 'aac', '-shortest', 'iphone-hdr.mov')
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'sine=d=3', 'audio.wav')
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'testsrc2=s=1600x1200', '-frames:v', '1', '-q:v', '2', 'photo.jpg')
    Make 'exiftool.exe' @('-q', '-overwrite_original', '-GPSLatitude=25.2', '-GPSLongitude=55.2', '-Orientation#=6', 'photo.jpg')
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'testsrc2=d=1:s=200x150:r=10', 'anim.gif')
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'testsrc2=s=800x1100', '-frames:v', '1', 'scan.tif')
    Set-Content -LiteralPath 'logo.svg' -Encoding ASCII -Value '<svg xmlns="http://www.w3.org/2000/svg" width="64" height="32"><circle cx="16" cy="16" r="14" fill="red"/></svg>'
    Set-Content -LiteralPath 'notes.txt' -Encoding ASCII -Value 'print("hello")'
    # A small e-book (an EPUB is a zip with a fixed layout).
    Add-Type -AssemblyName System.IO.Compression, System.IO.Compression.FileSystem
    $zip = [System.IO.Compression.ZipFile]::Open((Join-Path $work 'book.epub'), 'Create')
    $files = [ordered]@{
        'mimetype'               = 'application/epub+zip'
        'META-INF/container.xml' = '<?xml version="1.0"?><container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container"><rootfiles><rootfile full-path="c.opf" media-type="application/oebps-package+xml"/></rootfiles></container>'
        'c.opf'                  = '<?xml version="1.0"?><package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="i"><metadata xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:identifier id="i">x</dc:identifier><dc:title>T</dc:title><dc:language>en</dc:language></metadata><manifest><item id="a" href="a.xhtml" media-type="application/xhtml+xml"/></manifest><spine><itemref idref="a"/></spine></package>'
        'a.xhtml'                = '<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml"><head><title>a</title></head><body><h1>Chapter One</h1><p>Some text.</p></body></html>'
    }
    foreach ($name in $files.Keys) {
        $level = if ($name -eq 'mimetype') { 'NoCompression' } else { 'Optimal' }
        $w = New-Object System.IO.StreamWriter($zip.CreateEntry($name, $level).Open())
        $w.Write($files[$name]); $w.Close()
    }
    $zip.Dispose()
    Pass "made in $work"
    Log ""

    # ------------------------------------------------------------ conversions
    Log "3. Conversions  (time / peak memory of QuickConvert + engine)"
    $tests = @(
        @('mp4', 'iphone-hdr.mov', 'MP4 full quality (iPhone HDR)'),
        @('whatsapp', 'iphone-hdr.mov', 'MP4 WhatsApp (HDR to normal)'),
        @('mp3', 'video.mov', 'video to MP3'),
        @('m4a', 'video.mov', 'video to M4A'),
        @('gif', 'video.mov', 'video to GIF'),
        @('clean', 'video.mov', 'Remove metadata (video)'),
        @('flac', 'audio.wav', 'audio to FLAC'),
        @('ogg', 'audio.wav', 'audio to OGG'),
        @('png', 'photo.jpg', 'photo to PNG'),
        @('webp', 'photo.jpg', 'photo to WEBP'),
        @('pdf', 'photo.jpg', 'photo to PDF'),
        @('ico', 'photo.jpg', 'photo to ICO'),
        @('clean', 'photo.jpg', 'Remove metadata (photo)'),
        @('mp4', 'anim.gif', 'GIF to MP4'),
        @('pdf', 'scan.tif', 'TIFF to PDF'),
        @('png', 'logo.svg', 'SVG to PNG'),
        @('pdf', 'book.epub', 'EPUB to PDF'),
        @('md', 'book.epub', 'EPUB to MD'),
        @('code-py', 'notes.txt', 'TXT to .py')
    )
    foreach ($t in $tests) {
        $before = @(Get-ChildItem -LiteralPath $work -File | ForEach-Object Name)
        $r = RunQC @($t[0], (Join-Path $work $t[1]))
        $new = @(Get-ChildItem -LiteralPath $work -File | Where-Object { $before -notcontains $_.Name -and $_.Length -gt 0 })
        $info = "{0,-32} {1,6:N1} s {2,6} MB" -f $t[2], $r.Seconds, $r.PeakMB
        Check ($r.Code -eq 0 -and $new.Count -gt 0) $info
    }
    # PDF features, using the PDF just made from the photo
    $pdf = Join-Path $work 'photo.pdf'
    $r = RunQC @('merge', $pdf, (Join-Path $work 'scan.pdf'))
    Check ($r.Code -eq 0 -and (Test-Path (Join-Path $work 'photo (merged).pdf'))) ("{0,-32} {1,6:N1} s {2,6} MB" -f 'Send to: Merge PDFs', $r.Seconds, $r.PeakMB)
    $merged = Join-Path $work 'photo (merged).pdf'
    foreach ($a in @(@('split', 'Split into pages'), @('png', 'PDF to PNG'), @('shrink', 'Reduce size'))) {
        $r = RunQC @($a[0], $merged)
        Check ($r.Code -eq 0) ("{0,-32} {1,6:N1} s {2,6} MB" -f $a[1], $r.Seconds, $r.PeakMB)
    }
    # PDF to MD needs a PDF with real text: the one made from the e-book.
    $r = RunQC @('md', (Join-Path $work 'book.pdf'))
    Check ($r.Code -eq 0 -and (Test-Path (Join-Path $work 'book (1).md'))) ("{0,-32} {1,6:N1} s {2,6} MB" -f 'PDF to MD', $r.Seconds, $r.PeakMB)
    # A PDF made of pictures has no text: refused with a message, no file made.
    $r = RunQC @('md', $merged)
    Check ($r.Code -ne 0 -and -not (Test-Path (Join-Path $work 'photo (merged).md'))) "PDF of pictures to MD is refused (no text)"
    $r = RunQC @('imgpdf', (Join-Path $work 'photo.jpg'), (Join-Path $work 'anim.gif'))
    Check ($r.Code -eq 0) ("{0,-32} {1,6:N1} s {2,6} MB" -f 'Send to: Images to PDF', $r.Seconds, $r.PeakMB)
    Log ""

    # ------------------------------------------------------------ nothing left
    Log "4. Nothing left behind"
    Start-Sleep -Seconds 5 # notifications: PowerShell shows them and exits
    Check ((QCProcesses).Count -eq 0) "no quickconvert.exe running"
    Check ((Engines).Count -eq 0) "no engine (ffmpeg, mutool, ...) running"
    $toasts = @(Get-CimInstance Win32_Process -Filter "Name='powershell.exe'" | Where-Object { $_.CommandLine -like '*-EncodedCommand*' })
    Check ($toasts.Count -eq 0) "no notification PowerShell still running"
    Check ((Leftovers $work).Count -eq 0) "no temporary .qc-* folder"
    Log ""

    # ------------------------------------------------------------ limits
    Log "5. Many files at once"
    $jobs = 1..6 | ForEach-Object {
        Copy-Item 'photo.jpg' "many$_.jpg"
        Start-Process -FilePath $exe -ArgumentList @('png', ('"' + (Join-Path $work "many$_.jpg") + '"')) -PassThru
    }
    $maxQC = 0
    while (@($jobs | Where-Object { -not $_.HasExited }).Count) {
        $n = (FFmpegs).Count; if ($n -gt $maxQC) { $maxQC = $n }
        Start-Sleep -Milliseconds 50
    }
    Check ($maxQC -le 4) "6 photos: at most 4 engines at once (saw $maxQC)"
    $jobs = 1..4 | ForEach-Object {
        Start-Process -FilePath $exe -ArgumentList @('whatsapp', ('"' + (Join-Path $work 'video.mov') + '"')) -PassThru
    }
    $maxHeavy = 0
    while (@($jobs | Where-Object { -not $_.HasExited }).Count) {
        $n = (FFmpegs).Count; if ($n -gt $maxHeavy) { $maxHeavy = $n }
        Start-Sleep -Milliseconds 50
    }
    Check ($maxHeavy -le 2) "4 videos: at most 2 encodings at once (saw $maxHeavy)"
    Log ""

    # ------------------------------------------------------------ killed
    Log "6. QuickConvert killed in the middle (Task Manager)"
    Make 'ffmpeg.exe' @('-v', 'error', '-f', 'lavfi', '-i', 'testsrc2=d=60:s=1920x1080:r=30', '-c:v', 'libx264', '-preset', 'ultrafast', 'long.mp4')
    $p = Start-Process -FilePath $exe -ArgumentList @('whatsapp', ('"' + (Join-Path $work 'long.mp4') + '"')) -PassThru
    $waited = 0
    while ((Engines).Count -eq 0 -and $waited -lt 100) { Start-Sleep -Milliseconds 100; $waited++ }
    Start-Sleep -Seconds 1
    Stop-Process -Id $p.Id -Force
    Start-Sleep -Seconds 2
    Check ((Engines).Count -eq 0) "ffmpeg stopped together with QuickConvert"
    Get-ChildItem -LiteralPath $work -Force -Directory -Filter '.qc-*' | Remove-Item -Recurse -Force
    Log ""
}
catch {
    Fail ("unexpected error: " + $_)
}
finally {
    Pop-Location
    Remove-Item -LiteralPath $work -Recurse -Force -ErrorAction SilentlyContinue
}

# ---------------------------------------------------------------- uninstall
if ($Uninstall) {
    Log "7. Uninstall"
    cmd /c "echo.| `"$(Join-Path $dest 'uninstall.cmd')`"" | Out-Null
    Start-Sleep -Seconds 2
    Check (-not (Test-Path $dest)) "program folder removed"
    Check (-not (Test-Path (Join-Path $root '.mov\shell\QuickConvert'))) "menu entries removed"
    Check (-not (Test-Path (Join-Path $sendto 'Merge PDFs.lnk'))) "Send to shortcuts removed"
    Log ""
}

Log ("Manual checks (need a person): right-click menus look right; 'Extract pages...' box appears in front; notifications appear.")
if ($script:failures -eq 0) { Log "RESULT: all checks passed" } else { Log "RESULT: $($script:failures) check(s) FAILED" }
if (-not (Test-Path $here)) { $here = [Environment]::GetFolderPath('Desktop') } # e.g. run from the folder just uninstalled
$out = Join-Path $here 'selftest-report.txt'
$report | Set-Content -LiteralPath $out -Encoding UTF8
Write-Host ""
Write-Host "Report saved to $out"
