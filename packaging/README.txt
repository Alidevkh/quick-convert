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
  Remove metadata  (photos) saves "photo (clean).jpg"
  Reduce size      (PDFs) saves "report (small).pdf"

A Windows notification tells you when it starts and when it is done.
Notifications show "Windows PowerShell" as the sender; that is normal.

Combine several files into one, from the "Send to" menu:
Select the files, right-click > Show more options > Send to >

  Images to PDF    all selected pictures in one PDF, one picture per page,
                   saved as "<first picture name>.pdf"
  Merge PDFs       all selected PDFs in one, saved as "<first name> (merged).pdf"

Files are put in file-name order (IMG_9 before IMG_10). To choose another
order, rename the files, e.g. "1 cover.pdf", "2 report.pdf".

You can select many files at once. At most 2 big video or PDF jobs run at
the same time; the rest wait their turn.


WHAT EACH OPTION DOES
---------------------
Videos (.mov .mp4 .mkv .avi .webm .wmv)
  MP4               If the video is already MP4-compatible (most phone videos),
                    it is repackaged in seconds without losing quality.
  MP3 (audio only)  The sound track, 192 kbps.
  GIF               Animated, 15 frames per second, at most 480 pixels wide.

Audio (.wav .m4a .flac .ogg .aac .wma)
  MP3, WAV

Pictures (.png .jpg .jpeg .webp .bmp .heic)
  PNG, JPG, WEBP    Phone photos come out the right way up.
                    Transparent areas become white in JPG.
  PDF               A one-page PDF, the page the size of the picture.
                    (To put several pictures in ONE PDF, use Send to > Images to PDF.)
  Remove metadata   Removes GPS location, camera, dates and other hidden
                    information. Keeps the picture quality, the rotation and
                    the colours exactly the same.

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
The notification gives the reason. Details are written to:
  %LOCALAPPDATA%\QuickConvert\log.txt


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
