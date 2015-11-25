package eye

import (
	"fmt"
	"net/http"
)

type App struct {
	Handlers *ControllerRegistor
	Server   *http.Server
}

func NewApp() *App {
	// cr := NewControllerRegister()
	// app := &App{Handlers: cr, Server: &http.Server{}}
	return &App{}
}

func (p *App) Run() {
	addr := HttpAddr
	if HttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	}
	fmt.Println(addr)
}
