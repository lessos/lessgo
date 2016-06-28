package iossdb

import (
	"bytes"

	"testing"
)

func TestSendBuf(t *testing.T) {
	var tests = []struct {
		in  []interface{}
		out string
	}{
		// test signed integer type
		{[]interface{}{"set", "age", 19}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int8(19)}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int16(19)}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int32(19)}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int64(19)}, "3\nset\n3\nage\n2\n19\n\n"},

		// test unsigned integer type
		{[]interface{}{"set", "age", 19}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int8(19)}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int16(19)}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int32(19)}, "3\nset\n3\nage\n2\n19\n\n"},
		{[]interface{}{"set", "age", int64(19)}, "3\nset\n3\nage\n2\n19\n\n"},
	}

	for k, test := range tests {
		buf, err := send_buf(test.in)
		if err != nil {
			t.Errorf("[%d] [% #v] => [% #v] error: %s\n", k, test.in, buf, err.Error())
			break
		}
		if !bytes.Equal([]byte(test.out), buf) {
			t.Errorf("[%d] [% #v] => [% #v] , expect: % #v\n", k, test.in, string(buf), test.out)
			break
		}
	}
}
