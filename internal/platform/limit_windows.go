package platform

import (
	"syscall"
	"unsafe"
)

// heavySlots is how many heavy jobs (video encoding, PDF shrinking) may run
// at the same time across all QuickConvert processes.
const heavySlots = 2

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procCreateSemaphore  = kernel32.NewProc("CreateSemaphoreW")
	procReleaseSemaphore = kernel32.NewProc("ReleaseSemaphore")
)

// AcquireHeavy waits until a heavy-job slot is free and returns a function
// that frees it again.
//
// Windows starts one QuickConvert process per selected file. They all open
// the same named "semaphore" (a shared counter starting at 2): each job
// takes one from the counter, and waits while it is at 0.
func AcquireHeavy() (release func()) {
	name, _ := syscall.UTF16PtrFromString(`Local\QuickConvertHeavy`)
	handle, _, _ := procCreateSemaphore.Call(0, heavySlots, heavySlots, uintptr(unsafe.Pointer(name)))
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
