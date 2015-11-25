package eye

import ()

type ControllerRegistor struct {
}

type Controller struct {
	controllerName string
	actionName     string
	app            interface{}
}

type IController interface {
	Init(controllerName, actionName string, app interface{})
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
}

func NewControllerRegister() {
}

func (p *Controller) Init(controllerName, actionName string, app interface{}) {
	p.controllerName = controllerName
	p.actionName = actionName
	p.app = app
}
