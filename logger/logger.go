package logger

import (
	"bytes"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
	FATAL = 5
	TRACE = 6
)

var obj *logrus.Logger = logrus.New()
var defaultChannel *DataContainer = InitDataContainer()
var channelSwitch bool

type Options struct {
	AppName  string
	Path     string //日志路径
	Level    int    //输出级别
	Rotate   bool
	KeepDays int64
}

func (this *Options) New() {
	if this.Level == 0 {
		this.Level = INFO
	}
	SetLogLevel(this.Level)

	if this.Path == "" {
		this.Path = "/var/log"
	}

	if this.KeepDays == 0 {
		this.KeepDays = 7
	}

	os.MkdirAll(this.Path, 0777)

	if this.AppName != "" {
		logfile := this.Path + "/" + this.AppName + "-" + carbon.Now().ToDateString() + ".log"
		src, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|0644)
		if err == nil {
			obj.Out = src

			if this.Rotate {
				go func(Path, AppName string) {
					c := cron.New()
					c.AddFunc("0 0 0 * * ?", func() {
						logfile := Path + "/" + AppName + "-" + carbon.Now().ToDateString() + ".log"
						src, _ := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|0644)
						obj.Out = src

						var diff_time int64 = 3600 * 24 * this.KeepDays
						now_time := time.Now().Unix()
						delpath := this.Path
						err := filepath.Walk(delpath, func(delpath string, f os.FileInfo, err error) error {
							if f == nil {
								return err
							}
							file_time := f.ModTime().Unix()
							if (now_time - file_time) > diff_time {
								os.RemoveAll(delpath)
							}
							return nil
						})
						if err != nil {
							fmt.Printf("filepath.Walk() returned %v\r\n", err)
						}
					})
					c.Start()
					select {}
				}(this.Path, this.AppName)
			}
		} else {
			fmt.Println("打开日志文件失败:", err)
			writers := []io.Writer{os.Stdout}
			fileAndStdoutWriter := io.MultiWriter(writers...)
			obj.SetOutput(fileAndStdoutWriter)
		}
		obj.SetFormatter(new(LogFormatter))
	} else {
		writers := []io.Writer{os.Stdout}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		obj.SetOutput(fileAndStdoutWriter)
		obj.SetFormatter(new(LogFormatter))
	}
}

func SetLogLevel(level int) {
	switch level {
	case DEBUG:
		obj.SetLevel(logrus.DebugLevel) // Info()、Warn()、Error()、Debug()和Fatal(),更多详细信息
	case INFO:
		obj.SetLevel(logrus.InfoLevel) // Info()、Warn()、Error()和Fatal()
	case WARN:
		obj.SetLevel(logrus.WarnLevel) // Warn()、Error()和Fatal(),更多详细信息
	case ERROR:
		obj.SetLevel(logrus.ErrorLevel) // Error()和Fatal(),更多详细信息
	case FATAL:
		obj.SetLevel(logrus.FatalLevel) // Fatal(),更多详细信息
	default:
		obj.SetLevel(logrus.TraceLevel) // Info()、Warn()、Error()和Fatal(),更多详细信息
	}
}

func OpenLogChannel() {
	channelSwitch = true
	for {
		itemInterface := defaultChannel.Pop()
		if itemInterface != nil {
			item := itemInterface.(map[string]interface{})
			m := fmt.Sprint(item["Data"])
			if item["Type"] == "Info" {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Info(m[1:(len(m) - 1)])
			}
			if item["Type"] == "Debug" {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Debug(m[1:(len(m) - 1)])
			}
			if item["Type"] == "Error" {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Error(m[1:(len(m) - 1)])
			}
			if item["Type"] == "Warn" {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Warn(m[1:(len(m) - 1)])
			}
			if item["Type"] == "Fatal" {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Fatal(m[1:(len(m) - 1)])
			}
		}
	}
}

func Info(o ...interface{}) {
	msg := map[string]interface{}{
		"Type": "Info",
		"Data": o,
	}

	if obj.GetLevel() == logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := fmt.Sprint(o)
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Info(m[1:(len(m) - 1)])
	}
}

func Debug(o ...interface{}) {
	msg := map[string]interface{}{
		"Type": "Debug",
		"Data": o,
	}

	if obj.GetLevel() == logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := fmt.Sprint(o)
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Info(m[1:(len(m) - 1)])
	}
}

func Error(o ...interface{}) {
	msg := map[string]interface{}{
		"Type": "Error",
		"Data": o,
	}

	if obj.GetLevel() == logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := fmt.Sprint(o)
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Info(m[1:(len(m) - 1)])
	}
}

func Warn(o ...interface{}) {
	msg := map[string]interface{}{
		"Type": "Warn",
		"Data": o,
	}

	if obj.GetLevel() == logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := fmt.Sprint(o)
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Info(m[1:(len(m) - 1)])
	}
}

func Fatal(o ...interface{}) {
	msg := map[string]interface{}{
		"Type": "Fatal",
		"Data": o,
	}

	if obj.GetLevel() == logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}

	if channelSwitch {
		if defaultChannel.Push(msg) {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	} else {
		m := fmt.Sprint(o)
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Info(m[1:(len(m) - 1)])
		os.Exit(0)
	}
}

func printCaller() (string, string, int) {
	var pcs [32]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return frame.Function, frame.File, frame.Line
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//日志自定义格式
type LogFormatter struct{}

//格式详情
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t := carbon.Now().ToDateString() + "," + fmt.Sprint(time.Now().Nanosecond()/1e6)
	msg := fmt.Sprintf("%s [%s] %s\n", t, strings.ToUpper(entry.Level.String()), entry.Message)
	if entry.Data["Func"] != nil {
		msg = fmt.Sprintf("%s [%s:%d][%s][RTID:%d][%s] %s\n", t, entry.Data["File"], entry.Data["Line"], entry.Data["Func"], entry.Data["GID"], strings.ToUpper(entry.Level.String()), entry.Message)
	}
	return []byte(msg), nil
}

// 消息队列封装
type DataContainer struct {
	Queue chan interface{}
}

func InitDataContainer() (dc *DataContainer) {
	dc = &DataContainer{}
	dc.Queue = make(chan interface{})
	return dc
}

//非阻塞push
func (dc *DataContainer) Push(data interface{}) bool {
	click := time.After(time.Millisecond * 20)
	select {
	case dc.Queue <- data:
		return true
	case <-click:
		return false
	}
}

//非阻塞pop
func (dc *DataContainer) Pop() (data interface{}) {
	click := time.After(time.Millisecond * 20)
	select {
	case data = <-dc.Queue:
		return data
	case <-click:
		return nil
	}
}
