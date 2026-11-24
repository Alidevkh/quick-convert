// Package photo removes hidden information (metadata) from photos: GPS
// location, camera and lens, dates, software, thumbnails and comments.
//
// Two things are deliberately kept:
//   - the orientation ("which way is up"), or phone photos would show sideways;
//   - the colour profile, or iPhone photos would look washed out.
//
// The picture itself is not re-compressed, so quality stays exactly the same.
package photo

import (
	"context"

	"quickconvert/internal/app"
)

// Clean saves a copy of the photo without metadata as "<name> (clean).<ext>".
func Clean(ctx context.Context, j app.Job) (app.Result, error) {
	out := "out" + j.Ext
	_, err := j.Eng.Run(ctx, j.Work.Dir, "exiftool",
		"-q", "-q", // quiet: no messages or warnings
		"-charset", "filename=utf8", // file names are UTF-8 (matters on Windows)
		"-all=",              // remove all metadata...
		"-tagsFromFile", "@", // ...then copy back from the original:
		"-Orientation", //    the "which way is up" tag
		"-ICC_Profile", //    the colour profile
		"-o", out,      // write a new file; never change the input
		j.Work.In)
	if err != nil {
		return app.Result{}, err
	}
	saved, err := j.Work.Commit(out, j.Base+j.Action.Tag, j.Ext)
	if err != nil {
		return app.Result{}, err
	}
	return app.Result{Outputs: []string{saved}}, nil
}
