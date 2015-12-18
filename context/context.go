// Usage:
//	ctx := context.Context{Request:req,ResponseWriter:rw}
package context

import (
	"bytes"
	"net/http"
	"text/template"
)

// Http request context struct including BeegoInput, BeegoOutput, http.Request and http.ResponseWriter.
// BeegoInput and BeegoOutput provides some api to operate request and response more easily.
type Context struct {
	Input          *EyeInput
	Output         *EyeOutput
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// Redirect does redirection to localurl with http header status code.
// It sends http response header directly.
func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Output.Header("Location", localurl)
	ctx.ResponseWriter.WriteHeader(status)
}

// Abort stops this request.
// if beego.ErrorMaps exists, panic body.
func (ctx *Context) Abort(status int, body string) {
	ctx.ResponseWriter.WriteHeader(status)
	panic(body)
}

// Write string to response body.
// it sends response body.
func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

// Get cookie from request by a given key.
// It's alias of BeegoInput.Cookie.
func (ctx *Context) GetCookie(key string) string {
	return ctx.Input.Cookie(key)
}

// Set cookie for response.
// It's alias of BeegoOutput.Cookie.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	ctx.Output.Cookie(name, value, others...)
}

//json接口返回
func (ctx *Context) WriteJson(content []byte) error {
	ctx.Output.Header("Content-Type", "application/x-javascript; charset=utf-8")
	ctx.Output.Status = 200
	callback := template.JSEscapeString(ctx.Input.Query("_callback"))
	if len(callback) > 0 {
		bt := bytes.NewBufferString(" " + callback)
		bt.WriteString("(")
		bt.Write(content)
		bt.WriteString(");\r\n")
		ctx.Output.Body(bt.Bytes())
	} else {
		ctx.Output.Body(content)
	}

	return nil
}
