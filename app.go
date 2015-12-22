package eye

import (
	"fmt"
	"github.com/dzhcool/eye/endless"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/exec"
	"strings"
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
	addr := ":" + os.Getenv("GOPORT")
	lnet := os.Getenv("GONET")
	gosock := os.Getenv("GOSOCK")

	//启动程序监听
	if strings.ToUpper(os.Getenv("GRACEFUL")) == "ON" {
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
		if strings.ToUpper(lnet) == "TCP" {
			log.Println("[Eey]Listen as", lnet, addr)
			go func() {
				server := http.Server{
					Addr:         addr,
					Handler:      p.Handlers,
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 5 * time.Second,
				}
				err := server.ListenAndServe()
				if err != nil {
					log.Println("[Eey]Listen error:", err)
				}
				log.Println(fmt.Sprintf("[Eye]Server on %s stopped", addr))
				running <- true
			}()
		} else {
			log.Println("[Eey]Listen as", lnet, gosock)
			go func() {
				exec.Command("/bin/sh", "-c", "rm "+gosock).Run()
				unix, err := net.Listen("unix", gosock)
				exec.Command("/bin/sh", "-c", "chmod a+w "+gosock).Run()
				if err != nil {
					log.Println("[Eey]Listen error:", err)
				}
				fcgi.Serve(unix, p.Handlers)
				log.Println(fmt.Sprintf("[Eye]Server on %s stopped", addr))
				running <- true
			}()
		}
	}
	<-running
	log.Println("[Eye]All servers stopped. Exiting.")
	os.Exit(0)
}
