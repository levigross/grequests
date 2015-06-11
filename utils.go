package grequests

import (
	"log"
	"strconv"
)

// ConvertIToString will convert an interface to a string
func ConvertIToString(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case bool:
		return strconv.FormatBool(i.(bool))
	case int8:
		return strconv.FormatInt(int64(i.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(i.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(i.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(i.(int64)), 10)
	case int:
		return strconv.FormatInt(int64(i.(int)), 10)
	case uint8:
		return strconv.FormatUint(uint64(i.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(i.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(i.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(i.(uint64)), 10)
	case uint:
		return strconv.FormatUint(uint64(i.(uint)), 10)
	default:
		log.Printf("Cannot convert %T to String type", i)
		return ""
	}
}
