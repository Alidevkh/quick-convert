package media

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"quickconvert/internal/actions"
	"quickconvert/internal/app"
	"quickconvert/internal/engines"
	"quickconvert/internal/outpath"
	"quickconvert/internal/testutil"
)

// convert runs handler h on input for action id, the same way app.Run does.
func convert(t *testing.T, h app.Handler, input, id string) (app.Result, error) {
	t.Helper()
	a, ok := actions.Lookup(id)
	if !ok {
		t.Fatalf("unknown action %s", id)
	}
	w, err := outpath.NewWork(filepath.Dir(input))
	if err != nil {
		t.Fatal(err)
	}
	defer w.Cleanup()
	if err := w.LinkInput(input); err != nil {
		t.Fatal(err)
	}
	name := filepath.Base(input)
	j := app.Job{
		Input: input, Base: strings.TrimSuffix(name, filepath.Ext(name)),
		Ext: strings.ToLower(filepath.Ext(name)), Action: a, Work: w, Eng: engines.Default(),
	}
	return h(context.Background(), j)
}

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
