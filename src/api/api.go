package api

import "github.com/gin-gonic/gin"

func NewApp() *gin.Engine {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return r

}
