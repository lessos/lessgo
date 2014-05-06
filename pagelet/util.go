package pagelet

import (
    "math"
    "reflect"
)

// Equal is a helper for comparing value equality, following these rules:
//  - Values with equivalent types are compared with reflect.DeepEqual
//  - int, uint, and float values are compared without regard to the type width.
//    for example, Equal(int32(5), int64(5)) == true
//  - strings and byte slices are converted to strings before comparison.
//  - else, return false.
func Equal(a, b interface{}) bool {
    if reflect.TypeOf(a) == reflect.TypeOf(b) {
        return reflect.DeepEqual(a, b)
    }
    switch a.(type) {
    case int, int8, int16, int32, int64:
        switch b.(type) {
        case int, int8, int16, int32, int64:
            return reflect.ValueOf(a).Int() == reflect.ValueOf(b).Int()
        }
    case uint, uint8, uint16, uint32, uint64:
        switch b.(type) {
        case uint, uint8, uint16, uint32, uint64:
            return reflect.ValueOf(a).Uint() == reflect.ValueOf(b).Uint()
        }
    case float32, float64:
        switch b.(type) {
        case float32, float64:
            return reflect.ValueOf(a).Float() == reflect.ValueOf(b).Float()
        }
    case string:
        switch b.(type) {
        case []byte:
            return a.(string) == string(b.([]byte))
        }
    case []byte:
        switch b.(type) {
        case string:
            return b.(string) == string(a.([]byte))
        }
    }
    return false
}

// NotEqual evaluates the comparison a != b
func NotEqual(a, b interface{}) bool {
    return !Equal(a, b)
}

// LessThan evaluates the comparison a < b
func LessThan(a, b interface{}) bool {

    switch a.(type) {
    case int, int8, int16, int32, int64:
        switch b.(type) {
        case int, int8, int16, int32, int64:
            return reflect.ValueOf(a).Int() < reflect.ValueOf(b).Int()
        }
    case uint, uint8, uint16, uint32, uint64:
        switch b.(type) {
        case uint, uint8, uint16, uint32, uint64:
            return reflect.ValueOf(a).Uint() < reflect.ValueOf(b).Uint()
        }
    case float32, float64:
        switch b.(type) {
        case float32, float64:
            return reflect.ValueOf(a).Float() < reflect.ValueOf(b).Float()
        }
    }

    return false
}

// LessEqual evaluates the comparison a <= b
func LessEqual(a, b interface{}) bool {
    return LessThan(a, b) || Equal(a, b)
}

// GreaterThan evaluates the comparison a > b
func GreaterThan(a, b interface{}) bool {
    return !LessEqual(a, b)
}

// GreaterEqual evaluates the comparison a >= b
func GreaterEqual(a, b interface{}) bool {
    return !LessThan(a, b)
}

type Paginator struct {
    ItemCount         int
    CountPerPage      int
    PageCount         int
    CurrentPageNumber int
    FirstPageNumber   int
    PrevPageNumber    int
    NextPageNumber    int
    LastPageNumber    int
    RangeLen          int
    RangeStartNumber  int
    RangeEndNumber    int
    RangePages        []int
}

// Pager build a page collection
//
// @param int current  Number of Current page
// @param int count    Total number of items available
// @param int perpage  Number of items per page
// @param int rangelen Number of discrete page numbers that will be displayed
func Pager(current, count, perpage, rangelen int) Paginator {

    pager := Paginator{
        ItemCount:         count,
        CountPerPage:      perpage,
        CurrentPageNumber: current,
        RangeLen:          rangelen,
    }

    //
    if pager.ItemCount < 0 {
        pager.ItemCount = 0
    }

    if pager.CountPerPage < 1 {
        pager.CountPerPage = 10
    }

    if pager.RangeLen < 1 {
        pager.RangeLen = 1
    }

    pager.PageCount = pager.ItemCount / pager.CountPerPage
    if v := math.Mod(float64(pager.ItemCount), float64(pager.CountPerPage)); v > 0 {
        pager.PageCount++
    }

    // fixed current page number
    if pager.CurrentPageNumber < 1 {
        pager.CurrentPageNumber = 1
    } else if pager.CurrentPageNumber > pager.PageCount {
        pager.CurrentPageNumber = pager.PageCount
    }

    //
    pager.RangeStartNumber = 1
    if (pager.CurrentPageNumber - rangelen/2) > 0 {
        pager.RangeStartNumber = pager.CurrentPageNumber - rangelen/2
    }

    pager.RangeEndNumber = pager.PageCount
    if (pager.RangeStartNumber + rangelen) < pager.PageCount {
        pager.RangeEndNumber = pager.RangeStartNumber + rangelen - 1
    }

    // taking previous page
    if pager.CurrentPageNumber > 1 {
        pager.PrevPageNumber = pager.CurrentPageNumber - 1
    }

    // taking next page
    if pager.CurrentPageNumber < pager.PageCount {
        pager.NextPageNumber = pager.CurrentPageNumber + 1
    }

    // taking pages list
    for i := pager.RangeStartNumber; i <= pager.RangeEndNumber; i++ {
        pager.RangePages = append(pager.RangePages, i)
    }

    if pager.RangeStartNumber > 1 {
        pager.FirstPageNumber = 1
    }

    if pager.RangeEndNumber < pager.PageCount {
        pager.LastPageNumber = pager.PageCount
    }

    return pager
}
