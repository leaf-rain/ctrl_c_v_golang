package dingding

import (
	"fmt"
	"strconv"
	"strings"
)

var errorStr = "runtime error"
var panicStr = "runtime/panic.go"

type PanicInfo struct {
	Ty     string
	Method string
	Line   string
}

func GetPanicInfo(content string) *PanicInfo {
	slice := strings.Split(content, "\n")
	var info = PanicInfo{}
	var length = len(slice)
	var tmp []string
	for index, item := range slice {
		if strings.Contains(item, errorStr) {
			info.Ty = item
		}
		if strings.Contains(item, panicStr) {
			if index+1 < length {
				tmp = strings.Split(slice[index+1], "(")
				if len(tmp) == 1 {
					info.Method = slice[index+1]
				} else {
					info.Method = strings.Join(tmp[:len(tmp)-1], "(")
				}
			}
			if index+2 < length {
				tmp = strings.Split(slice[index+2], " ")
				if len(tmp) == 1 {
					info.Line = slice[index+2]
				} else {
					info.Line = strings.Join(tmp[:len(tmp)-1], " ")
				}
			}
			break
		}
	}
	if info.Line == "" || info.Ty == "" {
		return nil
	}
	return &info
}

func FloatDecimal(data float64) int64 {
	var str = fmt.Sprintf("%f", data)
	slice := strings.Split(str, ".")
	str = slice[len(slice)-1]
	return StrToInt64(str)
}

func StrToInt64(str string) int64 {
	var num, _ = strconv.ParseInt(str, 10, 64)
	return int64(num)
}
