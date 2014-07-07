package utils

import (
	"errors"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func NetFreePort(start, end int) (error, int, string) {

	if end > 65535 {
		end = 65535
	}
	if start < 1 {
		start = 1
	}
	if (end - start) < 100 {
		start = 1
		end = 65535
	}

	try := 100
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {

		if try < 0 {
			return errors.New("Error"), 0, "0"
		}
		iport := start + r.Intn(end-start)

		port := strconv.Itoa(iport)
		ln, err := net.Listen("tcp", ":"+port)
		if err == nil {
			ln.Close()
			return nil, iport, port
		}

		try--
	}
}
