// +build linux
package utilx

import (
	"runtime"
	"syscall"
	"time"
)

type Sysinfo struct {
	CpuNum    int       `json:"cn"`
	Uptime    int64     `json:"ut"`
	Loads     [3]uint64 `json:"ld"`
	MemTotal  uint64    `json:"mt"`
	MemFree   uint64    `json:"mf"`
	MemShared uint64    `json:"ms"`
	MemBuffer uint64    `json:"mb"`
	SwapTotal uint64    `json:"st"`
	SwapFree  uint64    `json:"sf"`
	Procs     uint16    `json:"pc"`
	TimeNow   string    `json:"tn"`
}

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
