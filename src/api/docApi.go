package api

import (
	"l2-golang-auth/doc"

	"github.com/gin-gonic/gin"
)

// DocContent Send Doc Data
func DocContent(c *gin.Context) {
	data, err := doc.GetAllDocPosts()
	if err != nil {
		c.JSON(200, gin.H{"error": true, "data": err})
	} else {
		c.JSON(200, gin.H{"error": false, "data": data})
	}
}

// DocIDContent Send Doc Data For that Id
func DocIDContent(c *gin.Context, id string) {
	data, err := doc.Fetcharticle(id)
	if err != nil {
		c.JSON(200, gin.H{"error": true, "data": err})
	} else {
		c.JSON(200, gin.H{"error": false, "data": data[id]})
	}
}

// Docsearch Send Doc Data
func Docsearch(c *gin.Context, id string) {
	data, err := doc.Searchtext(id)
	if err != nil {
		c.JSON(200, gin.H{"error": true, "data": err})
	} else {
		c.JSON(200, gin.H{"error": false, "data": data})
	}
}
