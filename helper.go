package main

import (
	"strconv"
)

// bool值转string
func BoolToStr(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// 接口值转string
func InterfaceToStr(v interface{}) string {
	switch v.(type) {
	case int:
		return strconv.Itoa(v.(int))
	case string:
		return v.(string)
	case float64:
		return strconv.FormatFloat(v.(float64),'f',-1,32)
	default:
		return ""
	}
}
