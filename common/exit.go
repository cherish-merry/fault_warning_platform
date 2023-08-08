package common

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

func AbnormalExit() {
	// 打印程序退出时的堆栈信息
	log.Fatal(string(debug.Stack()))
	// exit
	os.Exit(1)
}
