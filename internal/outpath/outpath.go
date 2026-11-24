// Package outpath decides where results are saved and makes sure nothing is
// ever overwritten.
//
// How one conversion uses it:
//
//  1. NewWork creates a temporary folder next to the original file,
//     e.g. "C:\Videos\.qc-1a2b3c4d.part".
//  2. LinkInput puts the original inside that folder under a short, plain
//     name ("in.mov"). Engines only ever see such simple names, which avoids
//     problems with Arabic letters, spaces or quotes in real file names.
//  3. The engine writes its result inside the folder (e.g. "out.mp4").
//  4. Commit moves the result out to the first free name next to the
//     original: "trip.mp4", or "trip (1).mp4" if that is taken, and so on.
//  5. Cleanup deletes the temporary folder, whether the job worked or not.
package outpath

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ErrNotWritable is returned when we cannot create files in the folder of
// the original (read-only folder, CD, protected location...).
var ErrNotWritable = errors.New("Can't save in this folder")

// maxTries is how many " (n)" numbers we try before giving up.
const maxTries = 10000

// Candidate builds an output file name.
// n == 0 gives "base.ext"; n > 0 gives "base (n).ext".
func Candidate(base string, n int, ext string) string {
	if n == 0 {
		return base + ext
	}
	return fmt.Sprintf("%s (%d)%s", base, n, ext)
}

// PageName builds the file name for one page of a multi-page result,
// e.g. "report-2.png", or "report (1)-2.png" when n == 1.
func PageName(base string, n, page int, ext string) string {
	return fmt.Sprintf("%s-%d%s", Candidate(base, n, ""), page, ext)
}

// Work is the temporary folder for one job.
type Work struct {
	Dir    string // full path of the temporary folder
	Parent string // folder of the original file; results are saved here
	In     string // name of the input copy inside Dir, e.g. "in.mov"

	// link creates newname as a hard link to oldname and fails if newname
	// already exists. It is os.Link normally; tests replace it.
	link func(oldname, newname string) error

	// move renames oldname to newname and fails if newname already exists.
	// Used when the drive has no hard links (FAT32/exFAT USB sticks, network
	// shares). On Windows it is MoveFile; elsewhere it reports "unsupported".
	move func(oldname, newname string) error
}

// NewWork creates the temporary folder inside parent.
func NewWork(parent string) (*Work, error) {
	random := make([]byte, 4)
	if _, err := rand.Read(random); err != nil {
		return nil, err
	}
	dir := filepath.Join(parent, ".qc-"+hex.EncodeToString(random)+".part")
	if err := os.Mkdir(dir, 0o755); err != nil {
		return nil, ErrNotWritable
	}
	return &Work{Dir: dir, Parent: parent, link: os.Link, move: moveNoReplace}, nil
}

// LinkInput makes the original file available inside the work folder as
// "in" + its extension (lower case). A hard link is instant even for huge
// videos; if the drive does not support hard links we copy instead.
// The original itself is never touched.
//
// A read-only original is always copied: on Windows, deleting a hard link
// to a read-only file would remove the read-only flag from the original.
func (w *Work) LinkInput(src string) error {
	w.In = "in" + strings.ToLower(filepath.Ext(src))
	dst := filepath.Join(w.Dir, w.In)
	if info, err := os.Stat(src); err == nil && info.Mode().Perm()&0o200 != 0 {
		if err := os.Link(src, dst); err == nil {
			return nil
		}
	}
	return copyNew(src, dst)
}

// Commit moves the file tmpName (inside the work folder) to the first free
// name "base.ext", "base (1).ext", ... in the parent folder, and returns
// the full path it was saved to. An existing file is never overwritten.
func (w *Work) Commit(tmpName, base, ext string) (string, error) {
	src := filepath.Join(w.Dir, tmpName)
	for n := 0; n < maxTries; n++ {
		dst := filepath.Join(w.Parent, Candidate(base, n, ext))
		err := w.claim(src, dst)
		if err == nil {
			os.Remove(src)
			return dst, nil
		}
		if !os.IsExist(err) {
			return "", err
		}
		// Name taken: try the next number.
	}
	return "", fmt.Errorf("no free file name for %s%s", base, ext)
}

// CommitPages saves several page files together. It picks one number n so
// that every page name is free ("report-1.png", "report-2.png"... or
// "report (1)-1.png"...), so the pages of one PDF always stay together.
func (w *Work) CommitPages(tmpNames []string, base, ext string) ([]string, error) {
	for n := 0; n < maxTries; n++ {
		if !w.allPagesFree(len(tmpNames), base, n, ext) {
			continue
		}
		var saved []string
		for i, tmp := range tmpNames {
			dst := filepath.Join(w.Parent, PageName(base, n, i+1, ext))
			if err := w.claim(filepath.Join(w.Dir, tmp), dst); err != nil {
				return saved, err
			}
			os.Remove(filepath.Join(w.Dir, tmp))
			saved = append(saved, dst)
		}
		return saved, nil
	}
	return nil, fmt.Errorf("no free file names for %s pages", base)
}

// allPagesFree reports whether none of the page names for number n exist.
func (w *Work) allPagesFree(pages int, base string, n int, ext string) bool {
	for i := 1; i <= pages; i++ {
		if _, err := os.Lstat(filepath.Join(w.Parent, PageName(base, n, i, ext))); err == nil {
			return false
		}
	}
	return true
}

// Cleanup deletes the work folder and anything left inside it.
func (w *Work) Cleanup() {
	os.RemoveAll(w.Dir)
}

// claim creates dst with the content of src, failing with an "exists"
// error if dst is already there. All three ways are safe when two jobs
// want the same name at the same moment:
//  1. a hard link (instant);
//  2. if the drive has no hard links, a move that refuses to replace
//     (instant, and never leaves a half-written file);
//  3. as a last resort, creating dst exclusively and copying the bytes.
func (w *Work) claim(src, dst string) error {
	err := w.link(src, dst)
	if err == nil || os.IsExist(err) {
		return err
	}
	err = w.move(src, dst)
	if err == nil || os.IsExist(err) {
		return err
	}
	return copyNew(src, dst)
}

// copyNew copies src to dst. dst must not exist yet (O_EXCL), so this can
// never overwrite anything. On failure the half-written dst is removed.
func copyNew(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(dst)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(dst)
		return err
	}
	return nil
}
