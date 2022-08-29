package gee

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})
	r.GET("/hello/:id", func(c *gin.Context) {
		param := c.Param("id")
		c.String(http.StatusOK, param)
	})
	v1 := r.Group("/")

	v1.GET("hell", func(c *gin.Context) {
		c.String(http.StatusOK, c.FullPath())
	})
	v1.GET("test", func(c *gin.Context) {
		c.String(http.StatusOK, c.FullPath())
	})

	err := r.Run(":9999")
	if err != nil {
		return
	}
}
