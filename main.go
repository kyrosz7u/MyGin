package main

import (
	"fmt"
	"gee"
	"net/http"
)

func strr(ctx *gee.Context) {
	ctx.String(http.StatusOK, "you're at %s", ctx.Path)
}

func main() {
	engine := gee.NewEngine()
	engine.Use(gee.Logger)
	engine.GET("/", strr)

	v1 := engine.NewGroup("/v1")
	v1.Use(func(c *gee.Context) {
		fmt.Println("v1's middleware executed")
		c.String(http.StatusOK, "v1's middleware executed\n")
	})
	v1.GET("/233", strr)

	v2 := engine.NewGroup("/v2")
	v2.Use(func(c *gee.Context) {
		fmt.Println("v2's middleware executed")
		c.String(http.StatusOK, "v2's middleware executed\n")
	})
	v2.GET("/133", strr)

	engine.GET("/:lang", func(c *gee.Context) {
		//TODO c.PostForm(c.Paras["lang"])是用来解析Post参数列表的吗？
		//fmt.Println(c.PostForm(c.Paras["lang"]))
		c.JSON(http.StatusOK, map[string]interface{}{
			"lang:": c.Paras["lang"],
		})
	})
	engine.Run(":9000")
}
