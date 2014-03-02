package pass

import (
    "../deps/go.crypto/scrypt"
    "crypto/rand"
    "encoding/base64"
    "io"
    "strconv"
)

var DefVersion = "L1"

func Hash(algo, passwd string) string {

    u := make([]byte, 16)
    _, err := io.ReadFull(rand.Reader, u)
    if err != nil {
        return ""
    }
    salt := base64.StdEncoding.EncodeToString(u)

    hash, _ := scrypt.Key([]byte(passwd), []byte(salt), 1<<16, 8, 1, 32)

    return "L1g81" +
        salt +
        base64.StdEncoding.EncodeToString(hash)
}

func Check(passwd, hashed string) bool {

    if len(hashed) < 40 {
        return false
    }

    if hashed[:2] == "L1" {
        N, _ := strconv.ParseUint(hashed[2:3], 36, 32)
        r, _ := strconv.ParseUint(hashed[3:4], 36, 32)
        p, _ := strconv.ParseUint(hashed[4:5], 36, 32)
        salt := hashed[5:29]

        hash, _ := scrypt.Key([]byte(passwd), []byte(salt), 1<<N, int(r), int(p), 32)

        return hashed[29:] == base64.StdEncoding.EncodeToString(hash)
    }

    return false
}
