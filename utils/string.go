package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"
	"time"
)

const (
	encodeBase36 = "abcdefghijklmnopqrstuvwxyz0123456789"
	encodeBase64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

func init() {
	mrand.Seed(time.Now().UTC().UnixNano())
}

func StringNewRand36(length int) string {

	if length < 1 {
		length = 1
	} else if length > 1000 {
		length = 1000
	}

	buf := make([]byte, length)
	buf[0] = encodeBase36[mrand.Intn(25)]

	for i := 1; i < length; i++ {
		buf[i] = encodeBase36[mrand.Intn(35)]
	}

	return string(buf)
}

func StringNewRand64(length int) string {

	if length < 1 {
		length = 1
	} else if length > 1000 {
		length = 1000
	}

	buf := make([]byte, length)

	for i := 0; i < length; i++ {
		buf[i] = encodeBase64[mrand.Intn(63)]
	}

	return string(buf)
}

func StringNewRand(len int) string {

	u := make([]byte, len/2)

	// Reader is a global, shared instance of a cryptographically strong pseudo-random generator.
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := io.ReadFull(rand.Reader, u)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", u)
}

// NewUUID generates a new UUID based on version 4.
func StringNewUUID() string {

	u := make([]byte, 16)

	// Reader is a global, shared instance of a cryptographically strong pseudo-random generator.
	// On Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the CryptGenRandom API.
	_, err := io.ReadFull(rand.Reader, u)
	if err != nil {
		panic(err)
	}

	// Set version (4) and variant (2).
	var version byte = 4 << 4
	var variant byte = 2 << 4
	u[6] = version | (u[6] & 15)
	u[8] = variant | (u[8] & 15)

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func StringEncode16(str string, slen uint) string {

	if slen < 1 {
		slen = 1
	} else if slen > 32 {
		slen = 32
	}

	h := md5.New()
	io.WriteString(h, str)

	return fmt.Sprintf("%x", h.Sum(nil))[:slen]
}
