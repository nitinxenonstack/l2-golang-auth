package noroute

import (
	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {

	c.HTML(404, "404-not-found", gin.H{
		"linktohome": "/login",
	})
	c.Abort()
}
