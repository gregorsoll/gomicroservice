package apicontroller

import "github.com/gin-gonic/gin"

func IndexRoutes(superRoute *gin.RouterGroup, baseroute string) {
	booksRouter := superRoute.Group(baseroute)
	{
		booksRouter.GET("/", Index)
		booksRouter.GET("/ping", PingPong)
	}
}
