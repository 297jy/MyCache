package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		fmt.Println("hello")
		go func() {
			time.Sleep(30 * time.Second)
			context.String(http.StatusOK, "hello World!")
		}()
	})
	_ = r.Run(":9090")
}
