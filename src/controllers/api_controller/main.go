package apicontroller

import "github.com/gin-gonic/gin"

func Index(c *gin.Context) {
	c.Writer.WriteString("api world")
	c.Status(200)
	c.Done()
}

func PingPong(c *gin.Context) {
	c.Writer.WriteString("api pong")
	c.Status(200)
	c.Done()
}
