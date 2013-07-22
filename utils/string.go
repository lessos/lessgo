package utils

import (
    "crypto/rand"
    "fmt"
    "io"
)

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
