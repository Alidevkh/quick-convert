package outpath

// The "move without replacing" used on drives without hard links (Windows).

import (
	"os"
	"syscall"
)

// moveNoReplace renames a file without ever replacing an existing one.
// Windows' MoveFile does exactly that, also on FAT32/exFAT and network shares.
func moveNoReplace(oldname, newname string) error {
	from, err := syscall.UTF16PtrFromString(oldname)
	if err != nil {
		return err
	}
	to, err := syscall.UTF16PtrFromString(newname)
	if err != nil {
		return err
	}
	if err := syscall.MoveFile(from, to); err != nil {
		return &os.LinkError{Op: "move", Old: oldname, New: newname, Err: err}
	}
	return nil
}
