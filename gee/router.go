package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	rootMap   map[string]*node
	handleMap map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{handleMap: make(map[string]HandlerFunc), rootMap: make(map[string]*node)}
}

func (r *Router) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "+" + pattern
	p := parsePath(pattern)
	_, ok := r.rootMap[method]
	if !ok {
		r.rootMap[method] = &node{children: make([]*node, 0)}
	}
	r.rootMap[method].insert(pattern, p, 0)
	r.handleMap[key] = handler
}

//查找路由树，获取动态匹配的路由，匹配成功返回node和参数
//参数的格式：/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，
///static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)
	root, ok := r.rootMap[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePath(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index][1:]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *Router) handle(ctx *Context) {
	node, paras := r.getRoute(ctx.Method, ctx.Path)

	if node == nil {
		ctx.String(http.StatusNotFound, "404 Not Found:%s", ctx.Path)
	} else {
		key := ctx.Method + "+" + node.pattern
		ctx.Paras = paras
		r.handleMap[key](ctx)
	}
}

func (r *Router) PrintRouters(method string) []*node {
	root, ok := r.rootMap[method]
	if !ok {
		return nil
	}
	list := make([]*node, 0)
	root.travel(&list)
	return list
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
