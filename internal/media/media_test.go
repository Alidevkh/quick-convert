package media

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"quickconvert/internal/app"
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
