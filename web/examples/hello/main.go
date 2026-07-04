package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-minstack/go-minstack/core"
	mgin "github.com/go-minstack/go-minstack/gin"
	"github.com/go-minstack/go-minstack/web"
)

func registerRoutes(r *gin.Engine) {
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, web.NewMessageDto("Hello from MinStack!"))
	})

	r.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, web.NewErrorDto(fmt.Errorf("something went wrong")))
	})
}

func main() {
	app := core.New(mgin.Module())
	app.Invoke(registerRoutes)
	app.Run()
}
