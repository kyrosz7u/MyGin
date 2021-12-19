package gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

//可以有多个Group绑定同一个Engine
type Group struct {
	engine      *Engine
	prefix      string        //Group的前缀
	midHandlers []HandlerFunc //中间件
}

//Go语言中可以用结构体嵌套的方式实现类似C++中的继承
//Engine相当于Group的子类？
type Engine struct {
	*Group
	groupList []*Group
	router    *Router
}

func NewEngine() *Engine {
	e := &Engine{
		groupList: make([]*Group, 0),
		router:    newRouter(),
	}
	e.Group = &Group{
		engine:      e,
		prefix:      "",
		midHandlers: make([]HandlerFunc, 0),
	}
	return e
}

func (e *Engine) NewGroup(prefix string) *Group {
	//e.midHandlers是全局的中间件
	//这里需要拷贝到Group中
	g := &Group{e, prefix, make([]HandlerFunc, 0)}
	e.groupList = append(e.groupList, g)
	return g
}

func (g *Group) addRoute(method, path string, handler HandlerFunc) {
	g.engine.router.addRoute(method, g.prefix+path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	g.engine.router.addRoute("POST", g.prefix+path, handler)
}

func (g *Group) GET(path string, handler HandlerFunc) {
	g.engine.router.addRoute("GET", g.prefix+path, handler)
}

//注册要使用的中间件
func (g *Group) Use(midWare HandlerFunc) {
	g.midHandlers = append(g.midHandlers, midWare)
}

func (e *Engine) FindGroup(URLPath string) *Group {
	list := parsePath(URLPath)
	for _, g := range e.groupList {
		if g.prefix == "/"+list[0] {
			return g
		}
	}
	return nil
}

func (e *Engine) Run(port string) {
	http.ListenAndServe(port, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//根据URL的查询对应的分组
	g := e.FindGroup(r.URL.Path)
	var ctx *Context

	if g == nil {
		//如果没查找到匹配的分组，使用全局中间件
		ctx = newContext(w, r, e.midHandlers)
	} else {
		//如果查找到匹配的分组，使用全局中间件+分组的中间件
		ctx = newContext(w, r, append(e.midHandlers, g.midHandlers...))
	}

	e.router.handle(ctx)
}
