package eye

import (
	"log"
	"path"
	"runtime"
)

const VERSION = "1.0.0"

var (
	EyeApp *App
)

// addr:Tcp:127.0.0.1:9001 or Tcp::9001
/*func Listen(addr string, Mux *Router) {
	sAddr := strings.Split(addr, ":")
	HttpProtocol := sAddr[0]
	HttpAddr := fmt.Sprintln("%s:%s", sAddr[1], sAddr[2])

	EyeApp.AddServer(HttpProtocol, HttpAddr, Mux)
}*/

//run
/*func Run(addr string, Mux *Router) {
	sAddr := strings.Split(addr, ":")
	HttpProtocol := sAddr[0]
	HttpAddr := fmt.Sprintln("%s:%s", sAddr[1], sAddr[2])
	EyeApp.AddServer(HttpProtocol, HttpAddr, Mux)

	EyeApp.Run()
}*/

func Router(rootpath string, c IController, mappingMethods ...string) *App {
	EyeApp.Handlers.Add(rootpath, c, mappingMethods...)
	return EyeApp
}

func RESTRouter(rootpath string, c IController) *App {
	Router(rootpath, c)
	Router(path.Join(rootpath, ":objectId"), c)
	return EyeApp
}

func Run() {
	EyeApp.Run()
}

//////////// 包初始化  /////////////
func init() {
	//初始化默认配置
	EyeApp = NewApp()

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("[Eye]Run As NumCPU:", runtime.NumCPU())
}
