# logger

依赖：golang版本大于等于1.16

Demo:

```javascript
package main

import "github.com/junbin-yang/golib/logger"

func main() {
	/*
     *	Path：日志存储路径,默认值/var/log；
     *  Level日志等级：默认INFO；
     *  AppName：文件名前缀；
     *  Rotate是否自动分割日志
     */
	//(&logger.Options{AppName: "dvsobj", Path: "/var/log", Level: 2, LogRotate: true}).New()

	// AppName为空时Stdout输出，不写入文件。
	(&logger.Options{}).New()

	// 开启日志异步写入。建议写web服务时开启。
	//go logger.OpenLogChannel()	

	logger.SetLogLevel(logger.DEBUG)
	logger.Info(1, "ha")
	logger.Debug("debugtest")
	logger.Warn("Warn")
	logger.Error("Error")
	logger.Fatal("Fatal")
}

```

