package eye

import (
	"fmt"
	"strconv"
	"strings"
)

const VERSION = "1.0.0"

var (
	EyeApp    *App
	HttpAddr  string
	HttpPort  int
	ListenTCP bool
	Graceful  bool
)

/*func Router(uri string, c ControllerInterface, mappingMethods ...string) *App {

}*/

type Eye struct {
	Host string
}

func newEye(host string) *Eye {
	return &Eye{Host: host}
}

func Run(params ...string) {
	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			HttpAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			HttpPort, _ = strconv.Atoi(strs[1])
		}
	}
}

func Rounter(uri string, controller *IController) {
	fmt.Println(uri, controller)
}

//////////// 包初始化  /////////////
func init() {
	//初始化配置
	EyeApp = NewApp()
	HttpAddr = ""
	HttpPort = 9002
	ListenTCP = true
	Graceful = true
}
