package gee

import (
	"fmt"
	"net/http"
)

//context是一个上下文处理器
//通过接收到了request来构造respond

type Context struct {
	resp   http.ResponseWriter
	req    *http.Request
	path   string
	method string
	// response info
	statusCode int
}

func New(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		resp:   w,
		req:    r,
		path:   r.URL.Path,
		method: r.Method}
}

func (c *Context) Status(int code) {

}

func (c *Context) Handle() {

}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
