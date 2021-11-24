# logger

基于github.com/sirupsen/logrus二次封装，实现以下功能：

1、自动分割文件

2、自动清除过期文件

3、异步写日志

4、接管Stdout和Stderr信息

依赖：golang版本大于等于1.16

Demo:

```javascript
package main

import (
	"github.com/junbin-yang/golib/logger"
    "fmt"
)

func main() {
	/*
     *	Path：日志存储路径,默认值/var/log；
     *  Level日志等级；
     *  AppName：文件名前缀；
     *  Rotate是否自动分割日志
     *  TakeStd是否接管Stdout和Stderr
     */
	//(&logger.Options{AppName: "dvsobj", Path: "/var/log", Level: 2, LogRotate: true, TakeStd: true}).New()

	// AppName为空时Stdout输出，不写入文件。
	(&logger.Options{}).New()

	// 开启日志异步写入。建议开启。
	//go logger.Asyn()	

	logger.SetLogLevel(logger.InfoLevel)
	logger.Info(1, "ha")
	logger.Debug("debugtest")
	logger.Warn("Warn")
	logger.Error("Error")
	logger.Fatal("Fatal")
	fmt.Println("stdout")
}

```

