// +build linux
package utilx

import (
	"runtime"
	"syscall"
	"time"
)

func SysinfoGet() Sysinfo {

	var si syscall.Sysinfo_t
	syscall.Sysinfo(&si)

	return Sysinfo{
		CpuNum:    runtime.NumCPU(),
		Uptime:    si.Uptime,
		Loads:     si.Loads,
		MemTotal:  si.Totalram / 1024,
		MemFree:   si.Freeram / 1024,
		MemShared: si.Sharedram / 1024,
		MemBuffer: si.Bufferram / 1024,
		SwapTotal: si.Totalswap / 1024,
		SwapFree:  si.Freeswap / 1024,
		Procs:     si.Procs,
		TimeNow:   time.Now().Format(time.RFC3339),
	}
}
