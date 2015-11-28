package eye

import (
	"log/syslog"
	"os"
)

var (
	loggers  map[string]*XLogger //日志
	logName  string
	logLevel string
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var LoggerLevel = map[string]int{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
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

func (this *XLogger) Debug(str string) {
	if this.logLevel <= LevelDebug {
		this.Logger().Info(" [debug] " + str)
	}
}
func (this *XLogger) Info(str string) {
	if this.logLevel <= LevelInfo {
		this.Logger().Info(" [info] " + str)
	}
}
func (this *XLogger) Warn(str string) {
	if this.logLevel <= LevelWarn {
		this.Logger().Info(" [warn] " + str)
	}
}
func (this *XLogger) Error(str string) {
	if this.logLevel <= LevelError {
		this.Logger().Info(" [error] " + str)
	}
}

func init() {
	logName = os.Getenv("PRJ_NAME")
	logLevel = os.Getenv("LOGLEVEL")
	loggers = make(map[string]*XLogger)
}
