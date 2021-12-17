package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	routeMap map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{routeMap: make(map[string]HandlerFunc)}
}

func (r *Router) addRoute(method, name string, handler HandlerFunc) {
	key := method + "+" + name
	r.routeMap[key] = handler
}

func (r *Router) handle(ctx *Context) {
	key := ctx.Method + "+" + ctx.Path
	if handler, ok := r.routeMap[key]; ok {
		handler(ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 Not Found:%s", ctx.Path)
	}
}

//将路径拆分成单词列表，用来查找前缀路由树
//只返回用于匹配的单词，/hello/*filepath这种，忽略filepath，返回/hello/*
func parsePath(path string) []string {
	vs := strings.Split(path, "/")
	ret := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			ret = append(ret, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return ret
}
