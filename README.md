# QuickConvert

Convert files on Windows straight from the right-click menu: no app to open,
nothing running in the background. The result is saved next to the original,
which is never changed or overwritten.

**[⬇ Download QuickConvert.zip](https://github.com/Alidev29/quick-convert/releases/latest/download/QuickConvert.zip)**
(Windows 10 / 11, 64-bit · about 140 MB · everything included)

[العربية ↓](#بالعربية)

---

## Install (1 minute)

1. Download **QuickConvert.zip** (link above) and **unzip** it
   (right-click > *Extract All…*).
2. Open the unzipped folder and double-click **`install.cmd`**.

That's it. No administrator rights, no internet needed afterwards. You can
delete the zip and the unzipped folder afterwards: QuickConvert is installed
in `%LOCALAPPDATA%\QuickConvert`.

**Use it:** right-click a file (on Windows 11 click *Show more options*
first, or hold <kbd>Shift</kbd> while right-clicking) > **Convert to**.

**Check it works on your PC (optional):** double-click **`selftest.cmd`** in
the unzipped folder. It tries every conversion on sample files and writes
`selftest-report.txt`.

**Uninstall:** double-click **`uninstall.cmd`** (in the unzipped folder or in
`%LOCALAPPDATA%\QuickConvert`).

## What it converts

| From | To |
|---|---|
| Video `.mov .mp4 .m4v .mkv .avi .webm .wmv .3gp .mts .m2ts .mpg .mpeg .flv` | **MP4 (full quality, for editing)**: iPhone video repackaged in seconds, nothing re-compressed · **MP4 (WhatsApp / any device)** · MP3 · M4A · GIF |
| Audio `.mp3 .m4a .wav .flac .ogg .opus .aac .wma .aiff .aif .amr` | MP3 · M4A · WAV · FLAC · OGG |
| Pictures `.jpg .png .heic .webp .bmp .gif .tif .avif` | PNG · JPG · WEBP · PDF · ICO (Windows icon) · GIF > MP4 |
| PDF | PNG / JPG (one picture per page, in a folder) · MD (text for AI tools) · **Reduce size** · **Split into pages** · **Extract pages…** |
| SVG / EPUB | SVG > PNG · EPUB > PDF / MD |
| Text | `.txt` ↔ `.py .c .cpp .js .ts .html .css .json .md .m .java .bat .sh …` |

**Remove metadata** (GPS location, dates, camera…) for photos, videos and music.

**Send to** (select several files > right-click > *Send to*):
**Images to PDF** (photos of documents into one PDF) · **Merge PDFs**.

More details: [`packaging/README.txt`](packaging/README.txt) (the guide that
comes inside the zip).

## Good to know

- **Nothing runs in the background.** QuickConvert starts when you click a
  menu entry and closes when the file is done. Selecting many files runs at
  most 2 video/PDF jobs and 4 small jobs at a time.
- **Your original is never touched.** If a result name is taken, it saves
  `photo (1).png` instead of overwriting.
- **Works offline.** No accounts, no uploads: everything happens on your PC.
- **"Windows protected your PC"?** The installer removes the download mark so
  this shouldn't appear. If it does: right-click
  `%LOCALAPPDATA%\QuickConvert\quickconvert.exe` > Properties > tick
  *Unblock*.
- **Something failed?** The notification gives the reason; details are in
  `%LOCALAPPDATA%\QuickConvert\log.txt`.

## For developers

QuickConvert is a small Go program (standard library only) that calls
FFmpeg, MuPDF, Ghostscript and ExifTool. The right-click menu is generated
from one list, [`internal/actions/actions.go`](internal/actions/actions.go).

```bash
go test ./...                 # tests (need ffmpeg, mutool, gs, exiftool on PATH)
scripts/fetch-engines.sh      # download the Windows engines into engines/windows
scripts/build.sh              # build dist/QuickConvert.zip
scripts/release.sh "notes"    # publish a GitHub release with the zip
```

A guided tour of the code (in Arabic): [`docs/CODE-GUIDE.md`](docs/CODE-GUIDE.md).
Design decisions: [`docs/adr/`](docs/adr/).

## License

QuickConvert is released under the [MIT License](LICENSE). The bundled
engines keep their own licenses (GPL / AGPL / Perl): see
[`packaging/THIRD-PARTY-NOTICES.txt`](packaging/THIRD-PARTY-NOTICES.txt).

---

<div dir="rtl">

## بالعربية

أداة لويندوز تحوّل الملفات من قائمة الزر الأيمن مباشرة: بدون فتح أي برنامج، ولا شيء يعمل في الخلفية. الملف الناتج يُحفظ بجانب الأصلي، والأصلي لا يُعدَّل ولا يُستبدل أبداً.

### التثبيت (دقيقة واحدة)

1. حمّل **[QuickConvert.zip](https://github.com/Alidev29/quick-convert/releases/latest/download/QuickConvert.zip)** وفك ضغطه (زر أيمن ← *Extract All*).
2. افتح المجلد وانقر مرتين على **`install.cmd`**.

لا يحتاج صلاحيات المدير ولا إنترنت. بعد التثبيت يمكنك حذف الـ zip والمجلد، فالبرنامج يُنسخ إلى `%LOCALAPPDATA%\QuickConvert`.

**الاستخدام:** زر أيمن على أي ملف (في ويندوز 11 اختر *Show more options* أولاً، أو اضغط <kbd>Shift</kbd> مع الزر الأيمن) ← **Convert to**.

**للتأكد أنه يعمل على جهازك (اختياري):** شغّل **`selftest.cmd`**، فيجرّب كل التحويلات ويكتب تقريراً في `selftest-report.txt`.

**إلغاء التثبيت:** شغّل **`uninstall.cmd`**.

### أبرز المزايا

- **فيديو الآيفون:** MP4 بجودة كاملة للمونتاج (في ثوانٍ وبلا إعادة ضغط)، أو MP4 يعمل على أي جهاز وجاهز للواتساب.
- **الصوت:** MP3 وM4A وWAV وFLAC وOGG.
- **الصور:** PNG وJPG وWEBP وPDF وأيقونات ICO، وصور الآيفون HEIC.
- **PDF:** إلى صور، وإلى نص Markdown للذكاء الاصطناعي، وتصغير الحجم، وتقسيم الصفحات، واستخراج صفحات معيّنة.
- **Send to:** جمع صور المستندات في ملف PDF واحد، ودمج ملفات PDF.
- **Remove metadata:** حذف الموقع والتاريخ ومعلومات الجهاز من الصور والفيديو والصوت.
- **الكتب الإلكترونية EPUB وملفات SVG،** وتحويل ملفات النص إلى ملفات برمجة والعكس.

كل شيء يعمل على جهازك بدون إنترنت، ولا يُرفع أي ملف لأي مكان.

</div>
