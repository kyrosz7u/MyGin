package gee

import (
	"fmt"
	"net/http"
)

type Engine struct {
	route map[string]http.HandlerFunc
}

func (e *Engine) addRoute(method, name string, handler http.HandlerFunc) {
	e.route[method+"+"+name] = handler
}

func (e *Engine) POST(name string, handler http.HandlerFunc) {
	e.addRoute("POST", name, handler)
}

func (e *Engine) GET(name string, handler http.HandlerFunc) {
	e.addRoute("GET", name, handler)
}

func (e *Engine) Run(port string) {
	http.ListenAndServe(port, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "+" + req.URL.Path
	if handler, ok := e.route[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 Not Found")
	}
}
