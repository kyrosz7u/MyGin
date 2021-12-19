package gee

import (
	"log"
	"time"
)

func Logger(ctx *Context) {
	log.Printf("Recived Request: Method:[%s] URL:[%s]", ctx.Method, ctx.Path)
	t := time.Now()
	//处理完剩下的再返回
	ctx.Next()
	log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.req.RequestURI, time.Since(t))
}
