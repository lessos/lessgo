package pagelet

import (
    "io/ioutil"
    "strings"
)

// Reads the lines of the given file.  Panics in the case of error.
func ReadLines(filename string) ([]string, error) {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return strings.Split(string(bytes), "\n"), nil
}
