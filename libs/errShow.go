package libs

import (
	"runtime/debug"

	. "github.com/soekchl/myUtils"
)

func ErrorShow() {
	err := recover()
	if err == nil {
		return
	}
	Error("Err:", err, " Debug:", string(debug.Stack()))
}
