package utilx

import (
	"net"
	"regexp"
)

func NetLocalAddress() string {

	addrs, _ := net.InterfaceAddrs()
	reg, _ := regexp.Compile(`^(.*)\.(.*)\.(.*)\.(.*)\/(.*)$`)
	for _, addr := range addrs {
		ips := reg.FindStringSubmatch(addr.String())
		if len(ips) != 6 || (ips[1] == "127" && ips[2] == "0") {
			continue
		}
		return ips[1] + "." + ips[2] + "." + ips[3] + "." + ips[4]
	}

	return "127.0.0.1"
}
