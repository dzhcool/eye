package eye

import (
	"github.com/dzhcool/eye/context"
	"net/http"
)

type Controller struct {
	controllerName string
	actionName     string
	AppController  interface{}
	Req            *context.Context
	Resp           []map[interface{}]interface{}
}

type IController interface {
	Init(context *context.Context, actionName string, app IController)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Options()
	Finish()
	Display() error
	HandlerFunc(fn string) bool
	Name() string
}

func (p *Controller) Init(context *context.Context, actionName string, app IController) {
	p.Req = context
	p.controllerName = app.Name()
	p.actionName = actionName
	p.AppController = app
}

func (p *Controller) Prepare() {
}

func (p *Controller) Get() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Post() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Delete() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Put() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Head() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Patch() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Options() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Finish() {
	p.Req.Output.SetStatus(405)
	http.Error(p.Req.ResponseWriter, "Method Not Allowed", 405)
}

func (p *Controller) Display() error {
	return nil
}

func (p *Controller) HandlerFunc(fn string) bool {
	return false
}

func (p *Controller) Name() string {
	return p.controllerName
}
