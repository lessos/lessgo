// +build darwin
package utilx

import (
	"runtime"
	"time"
)

func SysinfoGet() Sysinfo {

	return Sysinfo{
		CpuNum:  runtime.NumCPU(),
		TimeNow: time.Now().Format(time.RFC3339),
	}
}
