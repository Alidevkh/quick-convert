QuickConvert
============

Convert files from the right-click menu. The result is always saved next to
the original file, and the original is never changed or deleted.


INSTALL
-------
1. Unzip this folder anywhere (for example, your Downloads folder).
2. Double-click install.cmd.
   No administrator rights are needed. The program is copied to
   %LOCALAPPDATA%\QuickConvert


USE
---
Right-click a file. On Windows 11, click "Show more options" first
(or hold Shift while right-clicking). Then choose:

  Convert to >     the formats that make sense for that file
  Remove metadata  (photos, phone videos, music) removes GPS location,
                   dates, camera and titles; saves "photo (clean).jpg",
                   "trip (clean).mov"... Music cover pictures are kept.
  Reduce size      (PDFs) saves "report (small).pdf"
  Split into pages (PDFs) every page as its own PDF: report-1.pdf, report-2.pdf...
  Extract pages... (PDFs) asks which pages (e.g. 1-3, 5) and saves them, in that
                   order, as one PDF: "report (pages 1-3, 5).pdf"

A Windows notification tells you when it starts and when it is done.
Notifications show "Windows PowerShell" as the sender; that is normal.

Combine several files into one, from the "Send to" menu:
Select the files, right-click > Show more options > Send to >

  Images to PDF    all selected pictures in one PDF, one picture per page,
                   saved as "<first picture name>.pdf"
  Merge PDFs       all selected PDFs in one, saved as "<first name> (merged).pdf"

Files are put in file-name order (IMG_9 before IMG_10). To choose another
order, rename the files, e.g. "1 cover.pdf", "2 report.pdf".

You can select many files at once. At most 2 big video or PDF jobs (and 4
small jobs) run at the same time; the rest wait their turn.

Nothing runs in the background: QuickConvert starts when you click an
entry and closes when the file is done.


WHAT EACH OPTION DOES
---------------------
Videos (.mov .mp4 .m4v .mkv .avi .webm .wmv .3gp .mts .m2ts .mpg .mpeg .flv)
  MP4 (full quality, for editing)
                    For editing programs (Premiere, CapCut, DaVinci...).
                    iPhone and most phone videos are repackaged in seconds,
                    exactly the same picture and sound: nothing is
                    re-compressed, and HDR colours, every sound track, the
                    location, date and phone model are all kept.
                    Videos MP4 can't hold as they are are re-compressed at a
                    quality that looks the same as the original.
  MP4 (WhatsApp / any device)
                    Plays on every phone, computer and browser, ready to send
                    on WhatsApp: H.264, at most 1080p, iPhone HDR colours
                    converted so they look right everywhere. Slower than MP4
                    (the video is re-compressed). Saved as "trip (WhatsApp).mp4",
                    without location or other hidden information. Camcorder
                    and DVD videos are also de-interlaced (no comb lines).
  MP3 (audio only)  The sound track, 192 kbps.
  M4A (audio only)  The sound track as AAC (iPhone sound is copied, no loss).
  GIF               Animated, 15 frames per second, at most 480 pixels wide.

Audio (.mp3 .wav .m4a .flac .ogg .aac .wma .opus .aiff .aif .amr)
  MP3               192 kbps, plays everywhere
  M4A               AAC 192 kbps (Apple / iPhone)
  WAV               uncompressed
  FLAC              lossless, about half the size of WAV
  OGG               Vorbis, about 192 kbps

Pictures (.png .jpg .jpeg .webp .bmp .heic .gif .tif .tiff .avif)
  PNG, JPG, WEBP    Phone photos come out the right way up.
                    Transparent areas become white in JPG.
  PDF               A one-page PDF, the page the size of the picture.
                    (To put several pictures in ONE PDF, use Send to > Images to PDF.)
  ICO (icon)        A Windows icon (256, 48, 32 and 16 pixels). Pictures that
                    are not square get a transparent border; nothing is cut.
  MP4 (video)       (GIFs only) The animation as a video, which plays in
                    every app. Transparent parts become white.
  TIFF scans with several pages: PDF keeps every page; PNG, JPG and WEBP
  use the first page. Animated GIFs: PNG, JPG, WEBP and PDF use the first
  picture.
  Remove metadata   Removes GPS location, camera, dates and other hidden
                    information. Keeps the picture quality, the rotation and
                    the colours exactly the same.

SVG drawings (.svg)
  PNG               2048 pixels on the longest side, transparent background
                    kept.

E-books (.epub)
  PDF               The book as a PDF, text still selectable.
  MD                The text as Markdown, for AI tools.
  Protected (DRM) e-books from online stores can't be converted.

PDF
  PNG, JPG          One picture per page (report-1.png, report-2.png...).
  MD                The text as Markdown, for AI tools: headings, bold,
                    lists and simple tables are kept; page numbers and
                    repeated headers are removed. Scanned PDFs (pictures of
                    text) are not supported.
  Reduce size       Shrinks the pictures inside the PDF. If the file can't
                    get at least 10% smaller, no copy is made.

Text and code
  .txt -> .py .m .c .cpp .h .js .ts .html .css .java .json .md .bat .sh
  code files -> .txt
                    The content stays exactly the same; only the type changes.


IF SOMETHING GOES WRONG
-----------------------
"Windows protected your PC" (SmartScreen) when using Send to:
QuickConvert is not signed with a paid certificate. install.cmd removes the
"downloaded from the internet" mark from its own files, so this should not
appear. If it does: right-click %LOCALAPPDATA%\QuickConvert\quickconvert.exe
> Properties > tick "Unblock" > OK.

The notification gives the reason. Details are written to:
  %LOCALAPPDATA%\QuickConvert\log.txt


TEST ON THIS COMPUTER
--------------------
Double-click selftest.cmd (in this folder). It installs QuickConvert, makes
its own sample files in a temporary folder, tries every kind of conversion,
checks that nothing is left running, and saves selftest-report.txt here.
Many notifications pop up while it runs; that is normal. Your own files are
not used.


UNINSTALL
---------
Double-click uninstall.cmd (in this folder, or in %LOCALAPPDATA%\QuickConvert).


GETTING THE ENGINES YOURSELF
----------------------------
QuickConvert uses four free programs ("engines"). If this folder does not
already contain ffmpeg.exe, mutool.exe, gswin64c.exe and exiftool.exe,
download them and put the files listed below in this folder, next to
quickconvert.exe, BEFORE running install.cmd:

1. FFmpeg     https://www.gyan.dev/ffmpeg/builds/  -> "ffmpeg-release-essentials.zip"
              Copy bin\ffmpeg.exe and bin\ffprobe.exe
2. MuPDF      https://mupdf.com/releases  -> the Windows zip
              Copy mutool.exe
3. Ghostscript https://ghostscript.com/releases/gsdnld.html  -> 64-bit Windows
              Simply install it normally; QuickConvert finds it by itself.
4. ExifTool   https://exiftool.org/  -> "Windows Executable" (64-bit zip)
              Copy "exiftool(-k).exe" and rename it to exiftool.exe,
              and copy the exiftool_files folder next to it.
