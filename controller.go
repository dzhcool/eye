package eye

import (
	"fmt"
	"github.com/dzhcool/eye/context"
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
}

func (p *Controller) Post() {
}

func (p *Controller) Delete() {
}

func (p *Controller) Put() {
}

func (p *Controller) Head() {
}

func (p *Controller) Patch() {
}

func (p *Controller) Options() {
}

func (p *Controller) Finish() {
}

func (p *Controller) Display() error {
	fmt.Println("[controller]Display")
	return nil
}

func (p *Controller) HandlerFunc(fn string) bool {
	return false
}

func (p *Controller) Name() string {
	return p.controllerName
}
