package indexController

import "github.com/gin-gonic/gin"

func Index(c *gin.Context) {
	c.Writer.WriteString("Hello world")
	c.Status(200)
	c.Done()
}

func PingPong(c *gin.Context) {
	c.Writer.WriteString("Pingpong")
	c.Status(200)
	c.Done()
}
