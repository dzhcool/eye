package qhydra

import (
	"encoding/json"
	"log/syslog"
	"os"
	"time"
)

var QMsg *Qhydra = nil

type Qhydra struct {
	event    string
	tag      string
	hostname string
	priority syslog.Priority
	writer   *syslog.Writer
}

const EVENT_PREFIX string = "xxx_"

func NewQhydra(event string) *Qhydra {
	if QMsg == nil {
		QMsg = &Qhydra{
			event:    event,
			tag:      EVENT_PREFIX + event,
			priority: syslog.LOG_LOCAL4 | syslog.LOG_INFO,
		}
		QMsg.init()
	}
	return QMsg
}

func (this *Qhydra) init() {
	if writer, err := syslog.New(this.priority, this.tag); err != nil {
		panic("init qhydra fail:" + err.Error())
	} else {
		this.writer = writer
	}

	if hostname, err := os.Hostname(); err != nil {
		panic("init qhydra fail:" + err.Error())
	} else {
		this.hostname = hostname
	}
}

func (this *Qhydra) Trigger(data interface{}, key string) ([]byte, int, error) {

	msg := map[string]interface{}{
		"name": this.event,
		"data": data,
		"host": this.hostname,
		"key":  key,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}

	msgJson, _ := json.Marshal(msg)
	num, err := this.writer.Write(msgJson)
	return msgJson, num, err
}

func (this *Qhydra) Close() error {
	return this.writer.Close()
}
