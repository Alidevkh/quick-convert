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

func TestMP4ArgsRemuxHEVC(t *testing.T) {
	args := MP4Args("in.mov", "out.mp4", Info{VideoCodec: "hevc", AudioCodec: "aac", HasVideo: true, HasAudio: true})
	if !contains(args, "-c copy") || !contains(args, "-tag:v hvc1") || contains(args, "libx264") {
		t.Errorf("args = %v", args)
	}
}

func TestMP4ArgsReencodeVP9(t *testing.T) {
	args := MP4Args("in.webm", "out.mp4", Info{VideoCodec: "vp9", AudioCodec: "opus", HasVideo: true, HasAudio: true})
	if !contains(args, "libx264") || !contains(args, "192k") || contains(args, "-c copy") {
		t.Errorf("args = %v", args)
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
	if _, has := testutil.Codecs(t, r.Outputs[0])["subtitle"]; has {
		t.Error("subtitle track should be dropped")
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

// iPhones can add extra audio tracks (e.g. spatial audio) in formats MP4
// may not accept, so only the main audio track is kept.
func TestMP4ArgsKeepsOnlyFirstAudioTrack(t *testing.T) {
	args := MP4Args("in.mov", "out.mp4", Info{VideoCodec: "hevc", AudioCodec: "aac", HasVideo: true, HasAudio: true})
	if !contains(args, "-map 0:a:0?") || contains(args, "-map 0:a? ") {
		t.Errorf("args = %v", args)
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

// OGG/Opus covers can't be copied back into OGG by ffmpeg: map sound only.
func TestCleanArgsOggSoundOnly(t *testing.T) {
	args := strings.Join(CleanArgs("in.ogg", "out.ogg"), " ")
	if strings.Contains(args, "0:v") || !strings.Contains(args, "-map 0:a") {
		t.Errorf("args = %s", args)
	}
	if !strings.Contains(strings.Join(CleanArgs("in.mov", "out.mov"), " "), "0:v?") {
		t.Error("videos must keep the picture")
	}
}
