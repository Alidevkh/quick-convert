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
			CodecType   string `json:"codec_type"`
			CodecName   string `json:"codec_name"`
			Disposition struct {
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
			info.HasVideo, info.VideoCodec = true, s.CodecName
		case s.CodecType == "audio" && !info.HasAudio:
			info.HasAudio, info.AudioCodec = true, s.CodecName
		}
	}
	return info, nil
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

// Video converts a video to MP4, MP3 or GIF.
func Video(ctx context.Context, j app.Job) (app.Result, error) {
	info, err := Probe(ctx, j.Eng, j.Work.Dir, j.Work.In)
	if err != nil {
		return app.Result{}, &app.UserError{Msg: "This doesn't look like a video file"}
	}
	out := "out" + j.Action.OutExt
	var args []string
	switch j.Action.ID {
	case "mp4":
		if !info.HasVideo {
			return app.Result{}, &app.UserError{Msg: "This file has no video picture"}
		}
		args = MP4Args(j.Work.In, out, info)
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
		return app.Result{}, &app.UserError{Msg: "This doesn't look like an audio file"}
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
