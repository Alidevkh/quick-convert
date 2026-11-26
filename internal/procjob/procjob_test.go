package procjob

import (
	"runtime"
	"testing"
	"unsafe"
)

// The structure handed to Windows must have exactly Windows' layout.
func TestLimitLayout(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("layout checked for 64-bit only")
	}
	if n := unsafe.Sizeof(jobBasicLimit{}); n != 64 {
		t.Errorf("basic limit is %d bytes, want 64", n)
	}
	if n := unsafe.Sizeof(jobExtendedLimit{}); n != 144 {
		t.Errorf("extended limit is %d bytes, want 144", n)
	}
	if off := unsafe.Offsetof(jobBasicLimit{}.LimitFlags); off != 16 {
		t.Errorf("LimitFlags at %d, want 16", off)
	}
}
