package procjob

import (
	"os/exec"
	"sync"
	"syscall"
	"unsafe"
)

const (
	killOnJobClose    = 0x2000 // JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	extendedLimitInfo = 9      // JobObjectExtendedLimitInformation
	processSetQuota   = 0x0100 // rights needed to put a process in a job
	processTerminate  = 0x0001
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	procCreateJob   = kernel32.NewProc("CreateJobObjectW")
	procSetJobInfo  = kernel32.NewProc("SetInformationJobObject")
	procAssignToJob = kernel32.NewProc("AssignProcessToJobObject")
	createJob       sync.Once
	job             uintptr // 0 if it could not be created
)

// Track puts a started program in QuickConvert's job, so Windows ends it
// if QuickConvert ends first. Call it right after cmd.Start(). If anything
// fails, the program simply runs as before.
func Track(cmd *exec.Cmd) {
	createJob.Do(func() {
		h, _, _ := procCreateJob.Call(0, 0)
		if h == 0 {
			return
		}
		var info jobExtendedLimit
		info.BasicLimitInformation.LimitFlags = killOnJobClose
		ok, _, _ := procSetJobInfo.Call(h, extendedLimitInfo, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
		if ok == 0 {
			syscall.CloseHandle(syscall.Handle(h))
			return
		}
		job = h // kept open until QuickConvert ends; Windows then closes it
	})
	if job == 0 || cmd.Process == nil {
		return
	}
	p, err := syscall.OpenProcess(processSetQuota|processTerminate, false, uint32(cmd.Process.Pid))
	if err != nil {
		return
	}
	procAssignToJob.Call(job, uintptr(p))
	syscall.CloseHandle(p)
}
