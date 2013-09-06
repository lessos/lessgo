package keeper

import (
    "errors"
)

type ReplyType uint8

const (
    ReplyOK      ReplyType = 0
    ReplyError   ReplyType = 1
    ReplyTimeout ReplyType = 2

    ReplyNil     ReplyType = 10
    ReplyInteger ReplyType = 11
    ReplyString  ReplyType = 12
    ReplyJson    ReplyType = 13
    ReplyMulti   ReplyType = 14
    ReplyWatch   ReplyType = 15
)

type Reply struct {
    Err   error     `json:"err"`
    Type  ReplyType `json:"type"`
    Body  string    `json:"body"`
    Name  string    `json:"name"`
    Elems []*Reply  `json:"elems"`
}

// Str returns the reply value as a string or
// an error, if the reply type is not ReplyString.
func (r *Reply) Str() (string, error) {
    if r.Type == ReplyError {
        return "", r.Err
    }
    if !(r.Type == ReplyString) {
        return "", errors.New("string value is not available for this reply type")
    }

    return r.Body, nil
}

// Bytes is a convenience method for calling Reply.Str() and converting it to []byte.
func (r *Reply) Bytes() ([]byte, error) {
    if r.Type == ReplyError {
        return nil, r.Err
    }
    s, err := r.Str()
    if err != nil {
        return nil, err
    }

    return []byte(s), nil
}

// List returns a multi-bulk reply as a slice of strings or an error.
// The reply type must be ReplyMulti and its elements' types must all be either ReplyString or ReplyNil.
// Nil elements are returned as empty strings.
// Useful for list commands.
func (r *Reply) List() ([]string, error) {
    if r.Type == ReplyError {
        return nil, r.Err
    }
    if r.Type != ReplyMulti {
        return nil, errors.New("reply type is not ReplyMulti")
    }

    strings := make([]string, len(r.Elems))
    for i, v := range r.Elems {
        if v.Type == ReplyString {
            strings[i] = v.Body
        } else if v.Type == ReplyNil {
            strings[i] = ""
        } else {
            return nil, errors.New("element type is not ReplyString or ReplyNil")
        }
    }

    return strings, nil
}

// Map returns a multi-bulk reply as a map[string]string or an error.
// The reply type must be ReplyMulti,
// it must have an even number of elements,
// they must be in a "key value key value..." order and
// values must all be either ReplyString or ReplyNil.
// Nil values are returned as empty strings.
// Useful for hash commands.
func (r *Reply) Hash() (map[string]string, error) {
    if r.Type == ReplyError {
        return nil, r.Err
    }
    rmap := map[string]string{}

    if r.Type != ReplyMulti {
        return nil, errors.New("reply type is not ReplyMulti")
    }

    if len(r.Elems)%2 != 0 {
        return nil, errors.New("reply has odd number of elements")
    }

    for i := 0; i < len(r.Elems)/2; i++ {
        var val string

        key, err := r.Elems[i*2].Str()
        if err != nil {
            return nil, errors.New("key element has no string reply")
        }

        v := r.Elems[i*2+1]
        if v.Type == ReplyString {
            val = v.Body
            rmap[key] = val
        } else if v.Type == ReplyNil {
        } else {
            return nil, errors.New("value element type is not ReplyString or ReplyNil")
        }
    }

    return rmap, nil
}

// String returns a string representation of the reply and its sub-replies.
// This method is mainly used for debugging.
// Use method Reply.Str() for fetching a string reply.
func (r *Reply) String() string {
    switch r.Type {
    case ReplyError:
        return r.Err.Error()
    case ReplyString:
        return r.Body
    case ReplyNil:
        return "<nil>"
    case ReplyMulti:
        s := "[ "
        for _, e := range r.Elems {
            s = s + e.String() + " "
        }
        return s + "]"
    }

    // This should never execute
    return ""
}
