package eye

import (
	eyecontext "github.com/dzhcool/eye/context"
	// "log"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

var (
	HTTPMETHOD = map[string]string{
		"GET":     "GET",
		"POST":    "POST",
		"PUT":     "PUT",
		"DELETE":  "DELETE",
		"PATCH":   "PATCH",
		"OPTIONS": "OPTIONS",
		"HEAD":    "HEAD",
		"TRACE":   "TRACE",
		"CONNECT": "CONNECT",
	}
)

type ControllerRegistor struct {
	Mux map[string]*routerItem
}

type routerItem struct {
	Controller     IController
	ControllerType reflect.Type
	Methods        map[string]string
	Handler        http.Handler
}

func NewControllerRegister() *ControllerRegistor {
	return &ControllerRegistor{
		Mux: make(map[string]*routerItem),
	}
}

func (p *ControllerRegistor) Add(rootpath string, c IController, mappingMethods ...string) {
	rootpath = p.Path(strings.ToLower(rootpath))
	item, ok := EyeApp.Handlers.Mux[rootpath]
	if !ok {
		reflectVal := reflect.ValueOf(c)
		t := reflect.Indirect(reflectVal).Type()
		item = &routerItem{Controller: c, ControllerType: t, Methods: make(map[string]string), Handler: nil}
		EyeApp.Handlers.Mux[rootpath] = item
	}
	//"post,get:Users;delete:Del"
	if len(mappingMethods) > 0 {
		multi := strings.Split(mappingMethods[0], ";")
		for _, single := range multi {
			mappings := strings.Split(single, ":")
			if len(mappings) == 1 {
				item.Methods["get"] = mappings[0]
			}
			if len(mappings) == 2 {
				methods := strings.Split(mappings[0], ",")
				for _, method := range methods {
					item.Methods[strings.ToLower(method)] = mappings[1]
				}
			}
		}
	} else {
		for k, v := range HTTPMETHOD {
			fun := strings.ToUpper(v[0:1]) + strings.ToLower(v[1:])
			item.Methods[strings.ToLower(k)] = fun
		}
	}

}

func (p *ControllerRegistor) getRouter(rootpath string, method string) (string, error) {
	rootpath = p.Path(rootpath)
	if item, ok := EyeApp.Handlers.Mux[rootpath]; ok {
		if fun, ok := item.Methods[strings.ToLower(method)]; ok {
			return fun, nil
		}
	}

	return "", errors.New("Router not found")
}

func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w := &responseWriter{writer: rw}
	//inti
	context := &eyecontext.Context{
		ResponseWriter: w,
		Request:        r,
		Input:          eyecontext.NewInput(r),
		Output:         eyecontext.NewOutput(),
	}
	context.Output.Context = context
	context.Output.EnableGzip = false

	path := p.Path(r.URL.Path)
	ritem, ok := p.Mux[path]
	if !ok {
		panic("route is not exist")
	}
	vc := reflect.New(ritem.ControllerType)
	execController, ok := vc.Interface().(IController)
	if !ok {
		panic("controller is not IController")
	}
	actionName, err := p.getRouter(path, r.Method)
	if err != nil {
		panic("action is not exist")
	}

	//call controller Init()
	execController.Init(context, actionName, execController)

	//call controller Prepare()
	execController.Prepare()

	//call action
	if len(actionName) > 0 {
		in := make([]reflect.Value, 0)
		method := vc.MethodByName(actionName)
		method.Call(in)
	}
	if context.Output.Status == 0 {
		http.NotFound(rw, r)
	}
}

func (p *ControllerRegistor) Path(path string) string {
	path = strings.ToLower(strings.TrimRight(path, "/"))
	if len(path) <= 0 {
		path = "/"
	}
	return path
}

//responseWriter is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type responseWriter struct {
	writer  http.ResponseWriter
	started bool
	status  int
}

// Header returns the header map that will be sent by WriteHeader.
func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (w *responseWriter) Write(p []byte) (int, error) {
	w.started = true
	return w.writer.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.started = true
	w.writer.WriteHeader(code)
}
