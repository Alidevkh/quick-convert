package platform

// Limits how many jobs run at once (Windows).

import (
	"syscall"
	"unsafe"

	"quickconvert/internal/actions"
)

// slots is how many jobs of each kind may run at the same time, across all
// QuickConvert processes: 2 heavy ones (video encoding, PDF shrinking),
// 4 light ones (photos, metadata, text...).
var slots = map[actions.Kind]struct {
	name  string
	count uintptr
}{
	actions.Heavy: {`Local\QuickConvertHeavy`, 2},
	actions.Light: {`Local\QuickConvertLight`, 4},
}

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procCreateSemaphore  = kernel32.NewProc("CreateSemaphoreW")
	procReleaseSemaphore = kernel32.NewProc("ReleaseSemaphore")
)

// Acquire waits until a slot for this kind of job is free and returns a
// function that frees it again.
//
// Windows starts one QuickConvert process per selected file. They all open
// the same named "semaphore" (a shared counter, e.g. starting at 2): each
// job takes one from the counter, and waits while it is at 0.
func Acquire(kind actions.Kind) (release func()) {
	slot := slots[kind]
	name, _ := syscall.UTF16PtrFromString(slot.name)
	handle, _, _ := procCreateSemaphore.Call(0, slot.count, slot.count, uintptr(unsafe.Pointer(name)))
	if handle == 0 {
		return func() {} // could not create it: run without limit rather than fail
	}
	h := syscall.Handle(handle)
	syscall.WaitForSingleObject(h, syscall.INFINITE)
	return func() {
		procReleaseSemaphore.Call(uintptr(h), 1, 0)
		syscall.CloseHandle(h)
	}
}
