package common

import (
	"fmt"
	"os"
	"runtime/debug"
)

func AbnormalExit() {
	// 打印程序退出时的堆栈信息
	fmt.Println(string(debug.Stack()))
	// exit
	os.Exit(1)
}
