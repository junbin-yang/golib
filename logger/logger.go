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
	"syscall"
	"time"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

var obj *logrus.Logger = logrus.New()
var defaultChannel *DataContainer = InitDataContainer()
var channelSwitch bool

type Options struct {
	AppName  string
	Path     string
	Level    logrus.Level
	Rotate   bool
	KeepDays int64
	TakeStd  bool
}

func (this *Options) New() {
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
			obj.SetOutput(src)
			if this.TakeStd {
				os.Stderr = src
				os.Stdout = src
				syscall.Dup3(int(src.Fd()), 2, 0)
				syscall.Dup3(int(src.Fd()), 1, 0)
			}

			if this.Rotate {
				go func(Path, AppName string) {
					c := cron.New()
					c.AddFunc("0 0 0 * * ?", func() {
						logfile := Path + "/" + AppName + "-" + carbon.Now().ToDateString() + ".log"
						src, _ := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|0644)
						obj.SetOutput(src)

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

func SetLogLevel(level logrus.Level) {
	obj.SetLevel(level)
}

func Asyn() {
	channelSwitch = true
	for {
		itemInterface := defaultChannel.Pop()
		if itemInterface != nil {
			item := itemInterface.(map[string]interface{})
			m := fmt.Sprint(item["Data"])
			if item["Type"] == PanicLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Panic(m)
			}
			if item["Type"] == InfoLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Info(m)
			}
			if item["Type"] == DebugLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Debug(m)
			}
			if item["Type"] == ErrorLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Error(m)
			}
			if item["Type"] == WarnLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Warn(m)
			}
			if item["Type"] == FatalLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Fatal(m)
			}
			if item["Type"] == TraceLevel {
				obj.WithFields(map[string]interface{}{"Func": item["Func"], "File": item["File"], "Line": item["Line"], "GID": item["GID"]}).Trace(m)
			}
		}
	}
}

func Trace(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": TraceLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Trace(m)
	}
}

func Panic(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": PanicLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Panic(m)
	}
}

func Info(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": InfoLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Info(m)
	}
}

func Debug(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": DebugLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Debug(m)
	}
}

func Error(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": ErrorLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := ""
		for _, v := range o {
			m += fmt.Sprint(v) + " "
		}
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Error(m)
	}
}

func Warn(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": WarnLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}
	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := ""
		for _, v := range o {
			m += fmt.Sprint(v) + " "
		}
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Warn(m)
	}
}

func Fatal(o ...interface{}) {
	m := ""
	for _, v := range o {
		m += fmt.Sprint(v) + " "
	}

	msg := map[string]interface{}{
		"Type": FatalLevel,
		"Data": m,
	}

	if obj.GetLevel() >= logrus.DebugLevel {
		fun, file, line := printCaller()
		msg["Func"] = fun
		msg["File"] = file
		msg["Line"] = line
		msg["GID"] = getGID()
	}

	if channelSwitch {
		defaultChannel.Push(msg)
	} else {
		m := ""
		for _, v := range o {
			m += fmt.Sprint(v) + " "
		}
		obj.WithFields(map[string]interface{}{"Func": msg["Func"], "File": msg["File"], "Line": msg["Line"], "GID": msg["GID"]}).Fatal(m)
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
	t := carbon.Now().ToDateTimeString() + "," + fmt.Sprint(time.Now().Nanosecond()/1e6)
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
