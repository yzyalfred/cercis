package utils

import (
	"fmt"
	"github.com/yzyalfred/cercis/log"
	"runtime"
)

func TraceCode(code ...interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	data := ""
	for _, v := range code{
		data += fmt.Sprintf("%v", v)
	}
	data += string(buf[:n])
	log.Error("==> %s\n", data)
}
