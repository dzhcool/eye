package eye

import (
	"fmt"
	"github.com/dzhcool/eye/endless"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	Handlers *ControllerRegistor
	Server   *http.Server
}

func NewApp() *App {
	cr := NewControllerRegister()
	return &App{Handlers: cr, Server: &http.Server{}}
}

func (p *App) Run() {
	running := make(chan bool, 1)
	addr := Env["GOADDR"]

	//启动程序监听
	if Env["GRACEFUL"] == "1" {
		go func() {
			endless.DefaultReadTimeOut = 3 * time.Second
			endless.DefaultWriteTimeOut = 5 * time.Second
			err := endless.ListenAndServe(addr, p.Handlers)
			if err != nil {
				log.Println("[Eey]Listen error:", err)
			}
			log.Println(fmt.Sprintf("[Eye]Server on %s stopped", addr))
			running <- true
		}()
	} else {
		go func() {
			server := http.Server{
				Addr:         addr,
				Handler:      p.Handlers,
				ReadTimeout:  3 * time.Second,
				WriteTimeout: 5 * time.Second,
			}
			err := server.ListenAndServe()
			if err != nil {
				log.Println("[Eey]Listen error:", err)
			}
			log.Println(fmt.Sprintf("[Eye]Server on %s stopped", addr))
			running <- true
		}()
	}
	<-running
	log.Println("[Eye]All servers stopped. Exiting.")
	os.Exit(0)
}
