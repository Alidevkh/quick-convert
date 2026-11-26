// Package procjob makes sure the programs QuickConvert starts (ffmpeg,
// mutool, the "Which pages?" box...) stop when QuickConvert stops.
//
// Normally QuickConvert waits for each program to finish, so nothing is
// left running. But if QuickConvert itself is killed (Task Manager, logging
// off), a program it started would otherwise keep running on its own —
// e.g. ffmpeg using the whole processor for minutes.
//
// On Windows this uses a "job object": a group of processes that Windows
// ends together when the group's owner (QuickConvert) ends.
package procjob

// jobBasicLimit is Windows' JOBOBJECT_BASIC_LIMIT_INFORMATION.
type jobBasicLimit struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

// jobExtendedLimit is Windows' JOBOBJECT_EXTENDED_LIMIT_INFORMATION
// (144 bytes on 64-bit Windows).
type jobExtendedLimit struct {
	BasicLimitInformation jobBasicLimit
	IoInfo                [6]uint64 // IO_COUNTERS, unused
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}
