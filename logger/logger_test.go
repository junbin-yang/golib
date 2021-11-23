package logger

import "testing"

func Test_LOG(t *testing.T) {
	// Path：日志存储路径，默认值/var/log；Level日志等级：默认2；AppName：文件名前缀；Rotate是否自动分割日志
	(&Options{AppName: "test", Path: "/var/log", Level: INFO, Rotate: true}).New()
	// AppName为空时Stdout输出，不写入文件。
	//(&Options{}).New()

	// 是否开启日志通道，开启后日志输出为异步操作。写web服务时开启。
	//go OpenLogChannel()

	logger.Info(1, "ha")
	logger.Debug("debugtest")
	logger.Warn("Warn")
	//logger.SetLogLevel(INFO)
	logger.Error("Error")
	logger.Fatal("Fatal")

}
