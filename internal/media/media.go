// Package media converts video, audio and images with ffmpeg.
//
// Each Handler (Video, Audio, Image) works inside the job's temporary
// folder: the input is there as "in.<ext>", ffmpeg writes "out.<ext>",
// and the result is then saved next to the original.
package media

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	"quickconvert/internal/app"
	"quickconvert/internal/engines"
)

// Info describes the streams inside a media file.
type Info struct {
	VideoCodec string // e.g. "h264", "hevc", "vp9"; "" if none
	AudioCodec string // e.g. "aac", "mp3", "opus"; "" if none
	HasVideo   bool
	HasAudio   bool
	// ColorTransfer of the video, e.g. "bt709" (normal), "smpte2084"
	// (HDR10 / PQ) or "arib-std-b67" (HLG, used by iPhone HDR video).
	ColorTransfer string
	// FrameRate of the video in frames per second (0 if unknown).
	FrameRate float64
}

// HDR reports whether the video uses high dynamic range colours, which
// look washed out or dark on normal screens unless they are converted.
func (i Info) HDR() bool {
	return i.ColorTransfer == "smpte2084" || i.ColorTransfer == "arib-std-b67"
}

// Probe asks ffprobe which video and audio streams a file has.
// Album-cover pictures inside audio files are not counted as video.
func Probe(ctx context.Context, eng *engines.Set, dir, in string) (Info, error) {
	out, err := eng.Run(ctx, dir, "ffprobe",
		"-v", "error", // only print real errors
		"-show_streams", // list every stream
		"-of", "json",   // as JSON, easy to read here
		in)
	if err != nil {
		return Info{}, err
	}
	var data struct {
		Streams []struct {
			CodecType     string `json:"codec_type"`
			CodecName     string `json:"codec_name"`
			ColorTransfer string `json:"color_transfer"`
			AvgFrameRate  string `json:"avg_frame_rate"` // e.g. "30/1" or "120/1"
			Disposition   struct {
				AttachedPic int `json:"attached_pic"`
			} `json:"disposition"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &data); err != nil {
		return Info{}, err
	}
	var info Info
	for _, s := range data.Streams {
		switch {
		case s.CodecType == "video" && s.Disposition.AttachedPic == 0 && !info.HasVideo:
			info.HasVideo, info.VideoCodec, info.ColorTransfer = true, s.CodecName, s.ColorTransfer
			info.FrameRate = parseRate(s.AvgFrameRate)
		case s.CodecType == "audio" && !info.HasAudio:
			info.HasAudio, info.AudioCodec = true, s.CodecName
		}
	}
	return info, nil
}

// parseRate turns ffprobe's "30000/1001" into 29.97 (0 if unknown).
func parseRate(r string) float64 {
	num, den, ok := strings.Cut(r, "/")
	n, err1 := strconv.ParseFloat(num, 64)
	d, err2 := strconv.ParseFloat(den, 64)
	if !ok || err1 != nil || err2 != nil || d == 0 {
		return 0
	}
	return n / d
}

// start is the beginning of every ffmpeg command.
func start(in string) []string {
	return []string{
		"-hide_banner", // don't print ffmpeg's version banner
		"-nostdin",     // never wait for keyboard input
		"-y",           // overwrite out.* (it only exists in our temp folder)
		"-i", in,
	}
}

// MP4Args builds the ffmpeg command for any video → MP4.
//
// If the video and audio are already in formats MP4 supports, they are
// copied as they are ("remux"): this takes seconds, even for big files,
// and loses no quality. Otherwise the video is re-encoded to H.264.
func MP4Args(in, out string, i Info) []string {
	args := append(start(in),
		"-map", "0:v:0", // take the first video stream...
		"-map", "0:a:0?", // ...and the main audio track, if any (drops extra tracks, subtitles, data)
		"-map_metadata", "0", // keep date/title information
	)
	videoOK := i.VideoCodec == "h264" || i.VideoCodec == "hevc"
	audioOK := !i.HasAudio || i.AudioCodec == "aac" || i.AudioCodec == "mp3"
	if videoOK && audioOK {
		args = append(args, "-c", "copy") // copy streams without re-encoding
		if i.VideoCodec == "hevc" {
			args = append(args, "-tag:v", "hvc1") // label HEVC so Windows and phones play it
		}
	} else {
		args = append(args,
			"-c:v", "libx264", // H.264 video, plays everywhere
			"-preset", "medium", // balance of speed and file size
			"-crf", "20", // quality (lower = better; 20 looks like the original)
			"-pix_fmt", "yuv420p", // colour format all players support
			"-c:a", "aac", "-b:a", "192k", // AAC audio at 192 kbps
		)
	}
	return append(args, "-movflags", "+faststart", out) // lets the video start playing before fully loaded
}

// audioCodecs says how each audio format is encoded.
var audioCodecs = map[string][]string{
	"mp3":  {"-c:a", "libmp3lame", "-b:a", "192k"}, // MP3 at 192 kbps
	"m4a":  {"-c:a", "aac", "-b:a", "192k"},        // AAC at 192 kbps (the iPhone/Apple format)
	"wav":  {"-c:a", "pcm_s16le"},                  // uncompressed, CD quality
	"flac": {"-c:a", "flac"},                       // lossless, about half the size of WAV
	"ogg":  {"-c:a", "libvorbis", "-q:a", "6"},     // Vorbis, about 192 kbps
}

// AudioArgs saves the first sound track of a file (audio or video) in
// format "mp3", "m4a", "wav", "flac" or "ogg". For M4A, sound that is
// already AAC (iPhone videos) is copied as it is: instant and lossless.
func AudioArgs(in, out, format string, i Info) []string {
	args := append(start(in),
		"-vn",           // no video (and no album-cover picture)
		"-map", "0:a:0", // the first sound track
	)
	if format == "m4a" && i.AudioCodec == "aac" {
		return append(args, "-c:a", "copy", out) // already AAC: copy without re-encoding
	}
	return append(append(args, audioCodecs[format]...), out)
}

// MP3Args extracts the first audio stream as a 192 kbps MP3.
func MP3Args(in, out string) []string { return AudioArgs(in, out, "mp3", Info{}) }

// WAVArgs extracts the first audio stream as an uncompressed WAV.
func WAVArgs(in, out string) []string { return AudioArgs(in, out, "wav", Info{}) }

// WhatsAppArgs makes an MP4 that plays on any phone, computer or browser
// and is ready for WhatsApp: H.264 video in standard 8-bit colours, AAC
// stereo sound, at most 1080p (the short side), same frame rate.
// HDR video is tone-mapped to normal colours so it looks right everywhere.
func WhatsAppArgs(in, out string, i Info) []string {
	args := append(start(in),
		"-map", "0:v:0", // the video...
		"-map", "0:a:0?", // ...and the main sound track, if there is one
		"-map_metadata", "-1", // drop all metadata (GPS location, dates, camera): it is for sharing
	)
	var filters []string
	if i.FrameRate > 60.5 {
		// Slow-motion / 120 fps videos: more than 60 fps needs H.264 levels
		// that many phones and TVs can't play.
		filters = append(filters, "fps=60")
	}
	filters = append(filters,
		// Short side at most 1080 px (4K becomes 1080p; portrait stays
		// portrait). Sizes are kept even, as H.264 requires. out_range=tv:
		// the normal colour range every player shows correctly (screen
		// recordings are often "full range", which some players show washed out).
		"scale='if(gt(iw,ih),-2,min(1080,trunc(iw/2)*2))':'if(gt(iw,ih),min(1080,trunc(ih/2)*2),-2)':out_range=tv",
	)
	if i.HDR() {
		filters = append(filters,
			"zscale=t=linear:npl=100", // HDR signal to real light levels
			"format=gbrpf32le",        // high-precision colours for the maths
			"zscale=p=bt709",          // HDR colour range to normal (BT.709)
			// Squeeze bright highlights into the normal range. "mobius" keeps
			// normal brightness as it is, so HDR "reference white" stays white
			// (checked: hable made iPhone HLG video look grey/dark).
			"tonemap=tonemap=mobius:desat=0",
			"zscale=t=bt709:m=bt709:r=tv", // back to a normal video signal
		)
	}
	filters = append(filters, "format=yuv420p") // 8-bit 4:2:0: plays on every device
	return append(args,
		"-vf", strings.Join(filters, ","),
		"-c:v", "libx264", "-profile:v", "high", // H.264, the most compatible video format
		"-preset", "medium", "-crf", "23", // good quality, reasonable size
		"-color_primaries", "bt709", "-color_trc", "bt709", "-colorspace", "bt709", "-color_range", "tv", // label as normal colours
		"-c:a", "aac", "-b:a", "128k", "-ac", "2", // AAC stereo sound at 128 kbps
		"-movflags", "+faststart", // starts playing before fully downloaded
		out)
}

// GIFArgs makes an animated GIF: 15 frames per second, at most 480 pixels
// wide, using a colour palette made from the video itself so it looks good.
func GIFArgs(in, out string) []string {
	filter := "fps=15," +
		"scale='min(480,iw)':-1:flags=lanczos," + // shrink to 480 px wide (never enlarge)
		"split[a][b];[a]palettegen[p];[b][p]paletteuse" // build the best 256 colours, then use them
	return append(start(in), "-vf", filter, "-loop", "0", out) // loop forever
}

// orientationFilters turns the EXIF "which way is up" number (1-8) into
// the ffmpeg filters that physically rotate/flip the picture.
var orientationFilters = map[int]string{
	2: "hflip",
	3: "hflip,vflip", // upside down
	4: "vflip",
	5: "transpose=0",
	6: "transpose=1", // rotated 90° clockwise (most common phone photo)
	7: "transpose=3",
	8: "transpose=2",
}

// ImageArgs converts a still image to format "png", "jpg" or "webp".
// orientation is the EXIF orientation of the input (1 = already upright).
func ImageArgs(in, out, format string, orientation int) []string {
	args := []string{
		"-hide_banner", "-nostdin", "-y",
		"-noautorotate", // we rotate ourselves (below), so ffmpeg must not also do it
		"-i", in,
		"-frames:v", "1", // one picture
		"-map_metadata", "-1", // drop old metadata (the rotation is now applied to the pixels)
	}

	var filters []string
	if f, ok := orientationFilters[orientation]; ok {
		filters = append(filters, f)
	}
	switch format {
	case "jpg":
		// JPG has no transparency: put the picture on a white background.
		filters = append(filters, "format=rgba,split[fg][b];[b]drawbox=c=white:t=fill[bg];[bg][fg]overlay,format=yuvj444p")
		args = append(args, "-q:v", "2") // high JPG quality
	case "webp":
		args = append(args, "-c:v", "libwebp", "-quality", "90")
	}
	if len(filters) > 0 {
		args = append(args, "-vf", strings.Join(filters, ","))
	}
	return append(args, out)
}

// maxPageSide is the longest side, in pixels, of a photo put on a PDF page:
// sharp enough to read and print A4, without making the PDF huge.
const maxPageSide = "2480"

// PageArgs prepares a photo to become a PDF page: upright, at most 2480
// pixels on its longest side (never enlarged), transparency on white,
// saved as a JPEG.
func PageArgs(in, out string, orientation int) []string {
	args := []string{
		"-hide_banner", "-nostdin", "-y",
		"-noautorotate", // we rotate ourselves (below)
		"-i", in,
		"-frames:v", "1",
		"-map_metadata", "-1",
	}
	var filters []string
	if f, ok := orientationFilters[orientation]; ok {
		filters = append(filters, f)
	}
	// Shrink the longest side to 2480 px if bigger, keeping the proportions.
	// A photo that is already small keeps its exact size (min(1,...) = 1).
	filters = append(filters,
		"scale='if(gte(iw,ih),min("+maxPageSide+",iw),round(iw*min(1,"+maxPageSide+"/ih)))'"+
			":'if(gte(iw,ih),round(ih*min(1,"+maxPageSide+"/iw)),min("+maxPageSide+",ih))'",
		// White background; "format=rgb" keeps full black and white (no grey).
		"format=rgba,split[fg][b];[b]drawbox=c=white:t=fill[bg];[bg][fg]overlay=format=rgb",
		"format=yuvj420p", // standard JPEG colours, as every PDF reader expects
	)
	return append(args,
		"-vf", strings.Join(filters, ","),
		"-q:v", "3", // good JPEG quality (2 = best, 31 = worst)
		out)
}

// Orientation reads a photo's EXIF orientation with exiftool.
// Returns 1 (upright) if there is none or it cannot be read.
func Orientation(ctx context.Context, eng *engines.Set, dir, in string) int {
	out, err := eng.Run(ctx, dir, "exiftool",
		"-n",  // the number, not a description
		"-s3", // print only the value
		"-Orientation", in)
	if err != nil {
		return 1
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil || n < 1 || n > 8 {
		return 1
	}
	return n
}

// Readable prepares a picture for the other ffmpeg steps and returns the
// file name to use and its EXIF orientation.
//
// iPhone photos (HEIC) are stored as a grid of small tiles. ffmpeg puts the
// tiles together with its own internal step, and then refuses any extra
// step (rotate, shrink, white background...) in the same command
// ("Error opening output files: Invalid argument"). So a HEIC photo is
// first turned into a plain PNG on its own. ffmpeg also applies the HEIC
// photo's own rotation while doing so (as the HEIC standard says), so the
// PNG is already upright: orientation 1.
func Readable(ctx context.Context, eng *engines.Set, dir, in string) (name string, orientation int, err error) {
	ext := strings.ToLower(filepath.Ext(in))
	if ext != ".heic" && ext != ".heif" {
		return in, Orientation(ctx, eng, dir, in), nil
	}
	png := "decoded-" + strings.TrimSuffix(in, filepath.Ext(in)) + ".png"
	_, err = eng.Run(ctx, dir, "ffmpeg",
		"-hide_banner", "-nostdin", "-y",
		"-i", in, // no -noautorotate: let ffmpeg apply the photo's own rotation
		"-frames:v", "1",
		"-update", "1", // one picture, one file
		png)
	return png, 1, err
}

// Clean saves a copy of a video or audio file without metadata (GPS
// location, dates, camera, titles...) as "<name> (clean).<ext>". Picture
// and sound are copied as they are: no quality loss, and the rotation of
// phone videos is kept. Data tracks (iPhone location/timed metadata) are
// dropped; a music file's cover picture is kept.
func Clean(ctx context.Context, j app.Job) (app.Result, error) {
	out := "out" + j.Ext
	args := CleanArgs(j.Work.In, out)
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "ffmpeg", args...); err != nil {
		return app.Result{}, err
	}
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, j.Ext) // same type as the original
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}

// CleanArgs builds the ffmpeg command for Clean.
func CleanArgs(in, out string) []string {
	args := start(in)
	switch strings.ToLower(filepath.Ext(out)) {
	case ".ogg", ".opus":
		// ffmpeg can't write a cover picture back into OGG: sound only.
		args = append(args, "-map", "0:a")
	default:
		args = append(args,
			"-map", "0:v?", // the picture (or album cover)...
			"-map", "0:a?", // ...the sound...
			"-map", "0:s?", // ...and subtitles, but no data tracks
		)
	}
	return append(args,
		"-map_metadata", "-1", // drop all metadata
		"-map_chapters", "-1", // and chapter titles
		"-c", "copy", // copy picture and sound as they are
		out)
}

// Video converts a video to MP4, MP3 or GIF.
func Video(ctx context.Context, j app.Job) (app.Result, error) {
	info, err := Probe(ctx, j.Eng, j.Work.Dir, j.Work.In)
	if err != nil {
		return app.Result{}, &app.UserError{Msg: "This doesn't look like a video file", Err: err}
	}
	out := "out" + j.Action.OutExt
	var args []string
	switch j.Action.ID {
	case "mp4":
		if !info.HasVideo {
			return app.Result{}, &app.UserError{Msg: "This file has no video picture"}
		}
		args = MP4Args(j.Work.In, out, info)
	case "whatsapp":
		if !info.HasVideo {
			return app.Result{}, &app.UserError{Msg: "This file has no video picture"}
		}
		args = WhatsAppArgs(j.Work.In, out, info)
	case "mp3", "m4a": // the sound only
		if !info.HasAudio {
			return app.Result{}, &app.UserError{Msg: "This video has no audio"}
		}
		args = AudioArgs(j.Work.In, out, j.Action.ID, info)
	case "gif":
		if !info.HasVideo {
			return app.Result{}, &app.UserError{Msg: "This file has no video picture"}
		}
		args = GIFArgs(j.Work.In, out)
	}
	return runAndSave(ctx, j, out, args)
}

// Audio converts an audio file to MP3, M4A, WAV, FLAC or OGG.
func Audio(ctx context.Context, j app.Job) (app.Result, error) {
	info, err := Probe(ctx, j.Eng, j.Work.Dir, j.Work.In)
	if err != nil || !info.HasAudio {
		return app.Result{}, &app.UserError{Msg: "This doesn't look like an audio file", Err: err}
	}
	out := "out" + j.Action.OutExt
	return runAndSave(ctx, j, out, AudioArgs(j.Work.In, out, j.Action.ID, info))
}

// Image converts a picture to PNG, JPG or WEBP, keeping it upright.
func Image(ctx context.Context, j app.Job) (app.Result, error) {
	out := "out" + j.Action.OutExt
	in, orientation, err := Readable(ctx, j.Eng, j.Work.Dir, j.Work.In)
	if err != nil {
		return app.Result{}, err
	}
	args := ImageArgs(in, out, j.Action.ID, orientation)
	return runAndSave(ctx, j, out, args)
}

// runAndSave runs ffmpeg and saves out next to the original.
func runAndSave(ctx context.Context, j app.Job, out string, args []string) (app.Result, error) {
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "ffmpeg", args...); err != nil {
		return app.Result{}, err
	}
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, j.Action.OutExt)
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}
