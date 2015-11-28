package main

import (
	"controllers"
	"github.com/dzhcool/eye"
)

func main() {
	eye.Router("/", &controllers.MainController{}, "hello")
	eye.Router("/rest", &controllers.MainController{})
	eye.Run()
}
