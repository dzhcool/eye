package eye

import (
	"log/syslog"
	"runtime"
)

var (
	loggers  map[string]*XLogger //日志
	logName  string
	logLevel string
)

const callerLevel = 3
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone
)

var LoggerLevel = map[string]int{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"none":  LevelNone,
}

type XLogger struct {
	logName   string
	logLevel  int
	logWriter *syslog.Writer
}

func NewLogger(tag string) *XLogger {
	return Log(tag)
}

func Log(tag string) *XLogger {
	var logger *XLogger
	var exist bool

	if logger, exist = loggers[tag]; !exist {
		logger = &XLogger{}
		logger.init(tag)
		loggers[tag] = logger
	}

	return logger
}

func (this *XLogger) init(tag string) {
	this.logName = logName
	this.logLevel = LoggerLevel[logLevel]
	this.logWriter = getWriter(this.logName + "/" + tag)
}

func getWriter(logName string) *syslog.Writer {
	writer, _ := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL6, logName)
	return writer
}

func (this *XLogger) Logger() *syslog.Writer {
	if this.logName == "" {
		panic("XLogger log name missing")
	}
	if this.logWriter == nil {
		panic("XLogger log writer missing")
	}
	return this.logWriter
}

func (this *XLogger) Debug(str string, evts ...string) {
	evt := this.caller(callerLevel, evts...)
	if this.logLevel <= LevelDebug {
		this.Logger().Info("evt[" + evt + "] [debug] " + str)
	}
}
func (this *XLogger) Info(str string, evts ...string) {
	evt := this.caller(callerLevel, evts...)
	if this.logLevel <= LevelInfo {
		this.Logger().Info("evt[" + evt + "] [info] " + str)
	}
}
func (this *XLogger) Warn(str string, evts ...string) {
	evt := this.caller(callerLevel, evts...)
	if this.logLevel <= LevelWarn {
		this.Logger().Info("evt[" + evt + "] [warn] " + str)
	}
}
func (this *XLogger) Error(str string, evts ...string) {
	evt := this.caller(callerLevel, evts...)
	if this.logLevel <= LevelError {
		this.Logger().Info("evt[" + evt + "] [error] " + str)
	}
}

//获取调用者方法名
func (this *XLogger) caller(level int, evts ...string) string {
	evt := ""
	if len(evts) <= 0 {
		pc, _, _, _ := runtime.Caller(level)
		evt = runtime.FuncForPC(pc).Name()
	} else {
		evt = evts[0]
	}
	return evt
}

func init() {
	logName = Env["PRJ_NAME"]
	logLevel = Env["LOGLEVEL"]
	loggers = make(map[string]*XLogger)
}
