// Package media converts video, audio and images with ffmpeg.
//
// Each Handler (Video, Audio, Image) works inside the job's temporary
// folder: the input is there as "in.<ext>", ffmpeg writes "out.<ext>",
// and the result is then saved next to the original.
package media

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"os"
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
	// FieldOrder of the video: "progressive", or "tt", "bb"... for
	// interlaced video (camcorders, DVDs, TV recordings).
	FieldOrder string
	// Streams lists every stream in the file, in order.
	Streams []Stream
}

// Stream is one track inside a media file.
type Stream struct {
	Index         int    // position in the file (ffmpeg's "0:<Index>")
	Type          string // "video", "audio", "subtitle" or "data"
	Codec         string // e.g. "hevc", "aac"; "none" if ffmpeg doesn't know it
	Cover         bool   // an album-cover picture, not a real video
	Width, Height int    // pictures only
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
			Index         int    `json:"index"`
			Width         int    `json:"width"`
			Height        int    `json:"height"`
			CodecType     string `json:"codec_type"`
			CodecName     string `json:"codec_name"`
			ColorTransfer string `json:"color_transfer"`
			AvgFrameRate  string `json:"avg_frame_rate"` // e.g. "30/1" or "120/1"
			FieldOrder    string `json:"field_order"`
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
		info.Streams = append(info.Streams, Stream{Index: s.Index, Type: s.CodecType, Codec: s.CodecName,
			Cover: s.Disposition.AttachedPic == 1, Width: s.Width, Height: s.Height})
		switch {
		case s.CodecType == "video" && s.Disposition.AttachedPic == 0 && !info.HasVideo:
			info.HasVideo, info.VideoCodec, info.ColorTransfer = true, s.CodecName, s.ColorTransfer
			info.FrameRate = parseRate(s.AvgFrameRate)
			info.FieldOrder = s.FieldOrder
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

// MP4Args builds the ffmpeg command for any video → MP4, keeping
// everything, for editing.
//
// Every picture, sound and subtitle track that MP4 can hold is copied as it
// is ("remux"): this takes seconds, even for big files, and loses nothing —
// HDR colours, rotation, every sound track, and the location, date and
// phone model are kept. Only a track MP4 can't hold is converted, at
// near-original quality. Tracks that can't be converted at all (iPhone
// timed-location data, spatial audio unknown to ffmpeg) are left out.
func MP4Args(in, out string, i Info) []string {
	args := start(in)
	var codecs []string
	n := 0 // number of the next stream in the result
	add := func(s Stream, opts ...string) {
		args = append(args, "-map", "0:"+strconv.Itoa(s.Index))
		num := strconv.Itoa(n)
		for k := 0; k < len(opts); k += 2 { // e.g. "-c", "copy" becomes "-c:0", "copy"
			codecs = append(codecs, opts[k]+":"+num, opts[k+1])
		}
		n++
	}
	video := false
	for _, s := range i.Streams {
		switch {
		case s.Type == "video" && !s.Cover && !video: // the picture
			video = true
			switch s.Codec {
			case "hevc":
				add(s, "-c", "copy", "-tag", "hvc1") // label HEVC so Windows, Apple and editors accept it
			case "h264":
				add(s, "-c", "copy")
			default:
				add(s,
					"-c", "libx264", // H.264 video, accepted by every editor
					"-preset", "medium", // balance of speed and file size
					"-crf", "18", // quality (lower = better; 18 looks the same as the original)
					"-pix_fmt", "yuv420p") // colour format all players support
			}
		case s.Type == "audio" && !s.Cover:
			switch {
			case s.Codec == "none" || s.Codec == "" || strings.Contains(s.Codec, "apac"):
				// ffmpeg can't read it (iPhone spatial audio, called "apple_apac"
				// by newer ffmpeg): leave it out
			case s.Codec == "aac" || s.Codec == "mp3" || s.Codec == "alac" || s.Codec == "ac3" || s.Codec == "eac3":
				add(s, "-c", "copy")
			case s.Codec == "flac" || strings.HasPrefix(s.Codec, "pcm_"):
				add(s, "-c", "alac") // lossless sound stays lossless
			default:
				add(s, "-c", "aac", "-b", "256k") // AAC at 256 kbps
			}
		case s.Type == "subtitle":
			switch s.Codec {
			case "mov_text", "subrip", "ass", "ssa", "webvtt", "text":
				add(s, "-c", "mov_text") // MP4's own subtitle format
			} // picture subtitles (DVD, Blu-ray) can't go in MP4
		}
	}
	args = append(args, codecs...)
	return append(args,
		"-map_metadata", "0", // keep location, date, phone model...
		"-metadata", "major_brand=", "-metadata", "minor_version=", "-metadata", "compatible_brands=", // ...but not the old file-type labels
		"-strict", "unofficial", // also keep Dolby Vision (iPhone HDR) information
		"-movflags", "+faststart+use_metadata_tags", // start playing before fully loaded; keep Apple's metadata
		out)
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
	filters := tidyPicture(i)
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
func GIFArgs(in, out string, i Info) []string {
	filter := strings.Join(append(tidyPicture(i),
		"fps=15",
		"scale='min(480,iw)':-1:flags=lanczos", // shrink to 480 px wide (never enlarge)
		"split[a][b];[a]palettegen[p];[b][p]paletteuse", // build the best 256 colours, then use them
	), ",")
	return append(start(in), "-vf", filter, "-loop", "0", out) // loop forever
}

// tidyPicture returns the first filters for a video that is re-drawn
// (WhatsApp, GIF):
//   - interlaced video (camcorders, DVDs) is made progressive, otherwise
//     moving things show comb-like lines;
//   - "wide pixels" (some camcorders and DVDs store 16:9 pictures in a
//     narrower size) become square, so the picture isn't squashed in apps
//     that ignore the pixel shape (GIFs, some phones).
func tidyPicture(i Info) []string {
	var f []string
	switch i.FieldOrder {
	case "tt", "bb", "tb", "bt":
		f = append(f, "bwdif") // remove interlacing
	}
	return append(f, "scale=iw*sar:ih,setsar=1") // square pixels (no change if already square)
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
	if ext != ".heic" && ext != ".heif" && ext != ".avif" { // AVIF is built the same way
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
	if isOgg(j.Ext) {
		var err error
		if args, err = oggCleanArgs(ctx, j, out); err != nil {
			return app.Result{}, err
		}
	}
	if _, err := j.Eng.Run(ctx, j.Work.Dir, "ffmpeg", args...); err != nil {
		return app.Result{}, err
	}
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, j.Ext) // same type as the original
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}

// CleanArgs builds the ffmpeg command for Clean (all types except OGG).
func CleanArgs(in, out string) []string {
	args := append(start(in),
		"-map", "0:v?", // the picture (or album cover)...
		"-map", "0:a?", // ...the sound...
		"-map", "0:s?", // ...and subtitles, but no data tracks
		"-map_metadata", "-1", // drop all metadata
		"-map_chapters", "-1", // and chapter titles
		"-c", "copy", // copy picture and sound as they are
	)
	switch strings.ToLower(filepath.Ext(out)) {
	case ".m4v":
		// ffmpeg's ".m4v" writer is the old iPod one, which refuses HEVC;
		// an .m4v file is really an MP4.
		args = append(args, "-f", "mp4")
	case ".aiff", ".aif":
		args = append(args, "-write_id3v2", "1") // AIFF keeps its cover only this way
	}
	return append(args, out)
}

func isOgg(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == ".ogg" || ext == ".opus"
}

// oggCleanArgs cleans an OGG/Opus file. ffmpeg can't copy a cover picture
// into OGG as a picture, only as a tag, so the cover is taken out first and
// written back as that tag (in a small metadata file, since a picture is
// too long for a command line).
func oggCleanArgs(ctx context.Context, j app.Job, out string) ([]string, error) {
	args := start(j.Work.In)
	info, err := Probe(ctx, j.Eng, j.Work.Dir, j.Work.In)
	if err != nil {
		return nil, err
	}
	cover := false
	for _, s := range info.Streams {
		if !s.Cover {
			continue
		}
		if _, err := j.Eng.Run(ctx, j.Work.Dir, "ffmpeg", append(start(j.Work.In),
			"-map", "0:"+strconv.Itoa(s.Index), "-c", "copy", "-f", "image2", "cover.bin")...); err != nil {
			return nil, err
		}
		data, err := os.ReadFile(filepath.Join(j.Work.Dir, "cover.bin"))
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(j.Work.Dir, "cover.txt"), CoverMetadata(data, s.Codec, s.Width, s.Height), 0o644); err != nil {
			return nil, err
		}
		args = append(args, "-f", "ffmetadata", "-i", "cover.txt") // input 1: the cover tag
		cover = true
		break
	}
	args = append(args, "-map", "0:a", "-map_metadata", "-1") // the sound, without its metadata...
	if cover {
		args = append(args, "-map_metadata:s:a:0", "1") // ...but with the cover tag
	}
	return append(args, "-map_chapters", "-1", "-c", "copy", out), nil
}

// CoverMetadata makes an ffmpeg metadata file holding a cover picture as
// OGG stores it: a METADATA_BLOCK_PICTURE tag (a FLAC picture block,
// base64-encoded).
func CoverMetadata(data []byte, codec string, width, height int) []byte {
	mime := "image/" + codec
	if codec == "mjpeg" {
		mime = "image/jpeg"
	}
	var b bytes.Buffer
	u32 := func(v int) { binary.Write(&b, binary.BigEndian, uint32(v)) }
	u32(3) // picture type: front cover
	u32(len(mime))
	b.WriteString(mime)
	u32(0) // no description
	u32(width)
	u32(height)
	u32(24) // colour depth
	u32(0)  // not a palette picture
	u32(len(data))
	b.Write(data)
	// In this file format "=", ";", "#" and "\" must be escaped.
	value := strings.ReplaceAll(base64.StdEncoding.EncodeToString(b.Bytes()), "=", `\=`)
	return []byte(";FFMETADATA1\nMETADATA_BLOCK_PICTURE=" + value + "\n")
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
		args = GIFArgs(j.Work.In, out, info)
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

// GIFToMP4 turns an animated GIF into an MP4 video, which plays in every
// app (many apps don't animate GIFs).
func GIFToMP4(ctx context.Context, j app.Job) (app.Result, error) {
	out := "out.mp4"
	return runAndSave(ctx, j, out, GIFToMP4Args(j.Work.In, out))
}

// GIFToMP4Args builds the ffmpeg command for GIF → MP4.
func GIFToMP4Args(in, out string) []string {
	return append(start(in),
		"-vf", "format=rgba,split[fg][b];[b]drawbox=c=white:t=fill[bg];[bg][fg]overlay=format=rgb,"+ // transparent parts white
			"scale=trunc(iw/2)*2:trunc(ih/2)*2,"+ // H.264 needs an even width and height
			"format=yuv420p", // colour format all players support
		"-c:v", "libx264", "-crf", "20", // H.264, good quality
		"-fps_mode", "vfr", // keep each frame's own duration
		"-an",                     // a GIF has no sound
		"-movflags", "+faststart", // starts playing before fully loaded
		out)
}

// iconSizes are the pictures inside a Windows icon, in pixels.
var iconSizes = []string{"256", "48", "32", "16"}

// Icon makes a Windows icon (.ico) from a picture.
func Icon(ctx context.Context, j app.Job) (app.Result, error) {
	in, orientation, err := Readable(ctx, j.Eng, j.Work.Dir, j.Work.In)
	if err != nil {
		return app.Result{}, err
	}
	out := "out.ico"
	return runAndSave(ctx, j, out, IconArgs(in, out, orientation))
}

// IconArgs builds the ffmpeg command for a picture → icon: upright, made
// square with a transparent border (never stretched), in 4 sizes.
func IconArgs(in, out string, orientation int) []string {
	graph := "[0:v]"
	if f, ok := orientationFilters[orientation]; ok {
		graph += f + ","
	}
	graph += "format=rgba,split=" + strconv.Itoa(len(iconSizes)) // one copy per size
	for i := range iconSizes {
		graph += "[s" + strconv.Itoa(i) + "]"
	}
	var maps []string
	for i, size := range iconSizes {
		n := strconv.Itoa(i)
		// Fit inside the square, keeping the proportions (at least 1 px, so
		// very wide or tall pictures work too)...
		fit := "scale='if(gte(iw,ih)," + size + ",max(1,round(" + size + "*iw/ih)))'" +
			":'if(gte(iw,ih),max(1,round(" + size + "*ih/iw))," + size + ")'"
		// ...and centre it on a transparent square.
		graph += ";[s" + n + "]" + fit + ",pad=" + size + ":" + size + ":-1:-1:color=0x00000000[o" + n + "]"
		maps = append(maps, "-map", "[o"+n+"]")
	}
	args := []string{
		"-hide_banner", "-nostdin", "-y",
		"-noautorotate", // we rotate ourselves (above)
		"-i", in,
		"-filter_complex", graph,
	}
	args = append(args, maps...)
	return append(args,
		"-frames:v", "1", // one picture per size
		"-c:v", "png", // PNG inside the icon keeps transparency
		"-map_metadata", "-1",
		out)
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
