// +build darwin
package utilx

import (
	"runtime"
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

	return Sysinfo{
		CpuNum:  runtime.NumCPU(),
		TimeNow: time.Now().Format(time.RFC3339),
	}
}
