package logger

import "testing"
import "fmt"

func Test_LOG(t *testing.T) {
	// Path：日志存储路径，默认值/var/log；Level日志等级；AppName：文件名前缀；Rotate是否自动分割日志;是否接管Stdout和Stderr;
	//(&Options{AppName: "test", Path: "/var/log", Level: InfoLevel, Rotate: true, TakeStd: true}).New()
	// AppName为空时Stdout输出，不写入文件。
	(&Options{}).New()

	// 是否开启日志通道，开启后日志输出为异步操作。
	//go Asyn()

	SetLogLevel(TraceLevel)
	//Panic("panic test")
	//Fatal("fff")
	Error("Error", "wtf")
	Warn("Warn", "wtf")
	Info(1, "info", "wtf")
	Debug("debugtest", "wtf")
	Trace("trace test")
	fmt.Println("stdout")
}
