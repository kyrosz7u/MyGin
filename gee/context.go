package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//context是一个上下文处理器
//通过接收到了request来构造respond

type Context struct {
	resp   http.ResponseWriter
	req    *http.Request
	Path   string
	Method string
	// response info
	StatusCode int
	Paras      map[string]string
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		resp:   w,
		req:    r,
		Path:   r.URL.Path,
		Method: r.Method}
}

func (c *Context) PostForm(key string) string {
	return c.req.FormValue(key)
}

// 查询URL中的参数
func (c *Context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

//http.ResponseWriter.StatusCode的作用是设置响应的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.resp.WriteHeader(code)
}

//SetHeader: 设置响应报文的头部字段
//key: header
func (c *Context) SetHeader(key, value string) {
	c.resp.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.resp.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.resp)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.resp, err.Error(), 500)
	}
}
