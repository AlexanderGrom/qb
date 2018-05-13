package qb

import (
	"fmt"
	"strconv"
)

// Convert interface to string
func toString(x interface{}) string {
	switch x := x.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	case int:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(int64(x), 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(uint64(x), 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', 6, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', 6, 64)
	case []byte:
		return string(x)
	case []rune:
		return string(x)
	case nil:
		return ""
	default:
		return fmt.Sprint(x)
	}
}

// IntWeight returns number of digits in an int
func intWeight(x int) int {
	var p = 10
	for i := 1; i < 19; i++ {
		if x < p {
			return i
		}
		p *= 10
	}
	return 19
}
