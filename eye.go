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
	EyeApp = NewApp()
	//初始化默认配置
	StartPprof()

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("[Eye]Run As NumCPU:", runtime.NumCPU())
}
