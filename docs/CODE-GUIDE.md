# دليل الكود — QuickConvert

هذا الدليل يشرح أين يوجد كل جزء من البرنامج، وماذا يحدث عندما تضغط بالزر الأيمن على ملف.
الكود مكتوب بلغة **Go**، وكل ملف يبدأ بتعليق بالإنجليزية يشرح وظيفته.

## ماذا يحدث عند الضغط بالزر الأيمن؟

مثال: تضغط على `C:\Videos\trip.mov` ثم تختار **Convert to ▸ MP4**.

1. **ويندوز** يشغّل: `quickconvert.exe mp4 "C:\Videos\trip.mov"`
   (هذا الأمر كتبه ملف `install.cmd` في الريجستري).
2. **`cmd/quickconvert/main.go`** يقرأ الأمر ويستدعي `app.Run`.
3. **`internal/app/app.go`** — قلب البرنامج، يمشي بالترتيب:
   - يتأكد أن `mp4` مسموح لملفات `.mov` (من جدول `actions`).
   - إذا استغرق العمل أكثر من ثانية، يُظهر إشعار «Converting…».
   - إذا كانت المهمة ثقيلة (فيديو أو تصغير PDF) ينتظر دوره، بحيث لا تعمل أكثر من مهمتين معاً.
   - ينشئ مجلداً مؤقتاً بجانب الملف، مثل `C:\Videos\.qc-1a2b3c4d.part`.
   - يستدعي المحوّل المناسب (هنا `media.Video`).
   - يُظهر إشعار «Done» أو سبب الخطأ، ويكتب سطراً في ملف السجل.
   - يحذف المجلد المؤقت دائماً، سواء نجح العمل أو فشل.
4. **`internal/media/media.go`** يشغّل `ffmpeg` داخل المجلد المؤقت بأسماء بسيطة (`in.mov` ← `out.mp4`).
5. **`internal/outpath/outpath.go`** ينقل `out.mp4` إلى `C:\Videos\trip.mp4`. إذا كان الاسم موجوداً يحفظه باسم `trip (1).mp4`، ولا يكتب فوق أي ملف أبداً.

## خريطة الملفات

| الملف | وظيفته |
|---|---|
| `internal/actions/actions.go` | **الجدول الرئيسي**: أي امتداد يظهر له أي خيار في القائمة. لإضافة صيغة جديدة ابدأ من هنا. |
| `internal/app/app.go` | خطوات تنفيذ الطلب من البداية للنهاية (الموضّحة أعلاه). |
| `internal/app/errors.go` | نوعان من الرسائل: `UserError` (خطأ يُعرض للمستخدم) و`Info` (ليس خطأ، مثل «الملف مضغوط أصلاً»). |
| `internal/outpath/outpath.go` | المجلد المؤقت، وتسمية الملفات الناتجة، ومنع الكتابة فوق الملفات. |
| `internal/engines/engines.go` | البحث عن البرامج المساعدة (ffmpeg وغيره) وتشغيلها بدون نافذة سوداء. |
| `internal/media/media.go` | الفيديو والصوت والصور (ffmpeg). |
| `internal/photo/clean.go` | حذف البيانات المخفية من الصور (exiftool). |
| `internal/textcopy/textcopy.go` | تحويل txt ↔ py و m و c وغيرها (نسخ المحتوى كما هو). |
| `internal/pdf/render.go` | PDF ← صور PNG/JPG (mutool). |
| `internal/pdf/shrink.go` | تصغير حجم PDF (Ghostscript). |
| `internal/pdfmd/` | PDF ← Markdown: `stext.go` يقرأ النص ومواقعه، و`markdown.go` يحدد العناوين والقوائم والجداول. |
| `internal/platform/` | أشياء خاصة بويندوز: الإشعارات، وتحديد مهمتين ثقيلتين كحد أقصى. |
| `internal/imagepdf/imagepdf.go` | صور ← PDF (صفحة لكل صورة)، من Send to ومن Convert to. |
| `internal/pdfwriter/pdfwriter.go` | يكتب ملف PDF بنفسه من صور JPEG، بدون إعادة ضغطها. |
| `internal/pdf/merge.go` | دمج عدة ملفات PDF في ملف واحد (mutool). |
| `internal/app/natural.go` | ترتيب الملفات «الطبيعي»: IMG_9 قبل IMG_10. |
| `internal/menu/menu.go` | يكتب `install.cmd` و`uninstall.cmd` من الجدول الرئيسي. |
| `cmd/quickconvert/dispatch.go` | يربط كل خيار بالمحوّل المناسب له. |
| `scripts/build.sh` | يبني ملف `QuickConvert.zip`. |
| `scripts/fetch-engines.sh` | يحمّل البرامج المساعدة الأربعة لويندوز. |

## قائمة Send to (عدة ملفات في ملف واحد)

قائمة **Send to** تختلف عن قائمة الزر الأيمن العادية: ويندوز يعطي البرنامج **كل** الملفات المحددة في تشغيل واحد، مثل:
`quickconvert.exe imgpdf "a.jpg" "b.jpg" "c.jpg"`.
لذلك تمر عبر `app.RunMany` بدلاً من `app.Run`. هي تتجاهل الملفات غير المناسبة، وترتّب الباقي، وتضعها في المجلد المؤقت بأسماء `in-001.jpg` و`in-002.jpg` وهكذا.
الاختصارات نفسها (`Images to PDF.lnk` و`Merge PDFs.lnk`) يصنعها `install.cmd`.

## أمثلة على تعديلات شائعة

- **إضافة امتداد كود جديد (مثل `.cs`):** أضفه إلى قائمة `codeExts` في `actions.go`، ثم أعد البناء وشغّل `install.cmd` من جديد.
- **تغيير جودة MP3:** عدّل `"192k"` في دالة `MP3Args` داخل `media.go`.
- **تغيير دقة صور PDF:** عدّل الثابت `dpi` في `pdf/render.go`.
- **تغيير الحد الأدنى للتوفير عند تصغير PDF:** عدّل الثابت `minSaving` في `pdf/shrink.go`.

## الاختبارات والبناء

```bash
go test ./...          # تشغيل كل الاختبارات
scripts/build.sh       # بناء dist/QuickConvert.zip
```

كل جزء له ملف اختبار بجانبه (`*_test.go`)، واسم كل اختبار يشرح ما يتحقق منه، مثل `TestImageRotated` (هل صورة الجوال المقلوبة تخرج بالاتجاه الصحيح؟).
