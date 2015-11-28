package controllers

import (
	"github.com/dzhcool/eye"
)

type MainController struct {
	eye.Controller
}

func (p *MainController) Get() {
	p.Req.Output.Status = 200
	p.Req.WriteString("hello world, Get")
}

func (p *MainController) Post() {
	p.Req.Output.Status = 200
	p.Req.WriteString("hello world, Post")
}

func (p *MainController) Hello() {
	p.Req.Output.Status = 200
	p.Req.WriteString("hello world, hello")
}
