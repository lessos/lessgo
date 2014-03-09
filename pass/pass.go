package pass

import (
    "../deps/go.crypto/scrypt"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "strconv"
)

const (
    // Define the radix 64 encoding/decoding scheme,
    alphabet = "jbzB3WM6uYrPd20plhngE1U45QZLTOcsCy8mVwHkq9RFI/SKGeAXJifNaxt7oD+v"
)

var (
    // Define the default hashing algorithm
    AlgoDefault = "L001"
    // New base64 encoder
    b64Encoding = base64.NewEncoding(alphabet)
)

// HashDefault derives a key from password
func HashDefault(passwd string) (string, error) {

    u := make([]byte, 15) // 120-bit
    if _, err := io.ReadFull(rand.Reader, u); err != nil {
        return "", errors.New("Error: rand.Reader")
    }

    salt := b64Encoding.EncodeToString(u)
    if len(salt) != 20 {
        return "", errors.New("Error: base64.Encode")
    }

    hash, err := scrypt.Key([]byte(passwd), u, 1<<15, 8, 1, 36)
    if err != nil {
        return "", err
    }

    //  0,4   A     The string name of a hashing algorithm
    //  4,1   N     CPU cost parameter, 0-9a-z (0~35)
    //  5,1   r     Memory cost parameter, 0-9a-z (0~35)
    //  6,1   p     Parallelization parameter, 0-9a-z (0~35)
    //  7,20  salt  120-bit salt, convert to base64
    // 27,48  hash  288-bit derived key, convert to base64
    return AlgoDefault +
        "f81" +
        salt +
        b64Encoding.EncodeToString(hash), nil
}

// Check reports whether the given password and hashed key match
func Check(passwd, hash string) bool {

    if len(hash) < 40 {
        return false
    }

    if hash[:4] == AlgoDefault {
        N, _ := strconv.ParseUint(hash[4:5], 36, 32)
        r, _ := strconv.ParseUint(hash[5:6], 36, 32)
        p, _ := strconv.ParseUint(hash[6:7], 36, 32)
        salt, _ := b64Encoding.DecodeString(hash[7:27])

        key, _ := scrypt.Key([]byte(passwd), salt, 1<<N, int(r), int(p), 36)

        return hash[27:] == b64Encoding.EncodeToString(key)
    }

    return false
}
