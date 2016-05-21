## encoding/json

encoding and decoding of JSON objects

## Example
```go
package main

import (
	"fmt"
	"os"

	"github.com/lessos/lessgo/encoding/json"
)

type Object struct {
	Name string `json:"name"`
}

func main() {

	js := `{"name": "demo-value-of-string"}`

	var obj Object
	if err := json.Decode([]byte(js), &obj); err == nil {
		fmt.Println("Decode OK, obj.name =", obj.Name)
	}

	if bs, err := json.Encode(obj, "\t"); err == nil {
		fmt.Println("Encode string: ", string(bs))
	}

	if bsi, err := json.Indent([]byte(js), "\t\t"); err == nil {
		fmt.Println("Indent to : ", string(bsi))
	}

	if err := json.EncodeToFile(obj, "/tmp/output.file.json", "\t"); err == nil {
		fmt.Println("Encode to file OK")
	}

	var obj2 Object
	if err := json.DecodeFile("/tmp/output.file.json", &obj2); err == nil {
		fmt.Println("Decode file OK, obj2.name =", obj2.Name)
	}

	os.Remove("/tmp/output.file.json")
}
```
