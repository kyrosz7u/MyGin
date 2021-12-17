package main

import (
	"gee"
	"net/http"
)

func strr(ctx *gee.Context) {
	ctx.String(http.StatusOK, "you're at %s", ctx.Path)
}

func main() {
	engine := gee.New()
	engine.GET("/233", strr)
	engine.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	engine.Run(":9000")
}
