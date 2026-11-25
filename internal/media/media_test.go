package media

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"quickconvert/internal/app"
	"quickconvert/internal/engines"
	"quickconvert/internal/testutil"
)

var convert = testutil.RunHandler

func contains(args []string, want string) bool {
	return strings.Contains(strings.Join(args, " "), want)
}

// iPhone video: every picture and sound track MP4 can hold is copied as it
// is; the timed-location data track (which MP4 can't hold) is skipped.
func TestMP4ArgsCopiesIPhoneVideo(t *testing.T) {
	i := Info{HasVideo: true, HasAudio: true, VideoCodec: "hevc", AudioCodec: "aac", Streams: []Stream{
		{Index: 0, Type: "video", Codec: "hevc"},
		{Index: 1, Type: "audio", Codec: "aac"},
		{Index: 2, Type: "audio", Codec: "aac"},
		{Index: 3, Type: "data", Codec: "none"},
	}}
	args := strings.Join(MP4Args("in.mov", "out.mp4", i), " ")
	for _, want := range []string{"-map 0:0 -map 0:1 -map 0:2 ", "-c:0 copy", "-c:1 copy", "-c:2 copy",
		"-tag:0 hvc1", "-map_metadata 0", "use_metadata_tags", "-strict unofficial"} {
		if !strings.Contains(args, want) {
			t.Errorf("missing %q in %s", want, args)
		}
	}
	if strings.Contains(args, "0:3") || strings.Contains(args, "libx264") {
		t.Errorf("args = %s", args)
	}
}

// Sound MP4 can't hold is converted (AAC, or ALAC when the sound is
// lossless); a track ffmpeg doesn't know (e.g. spatial audio) is skipped.
func TestMP4ArgsAudioTracks(t *testing.T) {
	i := Info{HasVideo: true, HasAudio: true, VideoCodec: "h264", AudioCodec: "pcm_s16le", Streams: []Stream{
		{Index: 0, Type: "video", Codec: "h264"},
		{Index: 1, Type: "audio", Codec: "pcm_s16le"},
		{Index: 2, Type: "audio", Codec: "none"},
		{Index: 3, Type: "audio", Codec: "opus"},
		{Index: 6, Type: "audio", Codec: "apple_apac"}, // iPhone spatial audio: newer ffmpeg names it but can't decode it
		{Index: 4, Type: "subtitle", Codec: "subrip"},
		{Index: 5, Type: "subtitle", Codec: "hdmv_pgs_subtitle"},
	}}
	args := strings.Join(MP4Args("in.mkv", "out.mp4", i), " ")
	for _, want := range []string{"-map 0:0 -map 0:1 -map 0:3 -map 0:4 ", "-c:0 copy", "-c:1 alac", "-c:2 aac -b:2 256k", "-c:3 mov_text"} {
		if !strings.Contains(args, want) {
			t.Errorf("missing %q in %s", want, args)
		}
	}
	if strings.Contains(args, "0:2 ") || strings.Contains(args, "0:5") || strings.Contains(args, "0:6") {
		t.Errorf("args = %s", args)
	}
}

// Video MP4 can't hold is re-compressed at near-original quality.
func TestMP4ArgsReencodeVP9(t *testing.T) {
	i := Info{HasVideo: true, HasAudio: true, VideoCodec: "vp9", AudioCodec: "opus", Streams: []Stream{
		{Index: 0, Type: "video", Codec: "vp9"},
		{Index: 1, Type: "audio", Codec: "opus"},
	}}
	args := strings.Join(MP4Args("in.webm", "out.mp4", i), " ")
	if !strings.Contains(args, "-c:0 libx264") || !strings.Contains(args, "-crf:0 18") || !strings.Contains(args, "-c:1 aac") {
		t.Errorf("args = %s", args)
	}
}

// A real iPhone-style video: the picture is bit-for-bit the same, and the
// HDR colours, rotation, second sound track, location, date and phone model
// are all kept (for editing).
func TestVideoForEditingKeepsEverything(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	flat := filepath.Join(t.TempDir(), "flat.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=30:s=640x360",
		"-f", "lavfi", "-i", "sine=d=1", "-f", "lavfi", "-i", "sine=d=1:frequency=880",
		"-map", "0", "-map", "1", "-map", "2", "-c:v", "libx265", "-x265-params", "log-level=error",
		"-pix_fmt", "yuv420p10le", "-color_primaries", "bt2020", "-color_trc", "arib-std-b67", "-colorspace", "bt2020nc",
		"-c:a", "aac", "-metadata", "com.apple.quicktime.location.ISO6709=+25.2048+055.2708/",
		"-metadata", "com.apple.quicktime.model=iPhone 15 Pro", "-metadata", "creation_time=2026-01-01T06:00:00Z",
		"-movflags", "use_metadata_tags", flat)
	in := filepath.Join(dir, "IMG_0001.MOV")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-display_rotation", "90", "-i", flat, "-map", "0", "-c", "copy",
		"-map_metadata", "0", "-movflags", "use_metadata_tags", in)

	r, err := convert(t, Video, in, "mp4")
	if err != nil {
		t.Fatal(err)
	}
	out := r.Outputs[0]
	if out != filepath.Join(dir, "IMG_0001.mp4") {
		t.Fatalf("output %v", r.Outputs)
	}
	hash := func(f string) string {
		return testutil.Run(t, "ffmpeg", "-v", "error", "-i", f, "-map", "0:v:0", "-c", "copy", "-f", "hash", "-hash", "md5", "-")
	}
	if hash(in) != hash(out) {
		t.Error("the picture was re-compressed")
	}
	streams := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries",
		"stream=codec_type,codec_tag_string,color_transfer:stream_side_data=rotation", "-of", "compact", out)
	for _, want := range []string{"hvc1", "arib-std-b67", "rotation=90"} {
		if !strings.Contains(streams, want) {
			t.Errorf("missing %q: %s", want, streams)
		}
	}
	if n := strings.Count(streams, "codec_type=audio"); n != 2 {
		t.Errorf("%d sound tracks, want 2", n)
	}
	tags := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags", "-of", "compact", out)
	for _, want := range []string{"+25.2048+055.2708", "iPhone 15 Pro", "2026-01-01T06:00:00"} {
		if !strings.Contains(tags, want) {
			t.Errorf("missing %q: %s", want, tags)
		}
	}
	if strings.Contains(tags, "major_brand=qt") {
		t.Errorf("the QuickTime brand was copied as a tag: %s", tags)
	}
}

func TestVideoMovToMp4Remux(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideo(t, dir, "clip.mov", "libx264", true)
	r, err := convert(t, Video, in, "mp4")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "clip.mp4") {
		t.Errorf("output %v", r.Outputs)
	}
	c := testutil.Codecs(t, r.Outputs[0])
	if c["video"] != "h264" || c["audio"] != "aac" {
		t.Errorf("codecs %v", c)
	}
}

func TestVideoHEVCRemuxTagged(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideo(t, dir, "iphone.mov", "libx265", true)
	r, err := convert(t, Video, in, "mp4")
	if err != nil {
		t.Fatal(err)
	}
	tag := testutil.Run(t, "ffprobe", "-v", "error", "-select_streams", "v:0",
		"-show_entries", "stream=codec_tag_string", "-of", "csv=p=0", r.Outputs[0])
	if strings.TrimSpace(tag) != "hvc1" {
		t.Errorf("tag = %q", tag)
	}
}

func TestVideoWithSubtitleTrack(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideoWithSubs(t, dir, "phone.mov")
	r, err := convert(t, Video, in, "mp4")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["subtitle"] != "mov_text" {
		t.Errorf("subtitles should be kept: %v", c)
	}
}

func TestVideoToMp3(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideo(t, dir, "clip.mov", "libx264", true)
	r, err := convert(t, Video, in, "mp3")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["audio"] != "mp3" {
		t.Errorf("codecs %v", c)
	}
}

func TestVideoNoAudioToMp3(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideo(t, dir, "silent.mov", "libx264", false)
	_, err := convert(t, Video, in, "mp3")
	var ue *app.UserError
	if !errors.As(err, &ue) || ue.Msg != "This video has no audio" {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "silent.mp3")); !os.IsNotExist(err) {
		t.Error("silent.mp3 should not exist")
	}
}

func TestVideoToGifWidthCap(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideoSize(t, dir, "wide.mp4", "1280x720")
	r, err := convert(t, Video, in, "gif")
	if err != nil {
		t.Fatal(err)
	}
	if w, _ := testutil.Dims(t, r.Outputs[0]); w != 480 {
		t.Errorf("width %d", w)
	}
}

func TestAudioWavToMp3(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeAudio(t, dir, "song.wav")
	r, err := convert(t, Audio, in, "mp3")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "song.mp3") {
		t.Errorf("output %v", r.Outputs)
	}
}

func TestImageRotated(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "exiftool")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "phone.jpg", 40, 20, false)
	testutil.SetOrientation(t, in, 6)
	r, err := convert(t, Image, in, "png")
	if err != nil {
		t.Fatal(err)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 20 || h != 40 {
		t.Errorf("got %dx%d, want 20x40", w, h)
	}
}

func TestImageTransparentToJpgWhite(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "logo.png", 16, 16, true)
	r, err := convert(t, Image, in, "jpg")
	if err != nil {
		t.Fatal(err)
	}
	pixel, err := exec.Command("ffmpeg", "-v", "error", "-i", r.Outputs[0],
		"-frames:v", "1", "-f", "rawvideo", "-pix_fmt", "rgb24", "-").Output()
	if err != nil || len(pixel) < 3 {
		t.Fatalf("decode: %v", err)
	}
	for i, c := range pixel[:3] {
		if c < 245 {
			t.Errorf("channel %d = %d, want white", i, c)
		}
	}
}

func TestImageToWebp(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "pic.png", 16, 16, false)
	r, err := convert(t, Image, in, "webp")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["video"] != "webp" {
		t.Errorf("codecs %v", c)
	}
}

func TestUnicodeFolderEndToEnd(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := filepath.Join(t.TempDir(), "مجلد 'أ'")
	os.Mkdir(dir, 0o755)
	in := testutil.MakeVideo(t, dir, "رحلة صيف.mov", "libx264", true)
	r, err := convert(t, Video, in, "mp3")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "رحلة صيف.mp3") {
		t.Errorf("output %v", r.Outputs)
	}
}

// The rotation is applied to the pixels, so the result must not also carry
// an orientation tag (newer ffmpeg could otherwise make viewers rotate twice).
func TestImageRotatedHasNoOrientationTag(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "exiftool")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "phone.jpg", 40, 20, false)
	testutil.SetOrientation(t, in, 6)
	r, err := convert(t, Image, in, "jpg")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(testutil.Run(t, "exiftool", "-s3", "-Orientation", r.Outputs[0])); got != "" {
		t.Errorf("output still has Orientation %q", got)
	}
}

// iPhone photos are stored as tiles; ffmpeg refuses extra steps (rotate,
// white background...) on them in the same command. And the photo must
// come out portrait, as on the phone — not rotated twice.
func TestImageIPhoneHEICToPNG(t *testing.T) {
	testutil.NeedHEIC(t)
	testutil.Need(t, "exiftool")
	dir := t.TempDir()
	in := testutil.IPhoneHEIC(t, dir, "IMG_6736.HEIC")
	r, err := convert(t, Image, in, "png")
	if err != nil {
		t.Fatal(err)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 3024 || h != 4032 {
		t.Errorf("got %dx%d, want portrait 3024x4032", w, h)
	}
}

func TestImageIPhoneHEICToJPG(t *testing.T) {
	testutil.NeedHEIC(t)
	testutil.Need(t, "exiftool")
	dir := t.TempDir()
	in := testutil.IPhoneHEIC(t, dir, "IMG_6736.HEIC")
	r, err := convert(t, Image, in, "jpg")
	if err != nil {
		t.Fatal(err)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 3024 || h != 4032 {
		t.Errorf("got %dx%d, want portrait 3024x4032", w, h)
	}
}

func TestAudioFormats(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	wav := testutil.MakeAudio(t, dir, "song.wav")
	for id, codec := range map[string]string{"m4a": "aac", "flac": "flac", "ogg": "vorbis", "mp3": "mp3"} {
		r, err := convert(t, Audio, wav, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if c := testutil.Codecs(t, r.Outputs[0]); c["audio"] != codec {
			t.Errorf("%s: codec %v, want %s", id, c, codec)
		}
	}
	mp3 := testutil.MakeAudio(t, dir, "voice.mp3")
	r, err := convert(t, Audio, mp3, "wav")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["audio"] != "pcm_s16le" {
		t.Errorf("mp3 -> wav codec %v", c)
	}
}

// iPhone videos already have AAC sound: M4A copies it without re-encoding.
func TestVideoToM4ACopiesAAC(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideo(t, dir, "clip.mov", "libx264", true)
	r, err := convert(t, Video, in, "m4a")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "clip.m4a") {
		t.Fatalf("output %v", r.Outputs)
	}
	c := testutil.Codecs(t, r.Outputs[0])
	if c["audio"] != "aac" || c["video"] != "" {
		t.Errorf("codecs %v, want aac only", c)
	}
	if !contains(AudioArgs("in.mov", "out.m4a", "m4a", Info{HasAudio: true, AudioCodec: "aac"}), "-c:a copy") {
		t.Error("AAC sound should be copied, not re-encoded")
	}
}

// lumaOfFirstFrame returns the average brightness (0-255) of a video's first frame.
func lumaOfFirstFrame(t *testing.T, path string) float64 {
	t.Helper()
	raw, err := exec.Command("ffmpeg", "-v", "error", "-i", path, "-frames:v", "1",
		"-f", "rawvideo", "-pix_fmt", "gray", "-").Output()
	if err != nil || len(raw) == 0 {
		t.Fatalf("decode: %v", err)
	}
	sum := 0
	for _, b := range raw {
		sum += int(b)
	}
	return float64(sum) / float64(len(raw))
}

func streamField(t *testing.T, path, field string) string {
	t.Helper()
	return strings.TrimSpace(testutil.Run(t, "ffprobe", "-v", "error", "-select_streams", "v:0",
		"-show_entries", "stream="+field, "-of", "csv=p=0", path))
}

func TestWhatsAppFromHEVC4K(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "clip.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=5:s=3840x2160",
		"-f", "lavfi", "-i", "sine=d=1", "-c:v", "libx265", "-preset", "ultrafast", "-pix_fmt", "yuv420p",
		"-c:a", "aac", "-shortest", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "clip (WhatsApp).mp4") {
		t.Fatalf("output %v", r.Outputs)
	}
	c := testutil.Codecs(t, r.Outputs[0])
	if c["video"] != "h264" || c["audio"] != "aac" {
		t.Errorf("codecs %v", c)
	}
	if pf := streamField(t, r.Outputs[0], "pix_fmt"); pf != "yuv420p" {
		t.Errorf("pix_fmt %s", pf)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 1920 || h != 1080 {
		t.Errorf("size %dx%d, want 1920x1080", w, h)
	}
}

// A portrait phone video (stored landscape + a rotation tag) must stay portrait.
func TestWhatsAppPortrait(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	flat := testutil.MakeVideoSize(t, t.TempDir(), "flat.mp4", "1920x1080")
	in := filepath.Join(dir, "portrait.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-display_rotation", "90", "-i", flat, "-c", "copy", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 1080 || h != 1920 {
		t.Errorf("size %dx%d, want portrait 1080x1920", w, h)
	}
}

// HDR (PQ) at about 100 nits is normal white-ish; a naive conversion
// shows it as mid grey (~128). Tone mapping must keep it bright.
func TestWhatsAppHDRToneMapped(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "hdr.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "color=c=gray:s=320x240:d=1:r=5",
		"-c:v", "libx265", "-preset", "ultrafast", "-pix_fmt", "yuv420p10le",
		"-color_primaries", "bt2020", "-color_trc", "smpte2084", "-colorspace", "bt2020nc", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if l := lumaOfFirstFrame(t, r.Outputs[0]); l < 180 {
		t.Errorf("brightness %.0f, want > 180 (tone-mapped, not grey)", l)
	}
	if trc := streamField(t, r.Outputs[0], "color_transfer"); trc != "bt709" {
		t.Errorf("color_transfer %q, want bt709", trc)
	}
}

func TestWhatsAppNoAudio(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := testutil.MakeVideo(t, dir, "silent.mov", "libx264", false)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["video"] != "h264" || c["audio"] != "" {
		t.Errorf("codecs %v", c)
	}
}

// iPhone HDR video uses HLG. Its "reference white" (75% signal, BT.2408)
// must come out white on normal screens, not grey (~190).
func TestWhatsAppHLGReferenceWhite(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "iphone-hdr.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "color=c=0xBFBFBF:s=320x240:d=1:r=5",
		"-c:v", "libx265", "-preset", "ultrafast", "-pix_fmt", "yuv420p10le",
		"-color_primaries", "bt2020", "-color_trc", "arib-std-b67", "-colorspace", "bt2020nc", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if l := lumaOfFirstFrame(t, r.Outputs[0]); l < 220 {
		t.Errorf("reference white came out %.0f, want > 220", l)
	}
}

// Slow-motion / 120 fps iPhone videos: more than 60 fps needs H.264 levels
// many devices can't play. The WhatsApp MP4 is capped at 60 fps.
func TestWhatsAppCapsFrameRate(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "slomo.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=120:s=640x360",
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if fps := streamField(t, r.Outputs[0], "avg_frame_rate"); fps != "60/1" {
		t.Errorf("frame rate %s, want 60/1", fps)
	}
}

// A normal 30 fps video keeps its frame rate.
func TestWhatsAppKeeps30fps(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "normal.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=30:s=640x360",
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuv420p", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if fps := streamField(t, r.Outputs[0], "avg_frame_rate"); fps != "30/1" {
		t.Errorf("frame rate %s, want 30/1", fps)
	}
}

// Full-range sources (screen recordings, MJPEG) must come out in the normal
// "TV" range, which every player shows correctly.
func TestWhatsAppFullRangeBecomesTV(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "screen.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=10:s=640x360",
		"-c:v", "libx264", "-preset", "ultrafast", "-pix_fmt", "yuvj420p", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if cr := streamField(t, r.Outputs[0], "color_range"); cr != "tv" {
		t.Errorf("color_range %q, want tv", cr)
	}
}

// When the file can't even be read, the log must keep ffprobe's own words.
func TestProbeFailureKeepsEngineError(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "broken.mov")
	os.WriteFile(in, []byte("not a video"), 0o644)
	_, err := convert(t, Video, in, "mp4")
	var ee *engines.EngineError
	if !errors.As(err, &ee) {
		t.Errorf("engine error lost: %v", err)
	}
}

// A phone video with GPS location, date and camera: all removed, picture
// and sound untouched, still portrait.
func TestCleanVideoMetadata(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	flat := filepath.Join(t.TempDir(), "flat.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=10:s=320x240",
		"-f", "lavfi", "-i", "sine=d=1", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-shortest",
		"-metadata", "location=+25.2048+055.2708/", "-metadata", "creation_time=2026-01-01T10:00:00Z",
		"-metadata", "make=Apple", "-metadata", "model=iPhone 15", flat)
	in := filepath.Join(dir, "IMG_0001.MOV")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-display_rotation", "90", "-i", flat, "-c", "copy", "-map_metadata", "0", in)

	r, err := convert(t, Clean, in, "clean")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "IMG_0001 (clean).MOV") && r.Outputs[0] != filepath.Join(dir, "IMG_0001 (clean).mov") {
		t.Fatalf("output %v", r.Outputs)
	}
	tags := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags:stream_tags", "-of", "compact", r.Outputs[0])
	for _, gone := range []string{"location", "creation_time", "make", "model", "iPhone"} {
		if strings.Contains(tags, gone) {
			t.Errorf("%q still present: %s", gone, tags)
		}
	}
	if rot := strings.TrimSpace(testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "stream_side_data=rotation", "-of", "csv=p=0", r.Outputs[0])); !strings.HasPrefix(rot, "90") && !strings.HasPrefix(rot, "-90") {
		t.Errorf("rotation lost: %q", rot)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["video"] != "h264" || c["audio"] != "aac" {
		t.Errorf("codecs %v", c)
	}
	if orig := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags=location", "-of", "csv=p=0", in); !strings.Contains(orig, "25.2048") {
		t.Error("original was modified")
	}
}

func TestCleanAudioMetadata(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "song.mp3")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "sine=d=1", "-c:a", "libmp3lame",
		"-metadata", "artist=Someone", "-metadata", "title=Private", in)
	r, err := convert(t, Clean, in, "clean")
	if err != nil {
		t.Fatal(err)
	}
	tags := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags", "-of", "compact", r.Outputs[0])
	if strings.Contains(tags, "Someone") || strings.Contains(tags, "Private") {
		t.Errorf("tags still present: %s", tags)
	}
}

// The WhatsApp MP4 is for sharing: no location or other metadata.
func TestWhatsAppDropsMetadata(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "trip.mov")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=10:s=320x240",
		"-c:v", "libx264", "-pix_fmt", "yuv420p", "-metadata", "location=+25.2048+055.2708/", in)
	r, err := convert(t, Video, in, "whatsapp")
	if err != nil {
		t.Fatal(err)
	}
	if tags := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags", "-of", "compact", r.Outputs[0]); strings.Contains(tags, "location") {
		t.Errorf("location still present: %s", tags)
	}
}

// OGG/Opus covers are kept (ffmpeg can only write them back as a tag).
func TestCleanOggKeepsCover(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	for _, f := range []struct{ name, codec string }{{"song.ogg", "libvorbis"}, {"voice.opus", "libopus"}} {
		t.Run(f.name, func(t *testing.T) {
			dir := t.TempDir()
			src := t.TempDir()
			cover := filepath.Join(src, "cover.png")
			testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=s=64x48:d=1", "-frames:v", "1", cover)
			data, _ := os.ReadFile(cover)
			meta := filepath.Join(src, "meta.txt")
			os.WriteFile(meta, CoverMetadata(data, "png", 64, 48), 0o644)
			in := filepath.Join(dir, f.name)
			testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "sine=d=1:sample_rate=48000",
				"-f", "ffmetadata", "-i", meta, "-map", "0:a", "-map_metadata", "1", "-metadata", "artist=Someone",
				"-c:a", f.codec, in)
			if c := testutil.Codecs(t, in); c["video"] != "png" {
				t.Fatalf("test file has no cover: %v", c)
			}

			r, err := convert(t, Clean, in, "clean")
			if err != nil {
				t.Fatal(err)
			}
			if tags := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags:stream_tags", "-of", "compact", r.Outputs[0]); strings.Contains(tags, "Someone") {
				t.Errorf("tags still present: %s", tags)
			}
			got := testutil.Run(t, "ffmpeg", "-v", "error", "-i", r.Outputs[0], "-map", "0:v:0", "-c", "copy", "-f", "image2pipe", "-")
			if got != string(data) {
				t.Error("cover lost or changed")
			}
		})
	}
}

func TestCleanOggWithoutCover(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	in := filepath.Join(t.TempDir(), "song.ogg")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "sine=d=1", "-c:a", "libvorbis",
		"-metadata", "artist=Someone", in)
	r, err := convert(t, Clean, in, "clean")
	if err != nil {
		t.Fatal(err)
	}
	if tags := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "format_tags:stream_tags", "-of", "compact", r.Outputs[0]); strings.Contains(tags, "Someone") {
		t.Errorf("tags still present: %s", tags)
	}
}

// More video types: camcorder, old phone, DVD-style and Flash video.
func TestVideoMoreTypesToMP4(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	for _, c := range []struct{ name, vcodec, acodec, wantVideo string }{
		{"cam.mts", "libx264", "ac3", "h264"},    // camcorder: copied
		{"cam.m2ts", "libx264", "aac", "h264"},   // copied
		{"old.3gp", "h263", "aac", "h264"},       // re-compressed
		{"dvd.mpg", "mpeg2video", "mp2", "h264"}, // re-compressed
		{"web.flv", "libx264", "aac", "h264"},    // copied
		{"apple.m4v", "libx264", "aac", "h264"},  // copied
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			in := filepath.Join(dir, c.name)
			args := []string{"-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:s=176x144", "-f", "lavfi", "-i", "sine=d=1",
				"-c:v", c.vcodec, "-c:a", c.acodec, "-shortest"}
			if strings.HasSuffix(c.name, "ts") {
				args = append(args, "-f", "mpegts")
			}
			testutil.Run(t, "ffmpeg", append(args, in)...)
			r, err := convert(t, Video, in, "mp4")
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.Codecs(t, r.Outputs[0]); got["video"] != c.wantVideo || got["audio"] == "" {
				t.Errorf("codecs %v", got)
			}
		})
	}
}

func TestAudioMoreTypes(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	aiff := testutil.MakeAudio(t, dir, "mac.aiff")
	r, err := convert(t, Audio, aiff, "mp3")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["audio"] != "mp3" {
		t.Errorf("aiff codecs %v", c)
	}
	// AMR voice memos: ffmpeg can always read them; making one needs an
	// encoder only some builds have.
	amr := filepath.Join(dir, "memo.amr")
	if err := exec.Command("ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "sine=d=1:sample_rate=8000",
		"-ac", "1", "-c:a", "libopencore_amrnb", "-b:a", "12.2k", amr).Run(); err != nil {
		t.Skip("this ffmpeg can't make AMR test files")
	}
	r, err = convert(t, Audio, amr, "m4a")
	if err != nil {
		t.Fatal(err)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["audio"] != "aac" {
		t.Errorf("amr codecs %v", c)
	}
}

// A GIF becomes an MP4 that plays everywhere: H.264, even size,
// transparent parts white, every frame kept.
func TestGIFToMP4(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe")
	dir := t.TempDir()
	in := filepath.Join(dir, "funny.gif")
	testutil.Run(t, "ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "testsrc=d=1:r=10:s=101x81",
		"-vf", "format=rgba,geq=r='r(X,Y)':g='g(X,Y)':b='b(X,Y)':a='if(gt(X,50),0,255)',"+
			"split[a][b];[a]palettegen=reserve_transparent=1[p];[b][p]paletteuse=alpha_threshold=128", in) // right half transparent
	r, err := convert(t, GIFVideo, in, "mp4")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "funny.mp4") {
		t.Errorf("output %v", r.Outputs)
	}
	if c := testutil.Codecs(t, r.Outputs[0]); c["video"] != "h264" {
		t.Errorf("codecs %v", c)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 100 || h != 80 {
		t.Errorf("size %dx%d, want 100x80", w, h)
	}
	frames := strings.TrimSpace(testutil.Run(t, "ffprobe", "-v", "error", "-count_frames", "-select_streams", "v:0",
		"-show_entries", "stream=nb_read_frames", "-of", "csv=p=0", r.Outputs[0]))
	if frames != "10" {
		t.Errorf("%s frames, want 10", frames)
	}
	px := testutil.Run(t, "ffmpeg", "-v", "error", "-i", r.Outputs[0], "-frames:v", "1",
		"-vf", "crop=4:4:90:40,scale=1:1", "-f", "rawvideo", "-pix_fmt", "rgb24", "-")
	if len(px) != 3 || px[0] < 240 || px[1] < 240 || px[2] < 240 {
		t.Errorf("transparent part should be white, got %v", []byte(px))
	}
}

// New picture types convert like the others.
func TestImageMoreTypes(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "exiftool")
	dir := t.TempDir()
	for _, name := range []string{"scan.tif", "anim.gif"} {
		in := testutil.MakeImage(t, dir, name, 64, 48, false)
		r, err := convert(t, Image, in, "jpg")
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if w, h := testutil.Dims(t, r.Outputs[0]); w != 64 || h != 48 {
			t.Errorf("%s: %dx%d", name, w, h)
		}
	}
	avif := filepath.Join(dir, "web.avif")
	if err := exec.Command("ffmpeg", "-v", "error", "-y", "-f", "lavfi", "-i", "color=c=red:s=64x48", "-frames:v", "1",
		"-c:v", "libaom-av1", "-still-picture", "1", avif).Run(); err != nil {
		t.Skip("this ffmpeg can't make AVIF test files")
	}
	r, err := convert(t, Image, avif, "png")
	if err != nil {
		t.Fatal(err)
	}
	if w, h := testutil.Dims(t, r.Outputs[0]); w != 64 || h != 48 {
		t.Errorf("avif: %dx%d", w, h)
	}
}

// A picture becomes a Windows icon with the four usual sizes, square, with
// the empty space transparent.
func TestIcon(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "exiftool")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "logo.png", 300, 150, false)
	r, err := convert(t, Icon, in, "ico")
	if err != nil {
		t.Fatal(err)
	}
	if r.Outputs[0] != filepath.Join(dir, "logo.ico") {
		t.Errorf("output %v", r.Outputs)
	}
	sizes := testutil.Run(t, "ffprobe", "-v", "error", "-show_entries", "stream=width,height,pix_fmt", "-of", "csv=p=0", r.Outputs[0])
	for _, want := range []string{"256,256,rgba", "48,48,rgba", "32,32,rgba", "16,16,rgba"} {
		if !strings.Contains(sizes, want) {
			t.Errorf("missing %s: %s", want, sizes)
		}
	}
	corner := testutil.Run(t, "ffmpeg", "-v", "error", "-i", r.Outputs[0], "-map", "0:0", "-frames:v", "1",
		"-vf", "crop=1:1:0:0", "-f", "rawvideo", "-pix_fmt", "rgba", "-")
	if len(corner) != 4 || corner[3] != 0 {
		t.Errorf("the padding should be transparent: %v", []byte(corner))
	}
}

// A sideways phone photo becomes an upright icon.
func TestIconRotated(t *testing.T) {
	testutil.Need(t, "ffmpeg", "ffprobe", "exiftool")
	dir := t.TempDir()
	in := testutil.MakeImage(t, dir, "p.jpg", 64, 32, false)
	testutil.SetOrientation(t, in, 6)
	if !contains(IconArgs("in.jpg", "out.ico", 6), "transpose=1") {
		t.Error("rotation not applied")
	}
	if _, err := convert(t, Icon, in, "ico"); err != nil {
		t.Fatal(err)
	}
}
