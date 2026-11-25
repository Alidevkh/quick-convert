//go:build !windows

package outpath

// Stand-in for the move used on Windows drives without hard links.

import "errors"

// moveNoReplace is only needed on Windows (drives without hard links).
// Elsewhere it reports "unsupported" so the copy fallback is used.
func moveNoReplace(oldname, newname string) error {
	return errors.ErrUnsupported
}
